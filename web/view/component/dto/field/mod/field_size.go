package mod

import "ultigamecast/web/view/component/dto/field"

type FieldSizeModifier struct {
	Size     field.Size
	IconSize field.FaSize
}

func (f *FieldSizeModifier) Apply(fc *field.FieldConfig) {
	fc.Size = f.Size
	fc.FaIconSize = f.IconSize
}

func FieldSizeSmall() *FieldSizeModifier {
	return &FieldSizeModifier{
		Size:     field.SizeSmall,
		IconSize: field.FaSizeNormal,
	}
}

func FieldSizeMedium() *FieldSizeModifier {
	return &FieldSizeModifier{
		Size:     field.SizeMedium,
		IconSize: field.FaSizeLarge,
	}
}

func FieldSizeLarge() *FieldSizeModifier {
	return &FieldSizeModifier{
		Size:     field.SizeLarge,
		IconSize: field.FaSizeLarge,
	}
}
