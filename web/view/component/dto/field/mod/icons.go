package mod

import "ultigamecast/web/view/component/dto/field"

type IconsModifier struct {
	LeftIcon string
	RightIcon string
}

func (i *IconsModifier) Apply(fc *field.FieldConfig) {
	if i.LeftIcon != "" {
		fc.FaLeftIcon = i.LeftIcon
	}
	if i.RightIcon != "" {
		fc.FaRightIcon = i.RightIcon
	}
}

func Icons(l, r string) *IconsModifier {
	return &IconsModifier{
		LeftIcon: l,
		RightIcon: r,
	}
}

func IconLeft(l string) *IconsModifier {
	return Icons(l, "")
}

func IconRight(r string) *IconsModifier {
	return Icons("", r)
}