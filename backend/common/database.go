package common

import (
	"backend/config"
	"backend/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Global mysql database variables
var DB *gorm.DB

// Initialize mysql database
func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)
	// Hide password
	showDsn := fmt.Sprintf(
		"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		config.Conf.Mysql.Username,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
		config.Conf.Mysql.Charset,
		config.Conf.Mysql.Collation,
		config.Conf.Mysql.Query,
	)
	// Log.Info("Database connection DSN: ", showDsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Disable foreign keys (no real foreign key constraints will be created in mysql when specifying foreign keys)
		DisableForeignKeyConstraintWhenMigrating: true,
		// Enable table prefix
		// NamingStrategy: schema.NamingStrategy{
		// 	TablePrefix: config.Conf.Mysql.TablePrefix + "_",
		// },
	})
	if err != nil {
		Log.Panicf("Failed initializing mysql database: %v", err)
		panic(fmt.Errorf("Failed initializing mysql database: %v", err))
	}

	// Enable mysql log
	if config.Conf.Mysql.LogMode {
		db.Debug()
	}
	// Global DB assignment
	DB = db
	// Automatically migrate table structure
	dbAutoMigrate()
	Log.Infof("初始化mysql数据库完成! dsn: %s", showDsn)
}

// Automatically migrate table structure
func dbAutoMigrate() {
	DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Menu{},
		&model.Api{},
		&model.OperationLog{},
	)
}
