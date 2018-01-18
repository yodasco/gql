package serving

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// GetSelectedFields returns a string slice, a list of graphql selected field
// names in the given graphql `selectionPath`
// This function may be used in runtime to determine the list of fields selected
// by a user when running a specific graphql query
func GetSelectedFields(
	selectionPath []string,
	resolveParams graphql.ResolveParams,
) []string {
	fields := resolveParams.Info.FieldASTs
	for _, propName := range selectionPath {
		found := false
		for _, field := range fields {
			if field.Name.Value == propName && field.SelectionSet != nil {
				selections := field.SelectionSet.Selections
				fields = make([]*ast.Field, 0)
				for _, selection := range selections {
					fields = append(fields, selection.(*ast.Field))
				}
				found = true
				break
			}
		}
		if !found {
			return []string{}
		}
	}
	var collect []string
	for _, field := range fields {
		collect = append(collect, field.Name.Value)
	}
	return collect
}
