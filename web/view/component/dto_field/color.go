package dto_field

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
