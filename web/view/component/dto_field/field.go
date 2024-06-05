package dto_field

import (
	"github.com/a-h/templ"
)

type FieldErrorAccessor interface {
	FieldError(string) string
}

type FieldConfig struct {
	Dto      FieldErrorAccessor
	DtoField string

	Type       InputType
	Label      string
	Name       string
	FieldColor FieldColor

	HelpText      string
	HelpTextColor TextColor

	TooltipHelpText string

	Autocomplete string
	Placeholder  string
	Size         Size

	IconSize    Size
	FaLeftIcon  string
	FaRightIcon string
	FaIconSize  FaSize
}

func NewFieldConfig(dto FieldErrorAccessor, dtoField string, modifiers ...Modifier) *FieldConfig {
	config := &FieldConfig{
		Dto:        dto,
		DtoField:   dtoField,
		Size:       SizeNormal,
		IconSize:   SizeNormal,
		FaIconSize: FaSizeLarge,
	}

	if e := dto.FieldError(dtoField); e != "" {
		modifiers = append(modifiers, Error(e))
	}

	for _, m := range modifiers {
		m(config)
	}
	return config
}

func (fc *FieldConfig) LeftIconClass() templ.KeyValue[string, bool] {
	return templ.KV("has-icons-left", fc.FaLeftIcon != "")
}

func (fc *FieldConfig) RightIconClass() templ.KeyValue[string, bool] {
	return templ.KV("has-icons-right", fc.FaRightIcon != "")
}
