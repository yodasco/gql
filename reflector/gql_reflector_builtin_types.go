package reflector

import (
	"reflect"
	"time"

	"github.com/graphql-go/graphql"
)

var defaultTypeMap TypeMap

func init() {
	defaultTypeMap = buildDefaultTypeMap()
}

func trivialResolver(p graphql.ResolveParams) (interface{}, error) {
	reflected := reflect.ValueOf(p.Source)
	fieldName := p.Info.FieldName
	value := findFieldByTag(reflect.Indirect(reflected), "json", GqlName(fieldName))
	return value.Interface(), nil
}

func findFieldByTag(v reflect.Value, tagName string, fieldName GqlName) reflect.Value {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := GetFieldFirstTag(f, tagName)
		if tag == string(fieldName) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

// Format time as string in RFC3339
func timeResolver(p graphql.ResolveParams) (interface{}, error) {
	reflected := reflect.ValueOf(p.Source)
	fieldName := p.Info.FieldName
	value := findFieldByTag(reflect.Indirect(reflected), "json", GqlName(fieldName))
	t := value.Interface().(time.Time)
	return t.UTC().Format(time.RFC3339), nil
}

func buildDefaultTypeMap() TypeMap {
	return TypeMap{
		reflect.TypeOf(""): {
			Output:   graphql.String,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(int(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(int8(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(int16(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(int32(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(int64(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(uint(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(uint8(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(uint16(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(uint32(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(uint64(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(byte(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(rune(5)): {
			Output:   graphql.Int,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(float32(5)): {
			Output:   graphql.Float,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(float64(5)): {
			Output:   graphql.Float,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(true): {
			Output:   graphql.Boolean,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(complex(float32(5), 5)): { // complex64
			Output:   graphql.String,  // String repr is good enough for now
			Resolver: trivialResolver, // example: (1+1i) (1real and 1img)
		},
		reflect.TypeOf(complex(float64(5), 5)): { // complex128
			Output:   graphql.String,
			Resolver: trivialResolver,
		},
		reflect.TypeOf(time.Now()): {
			Output:   graphql.String,
			Resolver: timeResolver,
		},
	}
}

// GetDefaultTypeMap returns a default type map, including all the native types
func GetDefaultTypeMap() TypeMap {
	return defaultTypeMap
}
