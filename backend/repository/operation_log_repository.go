package repository

import (
	"github.com/esyede/goadmin/backend/common"
	"github.com/esyede/goadmin/backend/model"
	"github.com/esyede/goadmin/backend/vo"
	"fmt"
	"strings"
)

type IOperationLogRepository interface {
	GetOperationLogs(req *vo.OperationLogListRequest) ([]model.OperationLog, int64, error)
	BatchDeleteOperationLogByIds(ids []uint) error
	SaveOperationLogChannel(olc <-chan *model.OperationLog) // Handle OperationLogChan to log to database
}

type OperationLogRepository struct {
}

func NewOperationLogRepository() IOperationLogRepository {
	return OperationLogRepository{}
}

func (o OperationLogRepository) GetOperationLogs(req *vo.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	db := common.DB.Model(&model.OperationLog{}).Order("start_time DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// Paging
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err

}

func (o OperationLogRepository) BatchDeleteOperationLogByIds(ids []uint) error {
	err := common.DB.Where("id IN (?)", ids).Unscoped().Delete(&model.OperationLog{}).Error
	return err
}

// var Logs []model.OperationLog // Global variables need to be locked by multiple threads, so each thread maintains one by itself.
// Handle OperationLogChan to log to database
func (o OperationLogRepository) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	// Will only be executed once when the thread is started
	Logs := make([]model.OperationLog, 0)

	// Always execute - it will execute after receiving olc
	for log := range olc {
		Logs = append(Logs, *log)
		// Every 10 records are sent to the database
		if len(Logs) > 5 {
			common.DB.Create(&Logs)
			Logs = make([]model.OperationLog, 0)
		}
	}
}
