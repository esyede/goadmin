package controller

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/repository"
	"github.com/esyede/goadmin/backend/response"
	"github.com/esyede/goadmin/backend/vo"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IMenuController interface {
	GetMenus(c *gin.Context)                // Get menu list
	GetMenuTree(c *gin.Context)             // Get menu tree
	CreateMenu(c *gin.Context)              // Create menu
	UpdateMenuById(c *gin.Context)          // Update menu
	BatchDeleteMenuByIds(c *gin.Context)    // Batch delete menu
	GetUserMenusByUserId(c *gin.Context)    // Get user's accessible menu list
	GetUserMenuTreeByUserId(c *gin.Context) // Get user's accessible menu tree
}

type MenuController struct {
	MenuRepository repository.IMenuRepository
}

func NewMenuController() IMenuController {
	menuRepository := repository.NewMenuRepository()
	menuController := MenuController{MenuRepository: menuRepository}
	return menuController
}

// Get menu list
func (mc MenuController) GetMenus(c *gin.Context) {
	menus, err := mc.MenuRepository.GetMenus()
	if err != nil {
		response.Fail(c, nil, "Failed to get menu list: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, "Get menu list successfully")
}

// Get menu tree
func (mc MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.MenuRepository.GetMenuTree()
	if err != nil {
		response.Fail(c, nil, "Failed to get menu tree: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, "Get menu tree successfully")
}

// Create menu
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req vo.CreateMenuRequest
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

	// Get current user
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain current user information")
		return
	}

	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentId:   &req.ParentId,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.CreateMenu(&menu)
	if err != nil {
		response.Fail(c, nil, "Failed to create menu: "+err.Error())
		return
	}
	response.Success(c, nil, "Menu created successfully")
}

// Update menu
func (mc MenuController) UpdateMenuById(c *gin.Context) {
	var req vo.UpdateMenuRequest
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

	// Get the menuId in the path
	menuId, _ := strconv.Atoi(c.Param("menuId"))
	if menuId <= 0 {
		response.Fail(c, nil, "Menu ID is incorrect")
		return
	}

	// Get the current user
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain current user information")
		return
	}

	menu := model.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentId:   &req.ParentId,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.UpdateMenuById(uint(menuId), &menu)
	if err != nil {
		response.Fail(c, nil, "Update menu failed: "+err.Error())
		return
	}

	response.Success(c, nil, "Update menu successfully")

}

// Batch delete menu
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req vo.DeleteMenuRequest
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
	err := mc.MenuRepository.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		response.Fail(c, nil, "Failed to delete menu: "+err.Error())
		return
	}

	response.Success(c, nil, "Delete menu successfully")
}

// Get the user's accessible menu list based on the user ID
func (mc MenuController) GetUserMenusByUserId(c *gin.Context) {
	// Get the userId in the path
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		response.Fail(c, nil, "User ID is incorrect")
		return
	}

	menus, err := mc.MenuRepository.GetUserMenusByUserId(uint(userId))
	if err != nil {
		response.Fail(c, nil, "Failed to get list of user's accessible menus: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menus": menus}, "Obtaining the user's accessible menu list successfully")
}

// Get the user's accessible menu tree based on the user ID
func (mc MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
	// Get the userId in the path
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		response.Fail(c, nil, "User ID is incorrect")
		return
	}

	menuTree, err := mc.MenuRepository.GetUserMenuTreeByUserId(uint(userId))
	if err != nil {
		response.Fail(c, nil, "Failed to get user's accessible menu tree: "+err.Error())
		return
	}
	response.Success(c, gin.H{"menuTree": menuTree}, "Obtaining the user's accessible menu tree successfully")
}
