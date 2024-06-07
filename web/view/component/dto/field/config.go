package field

import (
	"github.com/a-h/templ"
)

type FieldErrorAccessor interface {
	FieldError(string) string
}

type FieldConfig struct {
	Dto      FieldErrorAccessor
	DtoField string
	
	Label      string
	Size         Size
	FieldColor FieldColor
	
	HelpText      string
	HelpTextColor TextColor
	TooltipHelpText string
	
	InputAttributes templ.Attributes

	IconSize    Size
	FaLeftIcon  string
	LeftIconColor TextColor
	FaRightIcon string
	RightIconColor TextColor
	FaIconSize  FaSize
}

func NewFieldConfig(dto FieldErrorAccessor, dtoField string, modifiers ...Modifier) *FieldConfig {
	config := &FieldConfig{
		Dto:        dto,
		DtoField:   dtoField,
		Size:       SizeNormal,
		IconSize:   SizeNormal,
		FaIconSize: FaSizeLarge,
		InputAttributes: templ.Attributes{},
	}


	for _, m := range modifiers {
		m.Apply(config)
	}

	if e := dto.FieldError(dtoField); e != "" {
		config.HelpText = e
		config.HelpTextColor = TextColorDanger
		config.FieldColor = FieldColorDanger
		if config.FaRightIcon == "" {
			config.FaRightIcon = "fa-exclamation-triangle"
			config.RightIconColor = TextColorDanger
		}
	}

	return config
}

func (fc *FieldConfig) LeftIconClass() templ.KeyValue[string, bool] {
	return templ.KV("has-icons-left", fc.FaLeftIcon != "")
}

func (fc *FieldConfig) RightIconClass() templ.KeyValue[string, bool] {
	return templ.KV("has-icons-right", fc.FaRightIcon != "")
}
