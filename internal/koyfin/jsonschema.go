package koyfin

import (
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

type SchemaModifier func(sc *jsonschema.Schema)

func ModifiedSchema(base *jsonschema.Schema, modifiers ...SchemaModifier) *jsonschema.Schema {
	targetSchema := base.CloneSchemas()
	// finally apply modifier
	for _, mod := range modifiers {
		mod(targetSchema)
	}
	return targetSchema
}

// create new schema from a base schema with a
func ModifyTargetPath(path string, modifiers ...SchemaModifier) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		targetSchema := sc
		var indx int

		// empty string to match base schema
		if path == "" {
			indx = 1
		}

		paths := strings.Split(path, ".")
		// traverse into inner schema properties
		for indx = 0; indx < len(paths); indx++ {
			if targetSchema.Properties != nil {
				targetSchema = targetSchema.Properties[paths[indx]]
			}
			break
		}
		// finally apply modifier
		for _, mod := range modifiers {
			mod(targetSchema)
		}
	}
}

// ChangeDescription modifies a schema's description field
func ChangeDescription(description string) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		sc.Description = description
	}
}

// SetType modifies a schema's Type field
func SetType(newType string) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		sc.Type = newType
		// TODO: empty all field specific value
	}
}

// SetRequiredFields modifies required fields schema
func SetRequiredFields(fields ...string) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		required := []string{}
		// ensure matching existing before settings required
		for _, f := range fields {
			_, ok := sc.Properties[f]
			if ok {
				required = append(required, f)
			}
		}
		sc.Required = required
	}
}

// SetEnumItems modifies a schema's type and sets enum values
func SetEnumItems[T any](values ...T) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		if len(values) == 0 {
			return
		}

		// check input type
		var newType string
		switch reflect.TypeOf(values[0]).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			newType = "integer"
		case reflect.Float32, reflect.Float64:
			newType = "number"
		case reflect.String:
			newType = "string"
		case reflect.Bool:
			newType = "boolean"
		}

		// unrecognized types - ignore
		if newType == "" {
			return
		}
		sc.Type = newType
		// values -> []any
		for _, val := range values {
			sc.Enum = append(sc.Enum, val)
		}
	}
}

// SetEnum modifies a schema's type and sets
func SetEnum[T any](values ...T) SchemaModifier {
	return SetEnumItems(values...)
}

// SetEnum modifies a schema's type and sets
func SetConstant(value any, showType bool) SchemaModifier {
	return func(sc *jsonschema.Schema) {
		// showing type is optional, if shown then make sure correct
		var newType string
		if showType {
			switch reflect.TypeOf(value).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				newType = "integer"
			case reflect.Float32, reflect.Float64:
				newType = "number"
			case reflect.String:
				newType = "string"
			case reflect.Bool:
				newType = "boolean"
			// cannot infer type - omit
			default:
				newType = ""
			}
		}
		sc.Type = newType
		sc.Const = jsonschema.Ptr(value)
	}
}
