package repository

import (
	"backend/common"
	"backend/dto"
	"backend/model"
	"backend/vo"
	"errors"
	"fmt"
	"strings"

	"github.com/thoas/go-funk"
)

type IApiRepository interface {
	GetApis(req *vo.ApiListRequest) ([]*model.Api, int64, error) // Get interface list
	GetApisById(apiIds []uint) ([]*model.Api, error)             // Get interface list based on the interface ID
	GetApiTree() ([]*dto.ApiTreeDto, error)                      // Get interface tree (classified by interface Category field)
	CreateApi(api *model.Api) error                              // Create interface
	UpdateApiById(apiId uint, api *model.Api) error              // Update interface
	BatchDeleteApiByIds(apiIds []uint) error                     // Batch deletion interface
	GetApiDescByPath(path string, method string) (string, error) // Get interface description based on the interface path and request method
}

type ApiRepository struct {
}

func NewApiRepository() IApiRepository {
	return ApiRepository{}
}

// Get interface list
func (a ApiRepository) GetApis(req *vo.ApiListRequest) ([]*model.Api, int64, error) {
	var list []*model.Api
	db := common.DB.Model(&model.Api{}).Order("created_at DESC")

	method := strings.TrimSpace(req.Method)

	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}

	path := strings.TrimSpace(req.Path)

	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}

	category := strings.TrimSpace(req.Category)

	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}

	creator := strings.TrimSpace(req.Creator)

	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
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

// Get interface list based on the interface ID
func (a ApiRepository) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := common.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

// Get interface tree (classified by interface Category field)
func (a ApiRepository) GetApiTree() ([]*dto.ApiTreeDto, error) {
	var apiList []*model.Api
	err := common.DB.Order("category").Order("created_at").Find(&apiList).Error

	// Get all categories
	var categoryList []string

	for _, api := range apiList {
		categoryList = append(categoryList, api.Category)
	}
	// Get classification after deduplication
	categoryUniq := funk.UniqString(categoryList)
	apiTree := make([]*dto.ApiTreeDto, len(categoryUniq))

	for i, category := range categoryUniq {
		apiTree[i] = &dto.ApiTreeDto{
			ID:       -i,
			Desc:     category,
			Category: category,
			Children: nil,
		}

		for _, api := range apiList {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, api)
			}
		}
	}

	return apiTree, err
}

// Create interface
func (a ApiRepository) CreateApi(api *model.Api) error {
	err := common.DB.Create(api).Error
	return err
}

// Update interface
func (a ApiRepository) UpdateApiById(apiId uint, api *model.Api) error {
	// Get interface information based on id
	var oldApi model.Api
	err := common.DB.First(&oldApi, apiId).Error

	if err != nil {
		return errors.New("Failed to obtain interface information based on interface ID")
	}

	err = common.DB.Model(api).Where("id = ?", apiId).Updates(api).Error

	if err != nil {
		return err
	}

	// After updating the method and path, update the policy in casbin.
	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies := common.CasbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		// The interface can only be operated if it exists in the policy of casbin.
		if len(policies) > 0 {
			// Delete first
			isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("Update permission interface failed")
			}

			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}

			// Add
			isAdded, _ := common.CasbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("Update permission interface failed")
			}
			// Load policy
			err := common.CasbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("The permission interface was updated successfully, but the permission interface policy loading failed.")
			} else {
				return err
			}
		}
	}
	return err
}

// Batch deletion interface
func (a ApiRepository) BatchDeleteApiByIds(apiIds []uint) error {

	apis, err := a.GetApisById(apiIds)
	if err != nil {
		return errors.New("Failed to obtain interface list based on interface ID")
	}

	if len(apis) == 0 {
		return errors.New("The interface list was not obtained based on the interface ID.")
	}

	err = common.DB.Where("id IN (?)", apiIds).Unscoped().Delete(&model.Api{}).Error
	// If the deletion is successful, delete the policy in casbin
	if err == nil {
		for _, api := range apis {
			policies := common.CasbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if len(policies) > 0 {
				isRemoved, _ := common.CasbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("Failed to delete permission interface")
				}
			}
		}
		// Reload strategy
		err := common.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("The permission interface was deleted successfully, but the permission interface policy failed to load.")
		} else {
			return err
		}
	}
	return err
}

// Get interface description based on the interface path and request method
func (a ApiRepository) GetApiDescByPath(path string, method string) (string, error) {
	var api model.Api
	err := common.DB.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
