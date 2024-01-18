package repository

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/vo"
	"errors"
	"fmt"
	"strings"
)

type IRoleRepository interface {
	GetRoles(req *vo.RoleListRequest) ([]model.Role, int64, error)       // Get role list
	GetRolesByIds(roleIds []uint) ([]*model.Role, error)                 // Get role based on the role ID
	CreateRole(role *model.Role) error                                   // Creating a Role
	UpdateRoleById(roleId uint, role *model.Role) error                  // Update role
	GetRoleMenusById(roleId uint) ([]*model.Menu, error)                 // Get role's permission menu
	UpdateRoleMenus(role *model.Role) error                              // Update the role's permissions menu
	GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error)   // Get permission interface of the role based on the role keyword
	UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error // Update the permission interface of the role (delete them all first and then add them)
	BatchDeleteRoleByIds(roleIds []uint) error                           // Delete role
}

type RoleRepository struct {
}

func NewRoleRepository() IRoleRepository {
	return RoleRepository{}
}

// Get role list
func (r RoleRepository) GetRoles(req *vo.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := common.DB.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
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
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// Get role based on the role ID
func (r RoleRepository) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// Creating a Role
func (r RoleRepository) CreateRole(role *model.Role) error {
	err := common.DB.Create(role).Error
	return err
}

// Update role
func (r RoleRepository) UpdateRoleById(roleId uint, role *model.Role) error {
	err := common.DB.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// Get role's permission menu
func (r RoleRepository) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := common.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// Update the role's permissions menu
func (r RoleRepository) UpdateRoleMenus(role *model.Role) error {
	err := common.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// Get permission interface of the role based on the role keyword
func (r RoleRepository) GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error) {
	policies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)

	// Get all interfaces
	var apis []*model.Api
	err := common.DB.Find(&apis).Error
	if err != nil {
		return apis, errors.New("Failed to obtain role permission interface")
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// Update the permission interface of the role (delete them all first and then add them)
func (r RoleRepository) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// First obtain the existing police of the role corresponding to the role ID in path (need to be deleted first)
	err := common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("The role's permission interface policy failed to load.")
	}
	rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if len(rmPolicies) > 0 {
		isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("Failed to update role's permission interface")
		}
	}
	isAdded, _ := common.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("Failed to update role's permission interface")
	}
	err = common.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("The role's permission interface was updated successfully, but the role's permission interface policy failed to load.")
	} else {
		return err
	}
}

// Delete role
func (r RoleRepository) BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*model.Role
	err := common.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// Delete the casbin policy if the deletion is successful.
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := common.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if len(rmPolicies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("Deleting the role was successful, but deleting the role-associated permission interface failed.")
				}
			}
		}

	}
	return err
}
