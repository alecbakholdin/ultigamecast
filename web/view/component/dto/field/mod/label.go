package mod

import "ultigamecast/web/view/component/dto/field"

type LabelModifier struct {
	Label string
}

func (l *LabelModifier) Apply(fc *field.FieldConfig) {
	fc.Label = l.Label
}

func Label(l string) *LabelModifier {
	return &LabelModifier{
		Label: l,
	}
}
