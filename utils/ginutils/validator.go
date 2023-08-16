package ginutils

import (
	"github.com/Conflux-Chain/go-conflux-sdk/types/cfxaddress"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var cfxAddress validator.Func = func(fl validator.FieldLevel) bool {
	addr, ok := fl.Field().Interface().(string)
	if ok {
		if _, err := cfxaddress.NewFromBase32(addr); err == nil {
			return true
		}
	}
	return false
}

func RegisterValidation() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("cfxaddress", cfxAddress)
	}
}
