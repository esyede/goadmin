package repository

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"

	"github.com/thoas/go-funk"
)

type IMenuRepository interface {
	GetMenus() ([]*model.Menu, error)                           // Get menu list
	GetMenuTree() ([]*model.Menu, error)                        // Get menu tree
	CreateMenu(menu *model.Menu) error                          // Create menu
	UpdateMenuById(menuId uint, menu *model.Menu) error         // Update menu
	BatchDeleteMenuByIds(menuIds []uint) error                  // Batch delete menu
	GetUserMenusByUserId(userId uint) ([]*model.Menu, error)    // Get user's permission (accessible) menu list based on the user ID
	GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) // Get user's permissions (accessible) menu tree based on the user ID
}

type MenuRepository struct {
}

func NewMenuRepository() IMenuRepository {
	return MenuRepository{}
}

// Get menu list
func (m MenuRepository) GetMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// Get menu tree
func (m MenuRepository) GetMenuTree() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := common.DB.Order("sort").Find(&menus).Error
	// The one with parentId 0 is the root menu
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if *m.ParentId == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// Create menu
func (m MenuRepository) CreateMenu(menu *model.Menu) error {
	err := common.DB.Create(menu).Error
	return err
}

// Update menu
func (m MenuRepository) UpdateMenuById(menuId uint, menu *model.Menu) error {
	err := common.DB.Model(menu).Where("id = ?", menuId).Updates(menu).Error
	return err
}

// Batch delete menu
func (m MenuRepository) BatchDeleteMenuByIds(menuIds []uint) error {
	var menus []*model.Menu
	err := common.DB.Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = common.DB.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// Get user's permission (accessible) menu list based on the user ID
func (m MenuRepository) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	// Get user
	var user model.User
	err := common.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// Get role
	roles := user.Roles
	// Menu collection of all characters
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		err := common.DB.Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// Get character's menu
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// The menu collection of all characters is deduplicated
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// Get menu with status 1
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// Get user's permissions (accessible) menu tree based on the user ID
func (m MenuRepository) GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) {
	menus, err := m.GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}
