package reflector

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/graphql-go/graphql"
)

const (
	// GqlExcludeTagName is the name of the struct field tag to use for exclusions.
	GqlExcludeTagName = "gqlexclude"
)

// ReflectType is a shorthand method for invoking ReflectTypeFq.
// It derives the most basic fields from the given instance (typically a struct)
// such as the name of the struct, and uses the default type mapping and no
// exclude tags at all.
func ReflectType(instance interface{}) graphql.Type {
	if instance == nil {
		panic("Cannot infer type of nil instance")
	}
	t := reflect.TypeOf(instance)
	return ReflectTypeFq(
		GqlName(t.Name()),
		t,
		GetDefaultTypeMap(),
		ExcludeFieldTag(""),
	)
}

// ReflectTypeWithTypeMap is a shorthand to invoking ReflectTypeEq with
// reasonable default values and your own provided type map.
func ReflectTypeWithTypeMap(instance interface{}, typeMap TypeMap) graphql.Type {
	if instance == nil {
		panic("Cannot infer type of nil instance")
	}
	t := reflect.TypeOf(instance)
	return ReflectTypeFq(
		GqlName(t.Name()),
		t,
		typeMap,
		ExcludeFieldTag(""),
	)
}

// ReflectTypeFq (Reflect Type Fully Qualified) returns a Graphql type that
// represents the go reflect.Type structure (recorsively)
func ReflectTypeFq(
	name GqlName,
	t reflect.Type,
	typeMap TypeMap,
	exclude ExcludeFieldTag,
) graphql.Type {
	gqlType := getGqlType(t, typeMap)
	if gqlType != nil {
		return gqlType
	}
	switch t.Kind() {
	case reflect.String:
		return graphql.String
	case reflect.Interface:
		// for interfaces assume type String. Correct assumption?
		return graphql.String
	case reflect.Bool:
		return graphql.Boolean
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return graphql.Int
	case reflect.Float32, reflect.Float64:
		return graphql.Float
	case reflect.Struct:
		fields := ReflectFieldsFq(t, typeMap, exclude)
		return graphql.NewObject(graphql.ObjectConfig{
			Name:   generateGqlOTypeName(name),
			Fields: fields,
		})
	case reflect.Slice, reflect.Array:
		return graphql.NewList(ReflectTypeFq(name, t.Elem(), typeMap, exclude))
	case reflect.Invalid:
		panic(fmt.Sprintf("Invalid GQL kind %s. Field: %s", t.Kind(), t.Name()))
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer:
		panic(fmt.Sprintf("Unsupported GQL kind %s. Field: %s", t.Kind(), t.Name()))
	default:
		panic(fmt.Sprintf("Unknown GO kind %s. Field: %s", t.Kind(), t.Name()))
	}
}

// ReflectFieldsFq returns a `graphql.Fields` map of the t struct.
// t must of a struct
func ReflectFieldsFq(
	t reflect.Type,
	typeMap TypeMap,
	exclude ExcludeFieldTag,
) graphql.Fields {
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf(`ReflectFieldsFq can only work on struct types.
			Received instead %s`, t.Kind()))
	}
	fields := make(graphql.Fields)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if includeField(f, exclude) {
			name := GqlName(GetFieldFirstTag(f, "json"))
			fields[string(name)] = ReflectFieldFq(name, f.Type, typeMap, exclude)
		}
	}
	return fields
}

// ReflectFieldFq returns a Graphql field that represents
// the go reflect.Field (recorsively)
func ReflectFieldFq(
	name GqlName,
	t reflect.Type,
	typeMap TypeMap,
	exclude ExcludeFieldTag,
) *graphql.Field {
	gqlType := ReflectTypeFq(name, t, typeMap, exclude)
	resolver := getResolver(t, typeMap)
	return &graphql.Field{
		Name:    string(name),
		Type:    gqlType,
		Resolve: resolver,
	}
}

var gqlTypeNameOrder = 0

func generateGqlOTypeName(name GqlName) string {
	gqlTypeNameOrder++
	return fmt.Sprintf("%s_%d", name, gqlTypeNameOrder)
}

// Whether to include this StructField in the gql schema
// excludeTagName is the value of gqlexclude to search for exclusion list
func includeField(f reflect.StructField, exclude ExcludeFieldTag) bool {
	fieldName := GetFieldFirstTag(f, "json")
	if fieldName == "" {
		return false
	}
	gqlexclude := f.Tag.Get(GqlExcludeTagName)
	if gqlexclude == "" {
		// No exclusions
		return true
	}
	for _, s := range strings.Split(gqlexclude, ",") {
		if strings.Trim(s, " ") == string(exclude) {
			// excluded
			return false
		}
	}
	return true
}

// GetFieldFirstTag gets the StructField first tag value. Empty string if the tag
// does not exist.
// First by means of coma separated
func GetFieldFirstTag(field reflect.StructField, tag string) string {
	return strings.Trim(strings.Split(field.Tag.Get(tag), ",")[0], " ")
}

// Get te gql type of the go type t.
// If doesn't exist - return String as default
func getGqlType(t reflect.Type, typeMap TypeMap) graphql.Output {
	m, exists := typeMap[t]
	if exists {
		return m.Output
	}
	// no predefined output in the map
	return nil
}

func getResolver(t reflect.Type, typeMap TypeMap) graphql.FieldResolveFn {
	m, exists := typeMap[t]
	if !exists {
		// By default use the trivial resolver
		return trivialResolver
	}
	return m.Resolver
}
