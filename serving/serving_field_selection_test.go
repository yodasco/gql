package serving

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yodasco/gql/reflector"
)

func TestGetSelectedFields(t *testing.T) {
	type S struct {
		A string `json:"a"`
		B string `json:"b"`
		C string `json:"c"`
		D struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"d"`
	}

	as := assert.New(t)
	gqlt := reflector.ReflectType(S{})
	f := graphql.Field{
		Type: gqlt,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {

			// When selectionPath is nil, the top-most field should be selected
			// in this case it is "s"
			selected := GetSelectedFields(nil, p)
			as.Equal([]string{"s"}, selected)

			// Sub-selectin of s should be "a"
			selected = GetSelectedFields([]string{"s"}, p)
			as.Equal([]string{"a", "b", "d"}, selected)

			// Sub--subselectin of s should be "a"
			selected = GetSelectedFields([]string{"s", "d"}, p)
			as.Equal([]string{"x"}, selected)

			// When path is not found, response should be empty
			selected = GetSelectedFields([]string{"xxx"}, p)
			as.Equal([]string{}, selected)

			return S{
				A: "hello world",
			}, nil
		},
	}
	runQuery(t, f, "s", "{a b d{x}}")
}

func runQuery(
	t *testing.T,
	f graphql.Field,
	rootQuery,
	query string,
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
	req.NotNil(r)
	req.Empty(r.Errors)
}
