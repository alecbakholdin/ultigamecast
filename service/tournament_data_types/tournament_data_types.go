package tournament_data_types

type Option int64

const (
	Text     Option = 0
	Location Option = 1
	Link     Option = 2
)

var Options = []Option{Text, Location, Link}
var OptionViews = []struct {
	Value Option
	Icon  string
	Title string
}{
	{Text, "description", "Text"},
	{Location, "location_on", "Location"},
	{Link, "link", "Link"},
}

var OptionView = map[Option]struct {
	Icon  string
	Title string
}{
	Text:     {"description", "Text"},
	Location: {"location_on", "Location"},
	Link:     {"link", "Link"},
}
