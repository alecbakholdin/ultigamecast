package mod

import "ultigamecast/web/view/component/dto/field"

type ButtonsModifier struct {
	IncludeButtons bool
}

func (b *ButtonsModifier) Apply(fc *field.FieldConfig) {
	fc.IncludeButtons = true
}

func IncludeButtons() *ButtonsModifier {
	return &ButtonsModifier{
		IncludeButtons: true,
	}
}