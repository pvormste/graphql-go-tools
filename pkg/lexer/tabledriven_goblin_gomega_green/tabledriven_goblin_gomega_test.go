package tabledriven_goblin_gomega_green

import (
	. "github.com/franela/goblin"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"testing"
)

func TestLexer_Peek_Read(t *testing.T) {

	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	for _, tt := range [...]struct {
		name              string
		input             string
		expectPeekKeyword types.GomegaMatcher
		expectReadKeyword types.GomegaMatcher
		expectLiteral     types.GomegaMatcher
		expectPeekErr     types.GomegaMatcher
		expectReadErr     types.GomegaMatcher
	}{
		{
			name:              "integer",
			input:             "1337",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "integer with whitespace",
			input:             " 1337 ",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "integer with comma",
			input:             "1337,",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.INTEGER),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "float",
			input:             "13.37",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "float with whitespace",
			input:             " 13.37 ",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "float with comma",
			input:             "13.37,",
			expectPeekKeyword: Equal(keyword.FLOAT),
			expectReadKeyword: Equal(keyword.FLOAT),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     Not(HaveOccurred()),
		},
		{
			name:              "float invalid",
			input:             "1.3.3.7,",
			expectPeekKeyword: Equal(keyword.INTEGER),
			expectReadKeyword: Equal(keyword.UNDEFINED),
			expectPeekErr:     Not(HaveOccurred()),
			expectReadErr:     HaveOccurred(),
		},
	} {
		g.Describe("lexer tests", func() {
			g.It(tt.name, func() {
				lex := lexer.NewLexer()
				lex.SetInput(tt.input)
				out, err := lex.Peek(true)
				if tt.expectPeekErr != nil {
					Expect(err).To(tt.expectPeekErr)
				}
				if tt.expectPeekKeyword != nil {
					Expect(out).To(tt.expectPeekKeyword)
				}
				actualToken, err := lex.Read()
				if tt.expectReadErr != nil {
					Expect(err).To(tt.expectReadErr)
				}
				if tt.expectReadKeyword != nil {
					Expect(actualToken.Keyword).To(tt.expectReadKeyword)
				}
			})
		})
	}
}
