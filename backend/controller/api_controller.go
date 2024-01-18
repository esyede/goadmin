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

type IApiController interface {
	GetApis(c *gin.Context)             // Get interface list
	GetApiTree(c *gin.Context)          // Get interface tree (classified by interface Category field)
	CreateApi(c *gin.Context)           // Create interface
	UpdateApiById(c *gin.Context)       // Update interface
	BatchDeleteApiByIds(c *gin.Context) // Batch deletion interface
}

type ApiController struct {
	ApiRepository repository.IApiRepository
}

func NewApiController() IApiController {
	apiRepository := repository.NewApiRepository()
	apiController := ApiController{ApiRepository: apiRepository}
	return apiController
}

// Get interface list
func (ac ApiController) GetApis(c *gin.Context) {
	var req vo.ApiListRequest
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
	apis, total, err := ac.ApiRepository.GetApis(&req)
	if err != nil {
		response.Fail(c, nil, "Failed to get interface list")
		return
	}
	response.Success(c, gin.H{
		"apis": apis, "total": total,
	}, "Obtaining interface list successfully")
}

// Get interface tree (classified by interface Category field)
func (ac ApiController) GetApiTree(c *gin.Context) {
	tree, err := ac.ApiRepository.GetApiTree()
	if err != nil {
		response.Fail(c, nil, "Failed to obtain interface tree")
		return
	}
	response.Success(c, gin.H{
		"apiTree": tree,
	}, "Obtaining the interface tree successfully")
}

// Create interface
func (ac ApiController) CreateApi(c *gin.Context) {
	var req vo.CreateApiRequest
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

	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	// Create interface
	err = ac.ApiRepository.CreateApi(&api)
	if err != nil {
		response.Fail(c, nil, "Failed to create interface: "+err.Error())
		return
	}

	response.Success(c, nil, "Interface created successfully")
	return
}

// Update interface
func (ac ApiController) UpdateApiById(c *gin.Context) {
	var req vo.UpdateApiRequest
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

	// Get apiId in the path
	apiId, _ := strconv.Atoi(c.Param("apiId"))
	if apiId <= 0 {
		response.Fail(c, nil, "Interface ID is incorrect")
		return
	}

	// Get current user
	ur := repository.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain current user information")
		return
	}

	api := model.Api{
		Method:   req.Method,
		Path:     req.Path,
		Category: req.Category,
		Desc:     req.Desc,
		Creator:  ctxUser.Username,
	}

	err = ac.ApiRepository.UpdateApiById(uint(apiId), &api)
	if err != nil {
		response.Fail(c, nil, "Update interface failed: "+err.Error())
		return
	}

	response.Success(c, nil, "Update interface successful")
}

// Batch deletion interface
func (ac ApiController) BatchDeleteApiByIds(c *gin.Context) {
	var req vo.DeleteApiRequest
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

	// Delete interface
	err := ac.ApiRepository.BatchDeleteApiByIds(req.ApiIds)
	if err != nil {
		response.Fail(c, nil, "Failed to delete interface: "+err.Error())
		return
	}

	response.Success(c, nil, "Interface deleted successfully")
}
