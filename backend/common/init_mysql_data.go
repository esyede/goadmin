package common

import (
	"github.com/esyede/goadmin/backend/config"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/util"
	"errors"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

// Initialize mysql data
func InitData() {
	// Whether to initialize data
	if !config.Conf.System.InitData {
		return
	}

	// 1. Write character data
	newRoles := make([]*model.Role, 0)
	roles := []*model.Role{
		{
			Model:   gorm.Model{ID: 1},
			Name:    "Administrator",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "system",
		},
		{
			Model:   gorm.Model{ID: 2},
			Name:    "User",
			Keyword: "user",
			Desc:    new(string),
			Sort:    3,
			Status:  1,
			Creator: "system",
		},
		{
			Model:   gorm.Model{ID: 3},
			Name:    "Guest",
			Keyword: "guest",
			Desc:    new(string),
			Sort:    5,
			Status:  1,
			Creator: "system",
		},
	}

	for _, role := range roles {
		err := DB.First(&role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		err := DB.Create(&newRoles).Error
		if err != nil {
			Log.Errorf("Failed to write system role data: %v", err)
		}
	}

	// 2. write menu
	newMenus := make([]model.Menu, 0)
	var uint0 uint = 0
	var uint1 uint = 1
	componentStr := "component"
	systemUserStr := "/system/user"
	userStr := "user"
	peoplesStr := "peoples"
	treeTableStr := "tree-table"
	treeStr := "tree"
	exampleStr := "example"
	logOperationStr := "/log/operation-log"
	documentationStr := "documentation"
	var uint6 uint = 6
	menus := []model.Menu{
		{
			Model:     gorm.Model{ID: 1},
			Name:      "System",
			Title:     "System Management",
			Icon:      &componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  &systemUserStr,
			Sort:      10,
			ParentId:  &uint0,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 2},
			Name:      "User",
			Title:     "User Management",
			Icon:      &userStr,
			Path:      "user",
			Component: "/system/user/index",
			Sort:      11,
			ParentId:  &uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 3},
			Name:      "Role",
			Title:     "Role Management",
			Icon:      &peoplesStr,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      12,
			ParentId:  &uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 4},
			Name:      "Menu",
			Title:     "Menu Management",
			Icon:      &treeTableStr,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      13,
			ParentId:  &uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 5},
			Name:      "API",
			Title:     "API Management",
			Icon:      &treeStr,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentId:  &uint1,
			Roles:     roles[:1],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 6},
			Name:      "Log",
			Title:     "Log Management",
			Icon:      &exampleStr,
			Path:      "/log",
			Component: "Layout",
			Redirect:  &logOperationStr,
			Sort:      20,
			ParentId:  &uint0,
			Roles:     roles[:2],
			Creator:   "system",
		},
		{
			Model:     gorm.Model{ID: 7},
			Name:      "OperationLog",
			Title:     "OperationLog",
			Icon:      &documentationStr,
			Path:      "operation-log",
			Component: "/log/operation-log/index",
			Sort:      21,
			ParentId:  &uint6,
			Roles:     roles[:2],
			Creator:   "system",
		},
	}
	for _, menu := range menus {
		err := DB.First(&menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := DB.Create(&newMenus).Error
		if err != nil {
			Log.Errorf("Failed to write system menu data: %v", err)
		}
	}

	// 3. Write user
	newUsers := make([]model.User, 0)
	users := []model.User{
		{
			Model:        gorm.Model{ID: 1},
			Username:     "admin",
			Password:     util.GenPasswd("123456"),
			Mobile:       "081234567890",
			Avatar:       "https://i.pravatar.cc/300",
			Nickname:     new(string),
			Introduction: new(string),
			Status:       1,
			Creator:      "system",
			Roles:        roles[:1],
		},
		{
			Model:        gorm.Model{ID: 2},
			Username:     "faker",
			Password:     util.GenPasswd("123456"),
			Mobile:       "19999999999",
			Avatar:       "https://i.pravatar.cc/300",
			Nickname:     new(string),
			Introduction: new(string),
			Status:       1,
			Creator:      "system",
			Roles:        roles[:2],
		},
		{
			Model:        gorm.Model{ID: 3},
			Username:     "nike",
			Password:     util.GenPasswd("123456"),
			Mobile:       "081234567891",
			Avatar:       "https://i.pravatar.cc/300",
			Nickname:     new(string),
			Introduction: new(string),
			Status:       1,
			Creator:      "system",
			Roles:        roles[1:2],
		},
		{
			Model:        gorm.Model{ID: 4},
			Username:     "bob",
			Password:     util.GenPasswd("123456"),
			Mobile:       "081234567892",
			Avatar:       "https://i.pravatar.cc/300",
			Nickname:     new(string),
			Introduction: new(string),
			Status:       1,
			Creator:      "system",
			Roles:        roles[2:3],
		},
	}

	for _, user := range users {
		err := DB.First(&user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := DB.Create(&newUsers).Error
		if err != nil {
			Log.Errorf("Failed to write user data: %v", err)
		}
	}

	// 4. Write api
	apis := []model.Api{
		{
			Method:   "POST",
			Path:     "/base/login",
			Category: "base",
			Desc:     "User login",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/logout",
			Category: "base",
			Desc:     "User logout",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/base/refreshToken",
			Category: "base",
			Desc:     "Refresh JWT token",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/info",
			Category: "user",
			Desc:     "Get current logged in user information",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/user/list",
			Category: "user",
			Desc:     "Get user list",
			Creator:  "system",
		},
		{
			Method:   "PUT",
			Path:     "/user/changePwd",
			Category: "user",
			Desc:     "Update user login password",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/user/create",
			Category: "user",
			Desc:     "Create user",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/user/update/:userId",
			Category: "user",
			Desc:     "Update user",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/user/delete/batch",
			Category: "user",
			Desc:     "Delete user",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/list",
			Category: "role",
			Desc:     "Get role list",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/role/create",
			Category: "role",
			Desc:     "Create role",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/role/update/:roleId",
			Category: "role",
			Desc:     "Update role",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/menus/get/:roleId",
			Category: "role",
			Desc:     "Get the role's permissions menu",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/role/menus/update/:roleId",
			Category: "role",
			Desc:     "Update the role's permissions menu",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/role/apis/get/:roleId",
			Category: "role",
			Desc:     "Obtain the permission interface of the role",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/role/apis/update/:roleId",
			Category: "role",
			Desc:     "Update role permission interface",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/role/delete/batch",
			Category: "role",
			Desc:     "Delete roles in batches",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/list",
			Category: "menu",
			Desc:     "Get menu list",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/tree",
			Category: "menu",
			Desc:     "Get menu tree",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/menu/create",
			Category: "menu",
			Desc:     "Create menu",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/menu/update/:menuId",
			Category: "menu",
			Desc:     "Update menu",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/menu/delete/batch",
			Category: "menu",
			Desc:     "Delete menu",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/list/:userId",
			Category: "menu",
			Desc:     "Get the user's list of accessible menus",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/menu/access/tree/:userId",
			Category: "menu",
			Desc:     "Get the user's accessible menu tree",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/api/list",
			Category: "api",
			Desc:     "Get interface list",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/api/tree",
			Category: "api",
			Desc:     "Get interface tree",
			Creator:  "system",
		},
		{
			Method:   "POST",
			Path:     "/api/create",
			Category: "api",
			Desc:     "Create interface",
			Creator:  "system",
		},
		{
			Method:   "PATCH",
			Path:     "/api/update/:roleId",
			Category: "api",
			Desc:     "Update interface",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/api/delete/batch",
			Category: "api",
			Desc:     "Batch delete interface",
			Creator:  "system",
		},
		{
			Method:   "GET",
			Path:     "/log/operation/list",
			Category: "log",
			Desc:     "Get operation log list",
			Creator:  "system",
		},
		{
			Method:   "DELETE",
			Path:     "/log/operation/delete/batch",
			Category: "log",
			Desc:     "Delete operation logs in batches",
			Creator:  "system",
		},
	}
	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	for i, api := range apis {
		api.ID = uint(i + 1)
		err := DB.First(&api, api.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApi = append(newApi, api)

			// Administrators have all API permissions
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})

			// Non-administrators have basic permissions
			basePaths := []string{
				"/base/login",
				"/base/logout",
				"/base/refreshToken",
				"/user/info",
				"/menu/access/tree/:userId",
			}

			if funk.ContainsString(basePaths, api.Path) {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[2].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	if len(newApi) > 0 {
		if err := DB.Create(&newApi).Error; err != nil {
			Log.Errorf("Failed to write api data: %v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			Log.Errorf("Failed to write casbin data: %v", err)
		}
	}
}
