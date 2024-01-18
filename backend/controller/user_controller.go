package controller

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/config"
	"github.com/esyede/goadmin/backend/dto"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/repository"
	"github.com/esyede/goadmin/backend/response"
	"github.com/esyede/goadmin/backend/util"
	"github.com/esyede/goadmin/backend/vo"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
)

type IUserController interface {
	GetUserInfo(c *gin.Context)          // Get current logged in user information
	GetUsers(c *gin.Context)             // Get user list
	ChangePwd(c *gin.Context)            // Update user login password
	CreateUser(c *gin.Context)           // Create user
	UpdateUserById(c *gin.Context)       // update user
	BatchDeleteUserByIds(c *gin.Context) // Delete users in batches
}

type UserController struct {
	UserRepository repository.IUserRepository
}

// Constructor
func NewUserController() IUserController {
	userRepository := repository.NewUserRepository()
	userController := UserController{UserRepository: userRepository}
	return userController
}

// Get current logged in user information
func (uc UserController) GetUserInfo(c *gin.Context) {
	user, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain current user information: "+err.Error())
		return
	}
	userInfoDto := dto.ToUserInfoDto(user)
	response.Success(c, gin.H{
		"userInfo": userInfoDto,
	}, "Obtain current user information successfully")
}

// Get user list
func (uc UserController) GetUsers(c *gin.Context) {
	var req vo.UserListRequest
	// Parameter binding
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Parameter verification
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// Obtain
	users, total, err := uc.UserRepository.GetUsers(&req)
	if err != nil {
		response.Fail(c, nil, "Failed to get user list: "+err.Error())
		return
	}
	response.Success(c, gin.H{"users": dto.ToUsersDto(users), "total": total}, "Obtain user list successfully")
}

// Update user login password
func (uc UserController) ChangePwd(c *gin.Context) {
	var req vo.ChangePwdRequest

	// Parameter binding
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Parameter verification
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// The password sent from the front end is RSA encrypted, decrypt it first
	// Password decrypted via RSA
	decodeOldPassword, err := util.RSADecrypt([]byte(req.OldPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	decodeNewPassword, err := util.RSADecrypt([]byte(req.NewPassword), config.Conf.System.RSAPrivateBytes)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	req.OldPassword = string(decodeOldPassword)
	req.NewPassword = string(decodeNewPassword)

	// Get current user
	user, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Get user's real and correct password
	correctPasswd := user.Password
	// Determine whether the password requested by the front end is equal to the real password
	err = util.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		response.Fail(c, nil, "The original password is wrong")
		return
	}
	// Update password
	err = uc.UserRepository.ChangePwd(user.Username, util.GenPasswd(req.NewPassword))
	if err != nil {
		response.Fail(c, nil, "Failed to update password: "+err.Error())
		return
	}
	response.Success(c, nil, "Password updated successfully")
}

// Create user
func (uc UserController) CreateUser(c *gin.Context) {
	var req vo.CreateUserRequest
	// Parameter binding
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Parameter verification
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// Password decrypted via RSA
	// Decrypt if the password is not empty
	if req.Password != "" {
		decodeData, err := util.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
		if err != nil {
			response.Fail(c, nil, err.Error())
			return
		}
		req.Password = string(decodeData)
		if len(req.Password) < 6 {
			response.Fail(c, nil, "Password length must be at least 6 characters")
			return
		}
	}

	// The current user role sorting minimum value (the highest level role) and the current user
	currentRoleSortMin, ctxUser, err := uc.UserRepository.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// Get user role id sent from the front end
	reqRoleIds := req.RoleIds
	// Get role based on the role id
	rr := repository.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain role information based on role ID: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// The front end transmits the minimum value of user role sorting (the highest level role)
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts).(int))

	// The minimum role sorting value of the current user needs to be less than the minimum role sorting value passed by the front end (users cannot create users with a higher level than themselves or with the same level)
	if currentRoleSortMin >= reqRoleSortMin {
		response.Fail(c, nil, "Users cannot create users with a higher level than themselves or with the same level.")
		return
	}

	// If the password is empty, the default is 123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := model.User{
		Username:     req.Username,
		Password:     util.GenPasswd(req.Password),
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}

	err = uc.UserRepository.CreateUser(&user)
	if err != nil {
		response.Fail(c, nil, "Failed to create user: "+err.Error())
		return
	}
	response.Success(c, nil, "User created successfully")

}

// update user
func (uc UserController) UpdateUserById(c *gin.Context) {
	var req vo.CreateUserRequest
	// Parameter binding
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Parameter verification
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// Get userId in path
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		response.Fail(c, nil, "User ID is incorrect")
		return
	}

	// Get user information based on userId in path
	oldUser, err := uc.UserRepository.GetUserById(uint(userId))
	if err != nil {
		response.Fail(c, nil, "Failed to obtain user information that needs to be updated: "+err.Error())
		return
	}

	// Get current user
	ctxUser, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Get all roles of the current user
	currentRoles := ctxUser.Roles
	// Get current user role sorting and compare it with the role sorting sent from the front end.
	var currentRoleSorts []int
	// Current user role ID collection
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// Minimum value of current user role sorting (highest level role)
	currentRoleSortMin := funk.MinInt(currentRoleSorts).(int)

	// Get user role id sent from the front end
	reqRoleIds := req.RoleIds
	// Get role based on the role id
	rr := repository.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain role information based on role ID: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// The front end transmits the minimum value of user role sorting (the highest level role)
	reqRoleSortMin := funk.MinInt(reqRoleSorts).(int)

	user := model.User{
		Model:        oldUser.Model,
		Username:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// Determine whether to update yourself or update others
	if userId == int(ctxUser.ID) {
		// If you are updating yourself
		// Cannot disable yourself
		if req.Status == 2 {
			response.Fail(c, nil, "Can't disable myself")
			return
		}
		// Cannot change own role
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			response.Fail(c, nil, "Can't change own role")
			return
		}

		// You cannot update your own password, you can only update it in the personal center
		if req.Password != "" {
			response.Fail(c, nil, "Please go to the personal center to update your password")
			return
		}

		// Password assignment
		user.Password = ctxUser.Password

	} else {
		// If updating someone else
		// Users cannot update users whose role level is higher than their own or of the same level.
		// Get minimum value of user role sorting based on userIdID in path
		minRoleSorts, err := uc.UserRepository.GetUserMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			response.Fail(c, nil, "Failed to obtain user role sorting minimum value based on user ID")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			response.Fail(c, nil, "Users cannot update users whose role level is higher than their own or of the same level.")
			return
		}

		// Users cannot update other users' role levels to be higher or equal to themselves.
		if currentRoleSortMin >= reqRoleSortMin {
			response.Fail(c, nil, "Users cannot update the role level of other users to be higher than or equal to themselves.")
			return
		}

		// Password assignment
		if req.Password != "" {
			// Password decrypted via RSA
			decodeData, err := util.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
			if err != nil {
				response.Fail(c, nil, err.Error())
				return
			}
			req.Password = string(decodeData)
			user.Password = util.GenPasswd(req.Password)
		}

	}

	// Update user
	err = uc.UserRepository.UpdateUser(&user)
	if err != nil {
		response.Fail(c, nil, "Update user failed: "+err.Error())
		return
	}
	response.Success(c, nil, "Update user successfully")

}

// Delete users in batches
func (uc UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req vo.DeleteUserRequest
	//Parameter binding
	if err := c.ShouldBind(&req); err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	// Parameter verification
	if err := common.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(common.Trans)
		response.Fail(c, nil, errStr)
		return
	}

	// User ID passed from the front end
	reqUserIds := req.UserIds
	// Get minimum value of user role sorting based on user ID
	roleMinSortList, err := uc.UserRepository.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		response.Fail(c, nil, "Failed to obtain user role sorting minimum value based on user ID")
		return
	}

	// The current user role sorting minimum value (the highest level role) and the current user
	minSort, ctxUser, err := uc.UserRepository.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// Cannot delete itself
	if funk.Contains(reqUserIds, ctxUser.ID) {
		response.Fail(c, nil, "Users cannot delete themselves")
		return
	}

	// Users whose roles are ranked lower (higher level) than their own cannot be deleted.
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			response.Fail(c, nil, "Users cannot delete users whose role level is higher than their own")
			return
		}
	}

	err = uc.UserRepository.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		response.Fail(c, nil, "Failed to delete user: "+err.Error())
		return
	}

	response.Success(c, nil, "Delete user successfully")

}
