package reflector

import (
	"reflect"

	"github.com/graphql-go/graphql"
)

// TypeMap defines a mapping b/w a go type and it's graphql Output type
// and resolver function
type TypeMap map[reflect.Type]struct {
	Output   graphql.Output
	Resolver graphql.FieldResolveFn
}

// ExcludeFieldTag is a type to define struct field annotations that need to be excluded
type ExcludeFieldTag string

// GoName defines a strict field go name
type GoName string

// GqlName defines a struct field graphql name
type GqlName string
