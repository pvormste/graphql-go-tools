package astvalidation

import (
	"github.com/pvormste/graphql-go-tools/pkg/astvisitor"
)

var reservedFieldPrefix = []byte("__")

// Rule is hook to register callback functions on the Walker
type Rule func(walker *astvisitor.Walker)
