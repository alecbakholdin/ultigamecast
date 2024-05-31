package tournament_data_types

import "fmt"

type Option int64

const (
	Text     Option = 0
	Location Option = 1
	Link     Option = 2
)
var Options = []Option{Text, Location, Link}

func Title(o Option) string {
	switch Option(o) {
	case Text:
		return "Text"
	case Location:
		return "Location"
	case Link:
		return "Link"
	default:
		panic(fmt.Sprintf("Unexpected title option %d", o))
	}
}


func Icon(o Option) string {
	switch o {
	case Text:
		return "description"
	case Location:
		return "location_on"
	case Link:
		return "link"
	default:
		panic(fmt.Sprintf("Unexpected icon option %d", o))
	}
}

func OptionIconTitle(o Option) (Option, string, string) {
	return o, Icon(o), Title(o)
}