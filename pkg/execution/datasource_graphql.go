package execution

import (
	"bytes"
	"encoding/json"
	"github.com/buger/jsonparser"
	"github.com/davecgh/go-spew/spew"
	"github.com/jensneuse/graphql-go-tools/pkg/ast"
	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/jensneuse/graphql-go-tools/pkg/astvisitor"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer/literal"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type GraphQLDataSourcePlanner struct {
	walker                *astvisitor.Walker
	operation, definition *ast.Document
	args                  []Argument
	nodes                 []ast.Node
	resolveDocument       *ast.Document

	rootFieldRef          int
	rootFieldArgumentRefs []int
	variableDefinitions   []int
}

func (g *GraphQLDataSourcePlanner) DirectiveName() []byte {
	return []byte("GraphQLDataSource")
}

func (g *GraphQLDataSourcePlanner) Initialize(walker *astvisitor.Walker, operation, definition *ast.Document, args []Argument, resolverParameters []ResolverParameter) {
	g.walker, g.operation, g.definition, g.args = walker, operation, definition, args

	g.resolveDocument = &ast.Document{}
	g.rootFieldArgumentRefs = make([]int, len(resolverParameters))
	g.variableDefinitions = make([]int, len(resolverParameters))
	g.rootFieldRef = -1
	for i := 0; i < len(resolverParameters); i++ {
		g.resolveDocument.VariableValues = append(g.resolveDocument.VariableValues, ast.VariableValue{
			Name: g.resolveDocument.Input.AppendInputBytes(resolverParameters[i].name),
		})
		variableRef := len(g.resolveDocument.VariableValues) - 1
		variableValue := ast.Value{
			Kind: ast.ValueKindVariable,
			Ref:  variableRef,
		}
		g.resolveDocument.Arguments = append(g.resolveDocument.Arguments, ast.Argument{
			Name:  g.resolveDocument.Input.AppendInputBytes(resolverParameters[i].name),
			Value: variableValue,
		})
		g.rootFieldArgumentRefs[i] = len(g.resolveDocument.Arguments) - 1

		g.resolveDocument.Types = append(g.resolveDocument.Types, ast.Type{
			TypeKind: ast.TypeKindNamed,
			Name:     g.resolveDocument.Input.AppendInputBytes([]byte("String")),
			OfType:   -1,
		})

		stringTypeRef := len(g.resolveDocument.Types) - 1
		g.resolveDocument.Types = append(g.resolveDocument.Types, ast.Type{
			TypeKind: ast.TypeKindNonNull,
			OfType:   stringTypeRef,
		})

		nonNullTypeRef := len(g.resolveDocument.Types) - 1

		g.resolveDocument.VariableDefinitions = append(g.resolveDocument.VariableDefinitions, ast.VariableDefinition{
			VariableValue: variableValue,
			Type:          nonNullTypeRef,
		})
		g.variableDefinitions[i] = len(g.resolveDocument.VariableDefinitions) - 1
	}
}

func (g *GraphQLDataSourcePlanner) EnterInlineFragment(ref int) {
	if len(g.nodes) == 0 {
		return
	}
	current := g.nodes[len(g.nodes)-1]
	if current.Kind != ast.NodeKindSelectionSet {
		return
	}
	inlineFragmentType := g.resolveDocument.ImportType(g.operation.InlineFragments[ref].TypeCondition.Type, g.operation)
	g.resolveDocument.InlineFragments = append(g.resolveDocument.InlineFragments, ast.InlineFragment{
		TypeCondition: ast.TypeCondition{
			Type: inlineFragmentType,
		},
		SelectionSet: -1,
	})
	inlineFragmentRef := len(g.resolveDocument.InlineFragments) - 1
	g.resolveDocument.Selections = append(g.resolveDocument.Selections, ast.Selection{
		Kind: ast.SelectionKindInlineFragment,
		Ref:  inlineFragmentRef,
	})
	selectionRef := len(g.resolveDocument.Selections) - 1
	g.resolveDocument.SelectionSets[current.Ref].SelectionRefs = append(g.resolveDocument.SelectionSets[current.Ref].SelectionRefs, selectionRef)
	g.nodes = append(g.nodes, ast.Node{
		Kind: ast.NodeKindInlineFragment,
		Ref:  inlineFragmentRef,
	})
}

func (g *GraphQLDataSourcePlanner) LeaveInlineFragment(ref int) {
	g.nodes = g.nodes[:len(g.nodes)-1]
}

func (g *GraphQLDataSourcePlanner) EnterSelectionSet(ref int) {

	fieldOrInlineFragment := g.nodes[len(g.nodes)-1]

	set := ast.SelectionSet{}
	g.resolveDocument.SelectionSets = append(g.resolveDocument.SelectionSets, set)
	setRef := len(g.resolveDocument.SelectionSets) - 1

	switch fieldOrInlineFragment.Kind {
	case ast.NodeKindField:
		g.resolveDocument.Fields[fieldOrInlineFragment.Ref].HasSelections = true
		g.resolveDocument.Fields[fieldOrInlineFragment.Ref].SelectionSet = setRef
	case ast.NodeKindInlineFragment:
		g.resolveDocument.InlineFragments[fieldOrInlineFragment.Ref].HasSelections = true
		g.resolveDocument.InlineFragments[fieldOrInlineFragment.Ref].SelectionSet = setRef
	}

	g.nodes = append(g.nodes, ast.Node{
		Kind: ast.NodeKindSelectionSet,
		Ref:  setRef,
	})
}

func (g *GraphQLDataSourcePlanner) LeaveSelectionSet(ref int) {
	g.nodes = g.nodes[:len(g.nodes)-1]
}

func (g *GraphQLDataSourcePlanner) EnterField(ref int) {
	if g.rootFieldRef == -1 {
		g.rootFieldRef = ref

		fieldNameValue, ok := g.walker.FieldDefinitionDirectiveArgumentValueByName(ref, g.DirectiveName(), []byte("field"))
		if !ok {
			return
		}

		if fieldNameValue.Kind != ast.ValueKindString {
			return
		}

		field := ast.Field{
			Name: g.resolveDocument.Input.AppendInputBytes(g.definition.StringValueContentBytes(fieldNameValue.Ref)),
			Arguments: ast.ArgumentList{
				Refs: g.rootFieldArgumentRefs,
			},
			HasArguments: len(g.rootFieldArgumentRefs) != 0,
		}
		g.resolveDocument.Fields = append(g.resolveDocument.Fields, field)
		fieldRef := len(g.resolveDocument.Fields) - 1
		selection := ast.Selection{
			Kind: ast.SelectionKindField,
			Ref:  fieldRef,
		}
		g.resolveDocument.Selections = append(g.resolveDocument.Selections, selection)
		selectionRef := len(g.resolveDocument.Selections) - 1
		set := ast.SelectionSet{
			SelectionRefs: []int{selectionRef},
		}
		g.resolveDocument.SelectionSets = append(g.resolveDocument.SelectionSets, set)
		setRef := len(g.resolveDocument.SelectionSets) - 1
		operationDefinition := ast.OperationDefinition{
			Name:          g.resolveDocument.Input.AppendInputBytes([]byte("o")),
			OperationType: g.operation.OperationDefinitions[g.walker.Ancestors[0].Ref].OperationType,
			SelectionSet:  setRef,
			HasSelections: true,
			VariableDefinitions: ast.VariableDefinitionList{
				Refs: g.variableDefinitions,
			},
			HasVariableDefinitions: len(g.variableDefinitions) != 0,
		}
		g.resolveDocument.OperationDefinitions = append(g.resolveDocument.OperationDefinitions, operationDefinition)
		operationDefinitionRef := len(g.resolveDocument.OperationDefinitions) - 1
		g.resolveDocument.RootNodes = append(g.resolveDocument.RootNodes, ast.Node{
			Kind: ast.NodeKindOperationDefinition,
			Ref:  operationDefinitionRef,
		})

		g.nodes = append(g.nodes, ast.Node{
			Kind: ast.NodeKindOperationDefinition,
			Ref:  operationDefinitionRef,
		})
		g.nodes = append(g.nodes, ast.Node{
			Kind: ast.NodeKindSelectionSet,
			Ref:  setRef,
		})
		g.nodes = append(g.nodes, ast.Node{
			Kind: ast.NodeKindField,
			Ref:  fieldRef,
		})
	} else {
		field := ast.Field{
			Name: g.resolveDocument.Input.AppendInputBytes(g.operation.FieldNameBytes(ref)),
		}
		g.resolveDocument.Fields = append(g.resolveDocument.Fields, field)
		fieldRef := len(g.resolveDocument.Fields) - 1
		set := g.nodes[len(g.nodes)-1]
		selection := ast.Selection{
			Kind: ast.SelectionKindField,
			Ref:  fieldRef,
		}
		g.resolveDocument.Selections = append(g.resolveDocument.Selections, selection)
		selectionRef := len(g.resolveDocument.Selections) - 1
		g.resolveDocument.SelectionSets[set.Ref].SelectionRefs = append(g.resolveDocument.SelectionSets[set.Ref].SelectionRefs, selectionRef)
		g.nodes = append(g.nodes, ast.Node{
			Kind: ast.NodeKindField,
			Ref:  fieldRef,
		})
	}
}

func (g *GraphQLDataSourcePlanner) LeaveField(ref int) {

	if g.rootFieldRef == ref {

		buff := bytes.Buffer{}
		err := astprinter.Print(g.resolveDocument, nil, &buff)
		if err != nil {
			g.walker.StopWithInternalErr(err)
			return
		}
		arg := &StaticVariableArgument{
			Name:  literal.QUERY,
			Value: buff.Bytes(),
		}
		g.args = append([]Argument{arg}, g.args...)

		definition, exists := g.walker.FieldDefinition(ref)
		if !exists {
			return
		}
		directive, exists := g.definition.FieldDefinitionDirectiveByName(definition, []byte("GraphQLDataSource"))
		if !exists {
			return
		}
		value, exists := g.definition.DirectiveArgumentValueByName(directive, literal.URL)
		if !exists {
			return
		}
		variableValue := g.definition.StringValueContentBytes(value.Ref)
		arg = &StaticVariableArgument{
			Name:  literal.URL,
			Value: variableValue,
		}
		g.args = append([]Argument{arg}, g.args...)
		value, exists = g.definition.DirectiveArgumentValueByName(directive, literal.HOST)
		if !exists {
			return
		}
		variableValue = g.definition.StringValueContentBytes(value.Ref)
		arg = &StaticVariableArgument{
			Name:  literal.HOST,
			Value: variableValue,
		}
		g.args = append([]Argument{arg}, g.args...)
	}

	g.nodes = g.nodes[:len(g.nodes)-1]
}

func (g *GraphQLDataSourcePlanner) Plan() (DataSource, []Argument) {
	return &GraphQLDataSource{}, g.args
}

type GraphQLDataSource struct{}

func (g *GraphQLDataSource) Resolve(ctx Context, args ResolvedArgs, out io.Writer) {

	hostArg := args.ByKey(literal.HOST)
	urlArg := args.ByKey(literal.URL)
	queryArg := args.ByKey(literal.QUERY)

	if hostArg == nil || urlArg == nil || queryArg == nil {
		spew.Dump(args)
		log.Fatal("one of host,url,query arg nil")
		return
	}

	url := string(hostArg) + string(urlArg)
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}

	variables := map[string]json.RawMessage{}
	for i := 0; i < len(args); i++ {
		key := args[i].Key
		switch {
		case bytes.Equal(key, literal.HOST):
		case bytes.Equal(key, literal.URL):
		case bytes.Equal(key, literal.QUERY):
		default:
			variables[string(key)] = args[i].Value
		}
	}

	gqlRequest := GraphqlRequest{
		OperationName: "o",
		Variables:     variables,
		Query:         string(queryArg),
	}

	gqlRequestData, err := json.MarshalIndent(gqlRequest, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			TLSHandshakeTimeout: 0 * time.Second,
		},
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(gqlRequestData))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data = bytes.ReplaceAll(data, literal.BACKSLASH, nil)
	data, _, _, err = jsonparser.Get(data, "data")
	if err != nil {
		log.Fatal(err)
	}
	out.Write(data)
}
