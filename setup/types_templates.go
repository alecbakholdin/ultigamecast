package setup

import "text/template"

const fileFunctionFormat = `
func (m *{{.CollectionData.GoName}}) Get{{.GoFieldName}}Path() string {
	if m.{{.GoFieldName}} == "" || m.GetId() == "" {
		return ""
	} else {
		return path.Join("{{.BaseFilepath}}", m.GetId(), m.{{.GoFieldName}})
	}
}
`
var fileFieldFunctionTemplate = template.Must(template.New("file").Parse(fileFunctionFormat))

const dateTimeFunctionsFormat = `
func (m *{{.CollectionData.GoName}}) Get{{.GoName}}Str(format string, locName string) string {
	if dt, err := m.Get{{.GoName}}Dt(); err != nil || dt.IsZero() {
		return ""
	} else if loc, err := time.LoadLocation(locName); err != nil{
		return ""
	} else {
		return dt.Time().In(loc).Format(format)
	}
}

func (m *{{.CollectionData.GoName}}) Get{{.GoName}}Dt() (types.DateTime, error) {
	m.{{.GoName}}Datetime = strings.TrimSpace(m.{{.GoName}}Datetime)
	m.{{.GoName}}Timezone = strings.TrimSpace(m.{{.GoName}}Timezone)
	if m.{{.GoName}}Datetime != "" && m.{{.GoName}}Timezone != "" {
		var datetimeFormat string
		if len(m.{{.GoName}}Datetime) == len("2006-01-02") {
			datetimeFormat = "2006-01-02"
		} else {
			datetimeFormat = "2006-01-02T15:04"
		}
		if loc, err := time.LoadLocation(m.{{.GoName}}Timezone); err != nil {
			return types.DateTime{}, err
		} else if time, err := time.ParseInLocation(datetimeFormat, m.{{.GoName}}Datetime, cmp.Or(loc, time.Local)); err != nil {
			return types.DateTime{}, err
		} else {
			return types.ParseDateTime(time)
		}
	} else if m.{{.GoName}}Datetime != "" {
		return types.ParseDateTime(m.{{.GoName}}Datetime + ":00")
	} else {
		return m.{{.GoName}}, nil
	}
}
`
var fieldDateTimeTemplate = template.Must(template.New("datetime").Parse(dateTimeFunctionsFormat))

const copyFunctionFormat = `
func (d *{{.GoName}}) CopyFrom(s *{{.GoName}}) *{{.GoName}} {
{{range .Fields}}	d.{{.GoName}} = s.{{.GoName}}
{{end}}	return d
}

func (m *{{.GoName}}) Copy() *{{.GoName}} {
	return (&{{.GoName}}{}).CopyFrom(m)
}

`
var collectionCopyFunctionTemplate = template.Must(template.New("copy").Parse(copyFunctionFormat))