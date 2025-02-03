package validator

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shahzodshafizod/gocloud/pkg"
)

// alternative lib: https://github.com/asaskevich/govalidator

func OptionalDateOnly() pkg.Validator {
	return &validate{
		tag: "dateonly",
		fnc: func(ctx context.Context, fl validator.FieldLevel) bool {
			value := fl.Field().String()
			_, err := time.Parse(time.DateOnly, value)
			return err == nil
		},
	}
}

type validate struct {
	tag string
	fnc validator.FuncCtx
}

func (v *validate) GetTag() string {
	return v.tag
}

func (v *validate) GetFunc() validator.FuncCtx {
	return v.fnc
}
