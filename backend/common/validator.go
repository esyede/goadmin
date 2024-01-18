package common

import (
	"regexp"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

// Global Validate data verification column
var Validate *validator.Validate

// Global translator
var Trans ut.Translator

// Initialize Validator data verification
func InitValidate() {
	indonesian := id.New()
	uni := ut.New(indonesian, indonesian)
	trans, _ := uni.GetTranslator("id")
	Trans = trans
	Validate = validator.New()
	_ = id_translations.RegisterDefaultTranslations(Validate, Trans)
	_ = Validate.RegisterValidation("checkMobile", checkMobile)
	Log.Infof("Initialization of validator.v10 data validator completed")
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}
