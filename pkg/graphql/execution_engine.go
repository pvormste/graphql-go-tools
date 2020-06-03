package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/cespare/xxhash"
	"github.com/jensneuse/abstractlogger"

	"github.com/jensneuse/graphql-go-tools/pkg/astprinter"
	"github.com/jensneuse/graphql-go-tools/pkg/execution"
	"github.com/jensneuse/graphql-go-tools/pkg/execution/datasource"
	"github.com/jensneuse/graphql-go-tools/pkg/operationreport"
)

type DataSourceHttpJsonOptions struct {
	HttpClient         *http.Client
	WhitelistedSchemes []string
}

type DataSourceGraphqlOptions struct {
	HttpClient         *http.Client
	WhitelistedSchemes []string
}

type ExecutionOptions struct {
	ExtraArguments json.RawMessage
}

type ExecutionEngine struct {
	logger           abstractlogger.Logger
	basePlanner      *datasource.BasePlanner
	executorPool     *sync.Pool
	plannerPool      *sync.Pool
	printerPool      *sync.Pool
	hash64Pool       *sync.Pool
	schema           *Schema
	queryPlanCache   map[uint64]execution.RootNode
	queryPlanCacheMu sync.RWMutex
}

func NewExecutionEngine(logger abstractlogger.Logger, schema *Schema, plannerConfig datasource.PlannerConfiguration) (*ExecutionEngine, error) {

	basePlanner, err := datasource.NewBaseDataSourcePlanner(schema.rawInput, plannerConfig, logger)
	if err != nil {
		return nil, err
	}

	return &ExecutionEngine{
		logger:      logger,
		basePlanner: basePlanner,
		executorPool: &sync.Pool{
			New: func() interface{} {
				return execution.NewExecutor(nil)
			},
		},
		plannerPool: &sync.Pool{
			New: func() interface{} {
				return execution.NewPlanner(basePlanner)
			},
		},
		hash64Pool: &sync.Pool{
			New: func() interface{} {
				return xxhash.New()
			},
		},
		printerPool: &sync.Pool{
			New: func() interface{} {
				return &astprinter.Printer{}
			},
		},
		schema:         schema,
		queryPlanCache: make(map[uint64]execution.RootNode, 1024),
	}, nil
}

func (e *ExecutionEngine) AddHttpJsonDataSource(name string) error {
	return e.AddHttpJsonDataSourceWithOptions(name, DataSourceHttpJsonOptions{})
}

func (e *ExecutionEngine) AddHttpJsonDataSourceWithOptions(name string, options DataSourceHttpJsonOptions) error {
	httpJsonFactoryFactory := &datasource.HttpJsonDataSourcePlannerFactoryFactory{}

	if options.HttpClient != nil {
		httpJsonFactoryFactory.Client = options.HttpClient
	}

	if len(options.WhitelistedSchemes) > 0 {
		httpJsonFactoryFactory.WhitelistedSchemes = options.WhitelistedSchemes
	}

	return e.AddDataSource(name, httpJsonFactoryFactory)
}

func (e *ExecutionEngine) AddGraphqlDataSource(name string) error {
	return e.AddGraphqlDataSourceWithOptions(name, DataSourceGraphqlOptions{})
}

func (e *ExecutionEngine) AddGraphqlDataSourceWithOptions(name string, options DataSourceGraphqlOptions) error {
	graphqlFactoryFactory := &datasource.GraphQLDataSourcePlannerFactoryFactory{}

	if options.HttpClient != nil {
		graphqlFactoryFactory.Client = options.HttpClient
	}

	if len(options.WhitelistedSchemes) > 0 {
		graphqlFactoryFactory.WhitelistedSchemes = options.WhitelistedSchemes
	}

	return e.AddDataSource(name, graphqlFactoryFactory)
}

func (e *ExecutionEngine) AddDataSource(name string, plannerFactoryFactory datasource.PlannerFactoryFactory) error {
	return e.basePlanner.RegisterDataSourcePlannerFactory(name, plannerFactoryFactory)
}

func (e *ExecutionEngine) ExecuteWithWriter(ctx context.Context, operation *Request, writer io.Writer, options ExecutionOptions) error {
	var report operationreport.Report

	if !operation.IsNormalized() {
		normalizationResult, err := operation.Normalize(e.schema)
		if err != nil {
			return err
		}

		if !normalizationResult.Successful {
			return normalizationResult.Errors
		}
	}

	operationID, err := e.operationID(operation)
	if err != nil {
		return err
	}

	e.queryPlanCacheMu.RLock()
	plan, exists := e.queryPlanCache[operationID]
	e.queryPlanCacheMu.RUnlock()
	if !exists {
		planner := e.plannerPool.Get().(*execution.Planner)
		plan = planner.Plan(&operation.document, e.basePlanner.Definition, &report)
		e.plannerPool.Put(planner)
		if report.HasErrors() {
			return report
		}
		e.queryPlanCacheMu.Lock()
		e.queryPlanCache[operationID] = plan
		e.queryPlanCacheMu.Unlock()
	}

	variables, extraArguments := execution.VariablesFromJson(operation.Variables, options.ExtraArguments)
	executionContext := execution.Context{
		Context:        ctx,
		Variables:      variables,
		ExtraArguments: extraArguments,
	}

	poolExecutor := e.executorPool.Get().(*execution.Executor)
	defer e.executorPool.Put(poolExecutor)
	return poolExecutor.Execute(executionContext, plan, writer)
}

func (e *ExecutionEngine) operationID(operation *Request) (uint64, error) {
	hash64 := e.hash64Pool.Get().(hash.Hash64)
	printer := e.printerPool.Get().(*astprinter.Printer)
	err := printer.Print(&operation.document, &e.schema.document, hash64)
	result := hash64.Sum64()
	hash64.Reset()
	e.hash64Pool.Put(hash64)
	e.printerPool.Put(printer)
	return result, err
}

func (e *ExecutionEngine) Execute(ctx context.Context, operation *Request, options ExecutionOptions) (*ExecutionResult, error) {
	var buf bytes.Buffer
	err := e.ExecuteWithWriter(ctx, operation, &buf, options)
	return &ExecutionResult{&buf}, err
}

type ExecutionResult struct {
	buf *bytes.Buffer
}

func (r *ExecutionResult) Buffer() *bytes.Buffer {
	return r.buf
}

func (r *ExecutionResult) GetAsHTTPResponse() (res *http.Response) {
	if r.buf == nil {
		return
	}

	res = &http.Response{}
	res.Body = ioutil.NopCloser(r.buf)
	res.Header = make(http.Header)
	res.StatusCode = 200

	res.Header.Set("Content-Type", "application/json")

	return
}
