package field

type Modifier interface {
	Apply(*FieldConfig)
}

type TextColor string

const (
	TextColorDefault TextColor = ""
	TextColorLink    TextColor = "has-text-link"
	TextColorPrimary TextColor = "has-text-primary"
	TextColorInfo    TextColor = "has-text-info"
	TextColorSuccess TextColor = "has-text-success"
	TextColorWarning TextColor = "has-text-warning"
	TextColorDanger  TextColor = "has-text-danger"
)

type FieldColor string

const (
	FieldColorDefault FieldColor = ""
	FieldColorLink    FieldColor = "is-link"
	FieldColorPrimary FieldColor = "is-primary"
	FieldColorInfo    FieldColor = "is-info"
	FieldColorSuccess FieldColor = "is-success"
	FieldColorWarning FieldColor = "is-warning"
	FieldColorDanger  FieldColor = "is-danger"
)

type InputType string

const (
	InputTypeText     InputType = "text"
	InputTypeEmail    InputType = "email"
	InputTypePassword InputType = "password"
	InputTypeHidden   InputType = "hidden"
)

type Size string

const (
	SizeSmall  Size = "is-small"
	SizeNormal Size = ""
	SizeMedium Size = "is-medium"
	SizeLarge  Size = "is-large"
)

type FaSize string

const (
	FaSizeSmall  FaSize = "fa-sm"
	FaSizeNormal FaSize = ""
	FaSizeLarge  FaSize = "fa-lg"
)
