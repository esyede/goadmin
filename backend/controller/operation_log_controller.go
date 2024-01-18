package controller

import (
	"backend/common"
	"backend/repository"
	"backend/response"
	"backend/vo"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type IOperationLogController interface {
	GetOperationLogs(c *gin.Context)             // Get the operation log list
	BatchDeleteOperationLogByIds(c *gin.Context) // Batch delete operation logs
}

type OperationLogController struct {
	operationLogRepository repository.IOperationLogRepository
}

func NewOperationLogController() IOperationLogController {
	operationLogRepository := repository.NewOperationLogRepository()
	operationLogController := OperationLogController{operationLogRepository: operationLogRepository}
	return operationLogController
}

// Get the operation log list
func (oc OperationLogController) GetOperationLogs(c *gin.Context) {
	var req vo.OperationLogListRequest
	// Bind parameters
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
	logs, total, err := oc.operationLogRepository.GetOperationLogs(&req)
	if err != nil {
		response.Fail(c, nil, "Failed to obtain operation log list: "+err.Error())
		return
	}
	response.Success(c, gin.H{"logs": logs, "total": total}, "Obtaining operation log list successfully")
}

// Batch delete operation logs
func (oc OperationLogController) BatchDeleteOperationLogByIds(c *gin.Context) {
	var req vo.DeleteOperationLogRequest
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
	err := oc.operationLogRepository.BatchDeleteOperationLogByIds(req.OperationLogIds)
	if err != nil {
		response.Fail(c, nil, "Failed to delete log: "+err.Error())
		return
	}

	response.Success(c, nil, "Log deleted successfully")
}
