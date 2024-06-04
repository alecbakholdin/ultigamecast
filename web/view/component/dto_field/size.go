package dto_field

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

func MediumField() Modifier {
	return func(c *FieldConfig) {
		c.Size = SizeMedium
		c.FaIconSize = FaSizeLarge
	}
}

func SmallField() Modifier {
	return func(c *FieldConfig) {
		c.Size = SizeSmall
		c.FaIconSize = FaSizeNormal
	}
}

func LargeField() Modifier {
	return func(c *FieldConfig) {
		c.Size = SizeLarge
		c.FaIconSize = FaSizeLarge
	}
}
