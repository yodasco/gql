# gql
GraphQL utilities in Go

## ReflectType Example

Create `graphql.Field` using `reflector.ReflectType`

```go
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

```go

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


# Common pitfalls

## Adding your own data types to the type map.
This library supports all go built-in data types, so for example it understands that a go type of `int8` should be defined as a `graphql.Integer` etc.
It also supports simple derived types, for example `type Email string` is defined as a `graphql.String`.

If you get the error `failed to create new schema, error: price_usd_5 fields must be an object with field names as keys or a function which return such an object.` (where `price_usd_5` is just an example), that means that you have a field named `price_usd` with a data type that's not supported.
Here's an example how to fix this:

```go
    gqlt := reflector.ReflectTypeWithTypeMap(
        A{},
        getMyTypeMap())

func getMyTypeMap() reflect.TypeMap {
	// clone a local version of the default type map and add to it
	typeMap = make(reflector.TypeMap)
	for t, outputAndResolver := range reflector.GetDefaultTypeMap() {
		typeMap[t] = outputAndResolver
	}

	// suport for sql.NullString
	typeMap[reflect.TypeOf(sql.NullString{})] = reflector.OutputAndResolver{
		Output: graphql.String,
		Resolver: func(p graphql.ResolveParams) (interface{}, error) {
			value := reflector.GetValueFromResolveParams(p)
			v := value.Interface().(sql.NullString)
			if v.Valid {
				return v.String, nil
			}
			return nil, nil
		},
	}
}
```

