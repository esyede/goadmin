package common

import (
	"backend/config"
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// Global CasbinEnforcer
var CasbinEnforcer *casbin.Enforcer

// Initialize casbin policy manager
func InitCasbinEnforcer() {
	e, err := mysqlCasbin()
	if err != nil {
		Log.Panicf("Failed to initialize Casbin: %v", err)
		panic(fmt.Sprintf("Failed to initialize Casbin: %v", err))
	}

	CasbinEnforcer = e
	Log.Info("Initialization of Casbin completed!")
}

func mysqlCasbin() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(DB)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(config.Conf.Casbin.ModelPath, a)
	if err != nil {
		return nil, err
	}

	err = e.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return e, nil
}
