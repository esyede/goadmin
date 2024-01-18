package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name    string  `gorm:"type:varchar(20);not null;unique" json:"name"`
	Keyword string  `gorm:"type:varchar(20);not null;unique" json:"keyword"`
	Desc    *string `gorm:"type:varchar(100);" json:"desc"`
	Status  uint    `gorm:"type:tinyint(1);default:1;comment:'1 normal, 2 disabled'" json:"status"`
	Sort    uint    `gorm:"type:int(3);default:999;comment:'Role sorting (greater value means lower permissions. Value of 1 indicates a superadmin)'" json:"sort"`
	Creator string  `gorm:"type:varchar(20);" json:"creator"`
	Users   []*User `gorm:"many2many:user_roles" json:"users"`
	Menus   []*Menu `gorm:"many2many:role_menus;" json:"menus"` // Role menu many-to-many relationship
}
