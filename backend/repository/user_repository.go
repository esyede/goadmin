package repository

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/util"
	"github.com/esyede/goadmin/backend/vo"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/thoas/go-funk"
)

type IUserRepository interface {
	Login(user *model.User) (*model.User, error)       // Log in
	ChangePwd(username string, newPasswd string) error // Update password

	CreateUser(user *model.User) error                              // Create user
	GetUserById(id uint) (model.User, error)                        // Get a single user
	GetUsers(req *vo.UserListRequest) ([]*model.User, int64, error) // Get user list
	UpdateUser(user *model.User) error                              // update user
	BatchDeleteUserByIds(ids []uint) error                          // batch deletion

	GetCurrentUser(c *gin.Context) (model.User, error)                  // Get the current logged in user information
	GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) // Get the minimum value of the current user role sorting (the highest level role) and the current user information
	GetUserMinRoleSortsByIds(ids []uint) ([]int, error)                 // Get the minimum value of user role sorting based on user ID

	SetUserInfoCache(username string, user model.User) // Set user information cache
	UpdateUserInfoCacheByRoleId(roleId uint) error     // Update the user information cache of the role based on the role ID
	ClearUserInfoCache()                               // Clear all user information cache
}

type UserRepository struct {
}

// Cache current user information to avoid frequent database acquisitions
var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

// UserRepository constructor
func NewUserRepository() IUserRepository {
	return UserRepository{}
}

// Log in
func (ur UserRepository) Login(user *model.User) (*model.User, error) {
	// Get the user based on the user name (normal status: user status is normal)
	var firstUser model.User
	err := common.DB.
		Where("username = ?", user.Username).
		Preload("Roles").
		First(&firstUser).Error
	if err != nil {
		return nil, errors.New("User does not exist")
	}

	// Determine the user's status
	userStatus := firstUser.Status
	if userStatus != 1 {
		return nil, errors.New("User is banned")
	}

	// Determine the status of all roles owned by the user. If all roles are disabled, you cannot log in.
	roles := firstUser.Roles
	isValidate := false
	for _, role := range roles {
		// You can log in if you have a normal character
		if role.Status == 1 {
			isValidate = true
			break
		}
	}

	if !isValidate {
		return nil, errors.New("User role is disabled")
	}

	// Verify password
	err = util.ComparePasswd(firstUser.Password, user.Password)
	if err != nil {
		return &firstUser, errors.New("Wrong password")
	}
	return &firstUser, nil
}

// Get the current logged in user information
// Need caching to reduce database access
func (ur UserRepository) GetCurrentUser(c *gin.Context) (model.User, error) {
	var newUser model.User
	ctxUser, exist := c.Get("user")
	if !exist {
		return newUser, errors.New("User is not logged in")
	}
	u, _ := ctxUser.(model.User)

	// Get the cache first
	cacheUser, found := userInfoCache.Get(u.Username)
	var user model.User
	var err error
	if found {
		user = cacheUser.(model.User)
		err = nil
	} else {
		// Get the database if it is not in the cache
		user, err = ur.GetUserById(u.ID)
		// Cache if successful
		if err != nil {
			userInfoCache.Delete(u.Username)
		} else {
			userInfoCache.Set(u.Username, user, cache.DefaultExpiration)
		}
	}
	return user, err
}

// Get the minimum value of the current user role sorting (the highest level role) and the current user information
func (ur UserRepository) GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	// Get the current user
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		return 999, ctxUser, err
	}
	// Get all roles of the current user
	currentRoles := ctxUser.Roles
	// Get the current user role sorting and compare it with the role sorting sent from the front end.
	var currentRoleSorts []int
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
	}
	// Minimum value of current user role sorting (highest level role)
	currentRoleSortMin := uint(funk.MinInt(currentRoleSorts).(int))

	return currentRoleSortMin, ctxUser, nil
}

// Get a single user
func (ur UserRepository) GetUserById(id uint) (model.User, error) {
	fmt.Println("GetUserById---")
	var user model.User
	err := common.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

// Get user list
func (ur UserRepository) GetUsers(req *vo.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := common.DB.Model(&model.User{}).Order("created_at DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// Paging only when pageNum > 0 and pageSize > 0
	// Total number of records
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

// Update password
func (ur UserRepository) ChangePwd(username string, hashNewPasswd string) error {
	err := common.DB.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
	// If the password update is successful, update the current user information cache
	// Get the cache first
	cacheUser, found := userInfoCache.Get(username)
	if err == nil {
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			// Get user information cache without cache
			var user model.User
			common.DB.Where("username = ?", username).First(&user)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}

	return err
}

// Create user
func (ur UserRepository) CreateUser(user *model.User) error {
	err := common.DB.Create(user).Error
	return err
}

// Update user
func (ur UserRepository) UpdateUser(user *model.User) error {
	err := common.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	err = common.DB.Model(user).Association("Roles").Replace(user.Roles)

	// err := common.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user).Error

	// If the update is successful, update the user information cache
	if err == nil {
		userInfoCache.Set(user.Username, *user, cache.DefaultExpiration)
	}
	return err
}

// batch deletion
func (ur UserRepository) BatchDeleteUserByIds(ids []uint) error {
	// There is a many-to-many relationship between users and roles
	var users []model.User
	for _, id := range ids {
		// Get user based on ID
		user, err := ur.GetUserById(id)
		if err != nil {
			return errors.New(fmt.Sprintf("The user with ID %d was not obtained", id))
		}
		users = append(users, user)
	}

	err := common.DB.Select("Roles").Unscoped().Delete(&users).Error
	// If the user is successfully deleted, the user information cache will be deleted.
	if err == nil {
		for _, user := range users {
			userInfoCache.Delete(user.Username)
		}
	}
	return err
}

// Get the minimum value of user role sorting based on user ID
func (ur UserRepository) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	// Get user information based on user ID
	var userList []model.User
	err := common.DB.Where("id IN (?)", ids).Preload("Roles").Find(&userList).Error
	if err != nil {
		return []int{}, err
	}
	if len(userList) == 0 {
		return []int{}, errors.New("No user information was obtained")
	}
	var roleMinSortList []int
	for _, user := range userList {
		roles := user.Roles
		var roleSortList []int
		for _, role := range roles {
			roleSortList = append(roleSortList, int(role.Sort))
		}
		roleMinSort := funk.MinInt(roleSortList).(int)
		roleMinSortList = append(roleMinSortList, roleMinSort)
	}
	return roleMinSortList, nil
}

// Set user information cache
func (ur UserRepository) SetUserInfoCache(username string, user model.User) {
	userInfoCache.Set(username, user, cache.DefaultExpiration)
}

// Update the user information cache of the role based on the role ID
func (ur UserRepository) UpdateUserInfoCacheByRoleId(roleId uint) error {

	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Users").First(&role).Error
	if err != nil {
		return errors.New("Failed to get role information based on role ID")
	}

	users := role.Users
	if len(users) == 0 {
		return errors.New("The user with this role was not obtained based on the role ID.")
	}

	// Update user information cache
	for _, user := range users {
		_, found := userInfoCache.Get(user.Username)
		if found {
			userInfoCache.Set(user.Username, *user, cache.DefaultExpiration)
		}
	}

	return err
}

// Clear all user information cache
func (ur UserRepository) ClearUserInfoCache() {
	userInfoCache.Flush()
}
