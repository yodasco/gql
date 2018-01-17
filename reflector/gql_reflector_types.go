package reflector

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

// GqlOutputAndResolver defines the value side of the TypeMap map.
// It defines the couple of gql Output and gql Resolver
type GqlOutputAndResolver struct {
	Output   graphql.Output
	Resolver graphql.FieldResolveFn
}

// TypeMap defines a mapping b/w a go type and it's graphql Output type
// and resolver function
type TypeMap map[reflect.Type]GqlOutputAndResolver

// ExcludeFieldTag is a type to define struct field annotations that need to be excluded
type ExcludeFieldTag string

// GoName defines a strict field go name
type GoName string

// GqlName defines a struct field graphql name
type GqlName string
