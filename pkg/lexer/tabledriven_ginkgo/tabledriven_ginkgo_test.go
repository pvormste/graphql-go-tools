package tabledriven_ginkgo

import (
	"github.com/jensneuse/graphql-go-tools/pkg/lexer"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"testing"
)

func TestLexer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test lexer")
}

var _ = Describe("test Peek & Read", func() {
	type Case struct {
		input             string
		expectPeekKeyword types.GomegaMatcher
		expectReadKeyword types.GomegaMatcher
		expectLiteral     types.GomegaMatcher
		expectPeekErr     types.GomegaMatcher
		expectReadErr     types.GomegaMatcher
	}

	var lex *lexer.Lexer

	BeforeEach(func() {
		lex = lexer.NewLexer()
	})

	DescribeTable("lexer test cases", func(c Case) {

		lex.SetInput(c.input)
		out, err := lex.Peek(true)
		if c.expectPeekErr != nil {
			Expect(err).To(c.expectPeekErr)
		}
		if c.expectPeekKeyword != nil {
			Expect(out).To(c.expectPeekKeyword)
		}
		actualToken, err := lex.Read()
		if c.expectReadErr != nil {
			Expect(err).To(c.expectReadErr)
		}
		if c.expectReadKeyword != nil {
			Expect(actualToken.Keyword).To(c.expectReadKeyword)
		}

	},
		Entry("integer", Case{
			input:             "1337",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("integer with whitespace", Case{
			input:             " 1337 ",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("integer with comma", Case{
			input:             "foo,",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("float", Case{
			input:             "13.37",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("float with whitespace", Case{
			input:             " 13.37 ",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("float with comma", Case{
			input:             "13.37,",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		}),
		Entry("float invalid", Case{
			input:             "1.3.3.7,",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.UNDEFINED),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     HaveOccurred(),
		}),
	)
})
