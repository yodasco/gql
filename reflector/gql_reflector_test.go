package reflector

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicType(t *testing.T) {
	val := 6
	gqlt := ReflectGqlType("a", reflect.TypeOf(val), GetDefaultTypeMap(), ExcludeFieldTag(""))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return val, nil
		},
	}
	assertQuery(t, f, "s", "", `{"data":{"s": 6}}`, "")
}

func TestStringKind(t *testing.T) {
	type stringKind string
	s := stringKind("sss")
	gqlt := ReflectGqlType("a", reflect.TypeOf(s), GetDefaultTypeMap(), ExcludeFieldTag(""))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return s, nil
		},
	}
	assertQuery(t, f, "s", "", `{"data":{"s": "sss"}}`, "")
}
func TestSimpleStruct(t *testing.T) {
	type S struct {
		A string `json:"a"`
	}

	gqlt := ReflectGqlType("s", reflect.TypeOf(S{}), GetDefaultTypeMap(), ExcludeFieldTag(""))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return S{
				A: "hello world",
			}, nil
		},
	}
	assertQuery(t, f, "s", "{a}", `{"data":{"s":{"a":"hello world"}}}`, "")
}

func TestExclude(t *testing.T) {
	type S struct {
		A string `json:"a" gqlexclude:"ignore"`
		B string `json:"b" gqlexclude:"f,ignore ,h"`
		C string `json:"c" gqlexclude:"c"`
	}

	gqlt := ReflectGqlType("s", reflect.TypeOf(S{}), GetDefaultTypeMap(), ExcludeFieldTag("ignore"))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return S{
				A: "hello world",
				B: "hello world",
				C: "hello world",
			}, nil
		},
	}
	assertQuery(t, f, "s", "{a}", "", `Cannot query field "a" on type`)
	assertQuery(t, f, "s", "{b}", "", `Cannot query field "b" on type`)
	assertQuery(t, f, "s", "{c}", `{"data":{"s":{"c":"hello world"}}}`, "")
}

func TestFieldsWithoutJSON(t *testing.T) {
	type S struct {
		A string `json:"a"`
		B string
	}

	gqlt := ReflectGqlType("s", reflect.TypeOf(S{}), GetDefaultTypeMap(), ExcludeFieldTag("ignore"))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return S{
				A: "hello world",
				B: "hello world",
			}, nil
		},
	}
	assertQuery(t, f, "s", "{a}", `{"data":{"s":{"a":"hello world"}}}`, "")
	assertQuery(t, f, "s", "{b}", "", `Cannot query field "b" on type `)
}

func TestDatatypes(t *testing.T) {
	type DataTypes struct {
		Bool       bool       `json:"bool"`
		String     string     `json:"string"`
		Int        int        `json:"int"`
		Int8       int8       `json:"int_8"`
		Int16      int16      `json:"int_16"`
		Int32      int32      `json:"int_32"`
		Int64      int64      `json:"int_64"`
		Uint       uint       `json:"uint"`
		Uint8      uint8      `json:"uint_8"`
		Uint16     uint16     `json:"uint_16"`
		Uint32     uint32     `json:"uint_32"`
		Uint64     uint64     `json:"uint_64"`
		Uintptr    uintptr    // Doesn't seem useful
		Byte       byte       `json:"byte"`
		Rune       rune       `json:"rune"`
		Float32    float32    `json:"float_32"`
		Float64    float64    `json:"float_64"`
		Complex64  complex64  `json:"complex_64"`
		Complex128 complex128 `json:"complex_128"`
		Time       time.Time  `json:"time"`
	}

	gqlt := ReflectGqlType("data_types", reflect.TypeOf(DataTypes{}), GetDefaultTypeMap(), ExcludeFieldTag(""))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return DataTypes{
				Bool:       true,
				String:     "string",
				Int:        -5,
				Int8:       -5,
				Int16:      -5,
				Int32:      -5,
				Int64:      -5,
				Uint:       5,
				Uint8:      5,
				Uint16:     5,
				Uint32:     5,
				Uint64:     5,
				Byte:       5,
				Rune:       5,
				Float32:    5.5,
				Float64:    5.5,
				Complex64:  complex(float32(-1), 1),
				Complex128: complex(-1, 1),
				Time:       time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
			}, nil
		},
	}
	assertQuery(t, f, "dt", `{
		bool
		string
		int
		int_8
		int_16
		int_32
		int_64
		uint
		uint_8
		uint_16
		uint_32
		uint_64
		byte
		rune
		float_32
		float_64
		complex_64
		complex_128
		time
	}`, `{
		"data":
		{
			"dt":
			{
				"bool": true,
				"string": "string",
				"int": -5,
				"int_8": -5,
				"int_16": -5,
				"int_32": -5,
				"int_64": -5,
				"uint": 5,
				"uint_8": 5,
				"uint_16": 5,
				"uint_32": 5,
				"uint_64": 5,
				"byte": 5,
				"rune": 5,
				"float_32": 5.5,
				"float_64": 5.5,
				"complex_64": "(-1+1i)",
				"complex_128": "(-1+1i)",
				"time": "2009-11-10T23:00:00Z"
			}
		}
	}`, "")
}

func TestComplexStruct(t *testing.T) {
	type T1 struct {
		G  string        `json:"g"`
		I  interface{}   `json:"i"`
		Is []interface{} `json:"is"`
	}

	type T struct {
		A        string   `json:"a"`
		B        int      `json:"b"`
		S        []string `json:"s"`
		SingleT1 T1       `json:"single_t_1"`
		ManyT1s  []T1     `json:"many_t_1_s"`
		C        int      `json:"c" gqlexclude:"ignore_me"`
		C1       int      `json:"c1" gqlexclude:"ignore_me2,ignore_me,ignore_me3"`
	}

	gqlt := ReflectGqlType("t", reflect.TypeOf(T{}), GetDefaultTypeMap(), ExcludeFieldTag("ignore_me"))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return T{
				A: "hello world",
				B: 6,
				SingleT1: T1{
					G: "GGG",
					I: "a string",
					Is: []interface{}{
						"another string",
						5,
						"and yes",
					},
				},
				S: []string{"A", "B", "C"},
				ManyT1s: []T1{
					{G: "G1"},
					{G: "G2"},
				},
			}, nil
		},
	}
	assertQuery(t, f, "t", `{
		a
		b
		single_t_1 {
			g
			i
			is
		}
		s
		many_t_1_s {
			g
			i
			is
		}
	}`, `
	{
		"data":
		{
			"t":
			{
				"b":6,
				"many_t_1_s":
				[
					{
						"is":[],
						"g":"G1",
						"i": null
					},
					{
						"g": "G2",
						"i": null,
						"is":[]
					}
				],
				"s":["A", "B", "C"],
				"single_t_1":
				{
					"i": "a string",
					"is":["another string", "5", "and yes"],
					"g":"GGG"
				},
				"a": "hello world"
			}
		}
	}`, "")
}

func TestArray(t *testing.T) {
	gqlt := ReflectGqlType("a", reflect.TypeOf([]string{}), GetDefaultTypeMap(), ExcludeFieldTag(""))
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			return []string{"hello", "world"}, nil
		},
	}
	assertQuery(t, f, "s", "", `{"data":{"s":["hello", "world"]}}`, "")
}

// runs a graphql query and asserts the result
// in case there query should result in an error then set the expectedError
// argument to non-empty string. This string should be a substript  of the
// expected returned error (for example `Cannot query field "b"`)
func assertQuery(
	t *testing.T,
	f graphql.Field,
	rootQuery,
	query,
	expectedResult,
	expectedError string,
) {
	req := require.New(t)
	fields := graphql.Fields{
		rootQuery: &f,
	}
	tp := graphql.NewObject(graphql.ObjectConfig{
		Name:   fmt.Sprintf("root_%s", rootQuery),
		Fields: fields,
	})
	schemaConfig := graphql.SchemaConfig{Query: tp}
	schema, err := graphql.NewSchema(schemaConfig)
	req.Nil(err)

	query = fmt.Sprintf("query{%s%s}", rootQuery, query)
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		if expectedError == "" {
			req.Fail("failed to execute graphql operation", "%+v", r.Errors)
		} else {
			errs := fmt.Sprintf("%+v", r.Errors)
			if strings.Contains(errs, expectedError) {
				// OK, pass
				return
			}
			req.Fail("Expected error, but got a different error",
				"Expected: %s. Actual: %s",
				expectedError, errs)
		}
	} else {
		if expectedError != "" {
			req.Fail("Expected an error, but there was no error",
				"Expected error: %s", expectedError)
		}
	}
	result, err := json.Marshal(r)
	req.Nil(err)
	assert.JSONEq(t, expectedResult, string(result))
}
