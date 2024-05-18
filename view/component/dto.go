package view_component

import (
	"errors"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate, trans = initValidation()
)

func initValidation() (*validator.Validate, ut.Translator) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	err := en_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		panic(err)
	}
	return validate, trans
}

type DTO struct {
	Errors       validator.ValidationErrors
	CustomErrors []*CustomFieldError
	FormErrors   []string
}
type CustomFieldError struct {
	Field string
	Error string
}

var ErrBinding = errors.New("error validating struct")

func (d *DTO) FieldInvalid(field string) bool {
	return d.FieldError(field) != ""
}

func (d *DTO) FieldInvalidClass(field string) templ.KeyValue[string, bool] {
	return templ.KV("invalid", d.FieldInvalid(field))
}

func (d *DTO) FieldError(field string) string {
	errs := []string{}
	if d.Errors != nil {
		for _, e := range d.Errors {
			if e.Field() == field {
				err := e.Translate(trans)
				errs = append(errs, err)
			}
		}
	}
	if d.CustomErrors != nil {
		for _, e := range d.CustomErrors {
			if e.Field == field {
				errs = append(errs, e.Error)
			}
		}
	}
	return strings.Join(errs, ",")
}

func (d *DTO) Invalid() bool {
	valid := d.Errors == nil && d.CustomErrors == nil && d.FormError() == ""
	return !valid
}

func (d *DTO) Validate(obj interface{}) error {
	fmt.Println(obj)
	err := validate.Struct(obj)
	fmt.Println("err", err)
	if err == nil {
		return nil
	}
	if errs, ok := err.(validator.ValidationErrors); !ok {
		return ErrBinding
	} else {
		d.Errors = errs
	}
	return nil
}

func (d *DTO) AddCustomError(field, err string) {
	if d.CustomErrors == nil {
		d.CustomErrors = []*CustomFieldError{{Field: field, Error: err}}
	} else {
		d.CustomErrors = append(d.CustomErrors, &CustomFieldError{Field: field, Error: err})
	}
}

func (d *DTO) AddFormError(err string) {
	if d.FormErrors == nil {
		d.FormErrors = []string{err}
	} else {
		d.FormErrors = append(d.FormErrors, err)
	}
}

func (d *DTO) FormError() string {
	if d.FormErrors == nil {
		return ""
	}
	return strings.Join(d.FormErrors, ", ")
}
