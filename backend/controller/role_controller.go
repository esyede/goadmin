package controller

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/repository"
	"github.com/esyede/goadmin/backend/response"
	"github.com/esyede/goadmin/backend/vo"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
)

type IRoleController interface {
	GetRoles(c *gin.Context)             // Get role list
	CreateRole(c *gin.Context)           // Creating a Role
	UpdateRoleById(c *gin.Context)       // Update role
	GetRoleMenusById(c *gin.Context)     // Get role's permission menu
	UpdateRoleMenusById(c *gin.Context)  // Update the role's permissions menu
	GetRoleApisById(c *gin.Context)      // Get permission interface of the role
	UpdateRoleApisById(c *gin.Context)   // Update the permission interface of the role
	BatchDeleteRoleByIds(c *gin.Context) // Delete roles in batches
}

type RoleController struct {
	RoleRepository repository.IRoleRepository
}

func NewRoleController() IRoleController {
	roleRepository := repository.NewRoleRepository()
	roleController := RoleController{RoleRepository: roleRepository}
	return roleController
}

// Get role list
func (rc RoleController) GetRoles(c *gin.Context) {
	var req vo.RoleListRequest
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

	// Get role list
	roles, total, err := rc.RoleRepository.GetRoles(&req)
	if err != nil {
		response.Fail(c, nil, "Failed to get role list: "+err.Error())
		return
	}
	response.Success(c, gin.H{"roles": roles, "total": total}, "Obtaining role list successfully")
}

// Creating a Role
func (rc RoleController) CreateRole(c *gin.Context) {
	var req vo.CreateRoleRequest
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

	// Get current user's highest role level
	uc := repository.NewUserRepository()
	sort, ctxUser, err := uc.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain current user's highest role level: "+err.Error())
		return
	}

	// Users cannot create characters with a higher level or the same level as themselves
	if sort >= req.Sort {
		response.Fail(c, nil, "You cannot create a role with a higher level or the same level as yourself.")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// Creating a Role
	err = rc.RoleRepository.CreateRole(&role)
	if err != nil {
		response.Fail(c, nil, "Failed to create role: "+err.Error())
		return
	}
	response.Success(c, nil, "Role created successfully")

}

// Update role
func (rc RoleController) UpdateRoleById(c *gin.Context) {
	var req vo.CreateRoleRequest
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
	// Get roleId in path
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "Incorrect role ID")
		return
	}

	//The current user role sorting minimum value (the highest level role) and the current user
	ur := repository.NewUserRepository()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// Cannot update characters that are higher or equal to your own character level
	// Get role information based on the role ID in path
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}
	if minSort >= roles[0].Sort {
		response.Fail(c, nil, "You cannot update role that are higher or equal to your own role level.")
		return
	}

	// Cannot update the character level to be higher than the current user's level
	if minSort >= req.Sort {
		response.Fail(c, nil, "The role level cannot be updated to be higher than or the same as the current user's level")
		return
	}

	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}

	// Update role
	err = rc.RoleRepository.UpdateRoleById(uint(roleId), &role)
	if err != nil {
		response.Fail(c, nil, "Failed to update role: "+err.Error())
		return
	}

	// If the update is successful and the role's keyword is updated, update the policy in casbin
	if req.Keyword != roles[0].Keyword {
		// Get policy
		rolePolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if len(rolePolicies) == 0 {
			response.Success(c, nil, "Update role successfully")
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		// Replace keyword
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}

		// gormadapter does not implement the UpdatePolicies method, wait for gorm to update
		// isUpdated, _ := common.CasbinEnforcer.UpdatePolicies(rolePoliciesCopy, rolePolicies)
		// if !isUpdated {
		// 	response.Fail(c, nil, "The role was updated, but the permission interface associated with the role keyword failed to be updated.")
		// 	return
		// }

		// Here you need to add first and then delete (deleting first and then adding will cause an error)
		isAdded, _ := common.CasbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			response.Fail(c, nil, "The role was updated, but the permission interface associated with the role keyword failed to be updated.")
			return
		}
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			response.Fail(c, nil, "The role was updated, but the permission interface associated with the role keyword failed to be updated.")
			return
		}
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			response.Fail(c, nil, "The role was updated, but the permission interface policy of the role associated with the role keyword failed to load.")
			return
		}

	}

	// There are two ways to update the role to successfully process the user information cache: (The second method is used here because there may be many users under one role, and the second method can spread the pressure on the database)
	// 1. It can help users update the information cache of users with this role, use the following method
	// err = ur.UpdateUserInfoCacheByRoleId(uint(roleId))
	// 2. Directly clear the cache and let active users re-cache the latest user information.
	ur.ClearUserInfoCache()

	response.Success(c, nil, "Update role successfully")
}

// Get role's permission menu
func (rc RoleController) GetRoleMenusById(c *gin.Context) {
	// Get roleId in path
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "Incorrect role ID")
		return
	}
	menus, err := rc.RoleRepository.GetRoleMenusById(uint(roleId))
	if err != nil {
		response.Fail(c, nil, "Failed to obtain role's permission menu: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, "Obtaining the role's permission menu successfully")
}

// Update the role's permissions menu
func (rc RoleController) UpdateRoleMenusById(c *gin.Context) {
	var req vo.UpdateRoleMenusRequest
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
	// Get roleId in path
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "Incorrect role ID")
		return
	}
	// Get role information based on the role ID in path
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}

	// The current user role sorting minimum value (the highest level role) and the current user
	ur := repository.NewUserRepository()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// Non-administrators cannot update the permission menu of a role higher than or equal to their own role level.
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			response.Fail(c, nil, "You cannot update the permissions menu of a role that is higher than or equal to your own role level.")
			return
		}
	}

	// Get permission menu owned by the current user
	mr := repository.NewMenuRepository()
	ctxUserMenus, err := mr.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		response.Fail(c, nil, "Failed to get list of accessible menus for the current user: "+err.Error())
		return
	}

	// Get permission menu ID owned by the current user
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// The latest MenuIds collection is transmitted from the front end
	menuIds := req.MenuIds

	// The menu collection that the user needs to modify
	reqMenus := make([]*model.Menu, 0)

	// Non-administrators cannot set the role's permission menu to be more than the permission menu owned by the current user.
	if minSort != 1 {
		for _, id := range menuIds {
			if !funk.Contains(ctxUserMenusIds, id) {
				response.Fail(c, nil, fmt.Sprintf("Do not have permission to set menu with ID %d", id))
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// The administrator can set it at will
		// Query the menu based on menuIds
		menus, err := mr.GetMenus()
		if err != nil {
			response.Fail(c, nil, "Failed to get menu list: "+err.Error())
			return
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = rc.RoleRepository.UpdateRoleMenus(roles[0])
	if err != nil {
		response.Fail(c, nil, "Failed to update role's permissions menu: "+err.Error())
		return
	}

	response.Success(c, nil, "Updated role's permissions menu successfully")

}

// Get permission interface of the role
func (rc RoleController) GetRoleApisById(c *gin.Context) {
	// Get roleId in path
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "Incorrect role ID")
		return
	}
	// Get role information based on the role ID in path
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}
	// Get policy in casbin based on the role keyword
	keyword := roles[0].Keyword
	apis, err := rc.RoleRepository.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	response.Success(c, gin.H{"apis": apis}, "Obtaining the role's permission interface successfully")
}

// Update the permission interface of the role
func (rc RoleController) UpdateRoleApisById(c *gin.Context) {
	var req vo.UpdateRoleApisRequest
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

	// Get roleId in path
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		response.Fail(c, nil, "Incorrect role ID")
		return
	}
	// Get role information based on the role ID in path
	roles, err := rc.RoleRepository.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}

	// The current user role sorting minimum value (the highest level role) and the current user
	ur := repository.NewUserRepository()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// Non-administrators cannot update permission interfaces with roles higher than or equal to their own role level
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			response.Fail(c, nil, "You cannot update the permission interface of a role that is higher than or equal to your own role level.")
			return
		}
	}

	// Get permission interface owned by the current user
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy := common.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	// Get set of permission interfaces that can be set by the role corresponding to the role ID in path
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}

	// The latest ApiID collection is transmitted from the front end
	apiIds := req.ApiIds
	// Get interface details based on apiID
	ar := repository.NewApiRepository()
	apis, err := ar.GetApisById(apiIds)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain interface information based on interface ID")
		return
	}
	// Generate role policies that the front end wants to set
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}

	// Non-administrators cannot set the role's permission interface to be more than the permission interface owned by the current user.
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				response.Fail(c, nil, fmt.Sprintf("Do not have permission to set the interface with path %s and request method %s", reqPolicy[1], reqPolicy[2]))
				return
			}
		}
	}

	// Update the permission interface of the role
	err = rc.RoleRepository.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	response.Success(c, nil, "Update role's permission interface successfully")

}

// Delete roles in batches
func (rc RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req vo.DeleteRoleRequest
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

	// Get highest level role of the current user
	ur := repository.NewUserRepository()
	minSort, _, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	// The front end passes the character ID that needs to be deleted
	roleIds := req.RoleIds
	// Get role information
	roles, err := rc.RoleRepository.GetRolesByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain role information: "+err.Error())
		return
	}
	if len(roles) == 0 {
		response.Fail(c, nil, "No role information was obtained")
		return
	}

	// Characters that are higher or equal to your own character level cannot be deleted.
	for _, role := range roles {
		if minSort >= role.Sort {
			response.Fail(c, nil, "You cannot delete roles that are higher or equal to your own role level.")
			return
		}
	}

	// Delete role
	err = rc.RoleRepository.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		response.Fail(c, nil, "Failed to delete role")
		return
	}

	// If the role is successfully deleted, the cache will be cleared directly, allowing active users to re-cache the latest user information.
	ur.ClearUserInfoCache()
	response.Success(c, nil, "Role deleted successfully")
}
