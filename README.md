# gql
GraphQL utilities in Go

## ReflectType Example

Create `graphql.Field` using `reflector.ReflectType`

```
import (
	"github.com/graphql-go/graphql"

	"github.com/yodasco/gql/reflector"
)

// A is an example struct type
type A struct {
    X string `json:"x"`
}

func GetAField() *graphql.Field {
	args := graphql.FieldConfigArgument{
		"url": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "URL input (just an example)",
		},
	}

	gqlt := reflector.ReflectType(A{})
	field := graphql.Field{
		Type:        gqlt,
		Description: "Get an A",
		Args:        args,
		Resolve:     resolveA,
	}
	return &field
}

func resolveA(p graphql.ResolveParams) (interface{}, error) {
    return A{
        x: "hello world",
    }, nil
}
```

## ReflectTypeFq Example

Create `graphql.Field` using `reflector.ReflectGqlType`

```

import (
	"github.com/graphql-go/graphql"

	"github.com/yodasco/gql/reflector"
)

// A is an example struct type
type A struct {
    X string `json:"x"`
    Ignored string `json:"ignored" gqlexclude:"ignore"`
}

func GetAField() *graphql.Field {
	args := graphql.FieldConfigArgument{
		"url": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "URL input (just an example)",
		},
	}

	gqlt := reflector.ReflectTypeFq(
		"a",
		reflect.TypeOf(A{}),
		reflector.GetDefaultTypeMap(),
		reflector.ExcludeFieldTag("ignore"))
	field := graphql.Field{
		Type:        gqlt,
		Description: "Get an A",
		Args:        args,
		Resolve:     resolveA,
	}
	return &field
}
```