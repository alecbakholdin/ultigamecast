package setup

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func CreateTypes(app core.App) {
	dao := app.Dao()
	if baseCollections, err := dao.FindCollectionsByType(models.CollectionTypeBase); err != nil {
		panic(err)
	} else {
		processCollections(baseCollections)
	}
}

func processCollections(collections []*models.Collection) {
	if err := os.MkdirAll("./pbmodels", 0o777); err != nil && !os.IsExist(err) {
		panic(err)
	}

	for _, c := range collections {
		processCollection(c)
	}
}

type CollectionData struct {
	GoName       string
	DbName       string
	Fields       []*CollectionField
	Enums        []*CollectionEnum
	Files        []*CollectionFile
	Dates        []*CollectionField
	IncludesDate bool
}

type CollectionField struct {
	GoName         string
	DbName         string
	GoType         string
	Attributes     map[string]string
	CollectionData *CollectionData
}

type CollectionFile struct {
	BaseFilepath   string
	GoFieldName    string
	CollectionData *CollectionData
}

type CollectionEnum struct {
	Name   string
	Type   string
	Values []string
}

func processCollection(c *models.Collection) {
	data := CollectionData{
		GoName: snakeToPascalCase(c.Name),
		DbName: c.Name,
		Fields: make([]*CollectionField, 0),
		Enums:  make([]*CollectionEnum, 0),
		Files:  make([]*CollectionFile, 0),
		Dates:  make([]*CollectionField, 0),
	}
	for _, fieldSchema := range c.Schema.AsMap() {
		if f := getCollectionField(&data, fieldSchema, c); f != nil {
			data.Fields = append(data.Fields, f...)
		}
	}
	writeCollectionToFile(&data)
}

func getCollectionField(c *CollectionData, f *schema.SchemaField, pbc *models.Collection) []*CollectionField {
	fields := make([]*CollectionField, 1)
	cf := CollectionField{
		DbName: f.Name,
		GoName: snakeToPascalCase(f.Name),
		Attributes: map[string]string{
			"db":   f.Name,
			"json": f.Name,
		},
		CollectionData: c,
	}
	fields[0] = &cf
	switch f.Type {
	case "text":
		cf.GoType = "string"
	case "number":
		if opts, ok := f.Options.(*schema.NumberOptions); !ok {
			fmt.Printf("Missing number options for %s\n", f.Name)
		} else if opts.NoDecimal {
			cf.GoType = "int"
		} else {
			cf.GoType = "float32"
		}
	case "select":
		if opts, ok := f.Options.(*schema.SelectOptions); ok {
			enumName := fmt.Sprintf("%s%s", c.GoName, cf.GoName)
			c.Enums = append(c.Enums, &CollectionEnum{
				Name:   enumName,
				Type:   "string",
				Values: opts.Values,
			})
			cf.GoType = enumName
		} else {
			fmt.Printf("Missing select options for %s: %v\n", f.Name, f.Options)
		}
	case "bool":
		cf.GoType = "bool"
	case "date":
		cf.GoType = "types.DateTime"
		fields = append(fields,
			&CollectionField{
				GoName: fmt.Sprintf("%sTimezone", cf.GoName),
				GoType: "string",
				DbName: fmt.Sprintf("%s_timezone", cf.DbName),
				Attributes: map[string]string{
					"form": fmt.Sprintf("%s_timezone", cf.DbName),
					"json": fmt.Sprintf("%s_timezone", cf.DbName),
					"db":   "-",
				},
				CollectionData: c,
			},
			&CollectionField{
				GoName: fmt.Sprintf("%sDatetime", cf.GoName),
				GoType: "string",
				DbName: fmt.Sprintf("%s_datetime", cf.DbName),
				Attributes: map[string]string{
					"form": fmt.Sprintf("%s_datetime", cf.DbName),
					"json": fmt.Sprintf("%s_datetime", cf.DbName),
					"db":   "-",
				},
				CollectionData: c,
			})
		c.Dates = append(c.Dates, &cf)
		c.IncludesDate = true
	case "relation":
		cf.GoType = "string"
	case "file":
		cf.GoType = "string"
		c.Files = append(c.Files, &CollectionFile{
			BaseFilepath:   pbc.BaseFilesPath(),
			GoFieldName:    cf.GoName,
			CollectionData: c,
		})
	default:
		fmt.Printf("Unsupported type %s for field %s\n", f.Type, f.Name)
		return nil
	}
	if strings.ToLower(cf.DbName) == "slug" {
		cf.Attributes["param"] = fmt.Sprintf("%sSlug", snakeToCamelCase(c.DbName))
	} else if f.Type != "date" {
		cf.Attributes["form"] = cf.DbName
	}
	return fields
}

func snakeToCamelCase(s string) string {
	str := snakeToPascalCase(s)
	if len(str) < 1 {
		return strings.ToLower(str)
	}
	return strings.ToLower(string(str[0])) + str[1:]
}

func snakeToPascalCase(s string) string {
	capNext := true
	str := ""
	for _, c := range s {
		if c == '_' {
			capNext = true
		} else if capNext {
			capNext = false
			str += strings.ToUpper(string(c))
		} else {
			str += string(c)
		}
	}
	return str
}

var defaultImports = []string{"\"github.com/pocketbase/pocketbase/models\""}
var dateImports = []string{"\"time\"", "\"cmp\"", "\"github.com/pocketbase/pocketbase/tools/types\""}

func writeCollectionToFile(c *CollectionData) {
	file := mustGetCollectionFile(c)
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()
	imports := defaultImports
	if c.IncludesDate {
		imports = append(imports, dateImports...)
	}
	if len(c.Files) > 0 {
		imports = append(imports, "\"path\"")
	}
	sort.Strings(imports)

	writer.WriteString("package pbmodels\n\n")
	writer.WriteString(fmt.Sprintf("import (\n\t%s\n)\n\n", strings.Join(imports, "\n\t")))
	writer.WriteString(fmt.Sprintf("type %s struct {\n\tmodels.BaseModel\n\n", c.GoName))

	nameLen := len(slices.MaxFunc(c.Fields, func(a, b *CollectionField) int {
		return len(a.GoName) - len(b.GoName)
	}).GoName)
	typeLen := len(slices.MaxFunc(c.Fields, func(a, b *CollectionField) int {
		return len(a.GoType) - len(b.GoType)
	}).GoType)
	for _, field := range c.Fields {
		attrs := make([]string, 0)
		for key, val := range field.Attributes {
			attrs = append(attrs, fmt.Sprintf(`%s:"%s"`, key, val))
		}
		sort.Strings(attrs)
		writer.WriteString(fmt.Sprintf("\t%s %s `%s`\n", padToLength(field.GoName, nameLen), padToLength(field.GoType, typeLen), strings.Join(attrs, " ")))
	}
	writer.WriteString("}\n\n")
	for _, enum := range c.Enums {
		writer.WriteString(fmt.Sprintf("type %s string\n\n", enum.Name))
		writer.WriteString("const (\n")
		enumValLen := len(slices.MaxFunc(enum.Values, func(a, b string) int {
			return len(a) - len(b)
		}))
		for _, val := range enum.Values {
			writer.WriteString(fmt.Sprintf("\t%s%s %s = \"%s\"\n", enum.Name, padToLength(formatEnumValForType(val), enumValLen), enum.Name, val))
		}
		writer.WriteString(")\n\n")
	}
	writer.WriteString(fmt.Sprintf(`func (m *%s) TableName() string {
    return "%s"
}
`, c.GoName, c.DbName))
	for _, file := range c.Files {
		fileFieldFunctionTemplate.Execute(writer, file)
	}

	for _, dateField := range c.Dates {
		fieldDateTimeTemplate.Execute(writer, dateField)
	}
}

func padToLength(s string, targetLen int) string {
	for len(s) < targetLen {
		s += " "
	}
	return s
}

func formatEnumValForType(val string) string {
	capNext := true
	str := ""
	for _, c := range val {
		if c == '_' || c == ' ' {
			capNext = true
		} else if capNext {
			capNext = false
			str += strings.ToUpper(string(c))
		} else {
			str += string(c)
		}
	}
	return str
}

func mustGetCollectionFile(c *CollectionData) (file *os.File) {
	sanitizedName := fmt.Sprintf("%s.go", strings.ToLower(c.DbName))
	colFilepath := path.Join("./pbmodels", sanitizedName)
	if file, err := os.Create(colFilepath); err != nil {
		panic(err)
	} else {
		return file
	}
}
