package model

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	Name       string  `gorm:"type:varchar(50);comment:'Menu name'" json:"name"`
	Title      string  `gorm:"type:varchar(50);comment:'Menu title'" json:"title"`
	Icon       *string `gorm:"type:varchar(50);comment:'Icon'" json:"icon"`
	Path       string  `gorm:"type:varchar(100);comment:'Menu access path'" json:"path"`
	Redirect   *string `gorm:"type:varchar(100);comment:'Redirect path'" json:"redirect"`
	Component  string  `gorm:"type:varchar(100);comment:'Frontend component path'" json:"component"`
	Sort       uint    `gorm:"type:int(3) unsigned;default:999;comment:'Menu order (1-999)'" json:"sort"`
	Status     uint    `gorm:"type:tinyint(1);default:1;comment:'Menu status (normal/disabled, default: normal)'" json:"status"`
	Hidden     uint    `gorm:"type:tinyint(1);default:2;comment:'Hide menu from sidebar (1 hide, 2 show)'" json:"hidden"`
	NoCache    uint    `gorm:"type:tinyint(1);default:2;comment:'Whether the menu is cached by <keep-alive> (1 not-cached, 2 cached)'" json:"noCache"`
	AlwaysShow uint    `gorm:"type:tinyint(1);default:2;comment:'Ignore the previously defined rules and always display the root route (1 ignore, 2 dont ignore)'" json:"alwaysShow"`
	Breadcrumb uint    `gorm:"type:tinyint(1);default:1;comment:'Breadcrumb visibility (visible/hidden, visible by default)'" json:"breadcrumb"`
	ActiveMenu *string `gorm:"type:varchar(100);comment:'When using other routes, highlight the route in the sidebar'" json:"activeMenu"`
	ParentId   *uint   `gorm:"default:0;comment:'Parent menu number (0 means the root menu)'" json:"parentId"`
	Creator    string  `gorm:"type:varchar(20);comment:'Creator'" json:"creator"`
	Children   []*Menu `gorm:"-" json:"children"`                  // Submenu collection
	Roles      []*Role `gorm:"many2many:role_menus;" json:"roles"` // Role menu many-to-many relationship
}
