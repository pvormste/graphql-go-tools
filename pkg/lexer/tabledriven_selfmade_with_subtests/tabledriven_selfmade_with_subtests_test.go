package tabledriven_selfmade_green

import (
	"github.com/jensneuse/graphql-go-tools/pkg/lexer"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	"testing"
)

func TestLexer_Peek_Read(t *testing.T) {
	for _, tt := range [...]struct {
		name              string
		input             string
		expectPeekKeyword keyword.Keyword
		expectReadKeyword keyword.Keyword
		expectLiteral     string
		expectPeekErr     bool
		expectReadErr     bool
	}{
		{
			name:              "integer",
			input:             "1337",
			expectPeekKeyword: keyword.INTEGER,
			expectReadKeyword: keyword.INTEGER,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "integer with whitespace",
			input:             " 1337 ",
			expectPeekKeyword: keyword.INTEGER,
			expectReadKeyword: keyword.INTEGER,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "integer with comma",
			input:             "1337,",
			expectPeekKeyword: keyword.INTEGER,
			expectReadKeyword: keyword.INTEGER,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "float",
			input:             "13.37",
			expectPeekKeyword: keyword.FLOAT,
			expectReadKeyword: keyword.FLOAT,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "float with whitespace",
			input:             " 13.37 ",
			expectPeekKeyword: keyword.FLOAT,
			expectReadKeyword: keyword.FLOAT,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "float with comma",
			input:             "13.37,",
			expectPeekKeyword: keyword.FLOAT,
			expectReadKeyword: keyword.FLOAT,
			expectPeekErr:     false,
			expectReadErr:     false,
		},
		{
			name:              "float invalid",
			input:             "1.3.3.7,",
			expectPeekKeyword: keyword.INTEGER,
			expectReadKeyword: keyword.UNDEFINED,
			expectPeekErr:     false,
			expectReadErr:     true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			lex := lexer.NewLexer()
			lex.SetInput(tt.input)
			out, err := lex.Peek(true)
			if tt.expectPeekErr && err == nil {
				t.Errorf("want error, got nil")
			}
			if !tt.expectPeekErr && err != nil {
				t.Error(err)
			}
			if !tt.expectPeekErr && out != tt.expectPeekKeyword {
				t.Errorf("want: %s, got: %s", tt.expectPeekKeyword.String(), out.String())
			}
			actualToken, err := lex.Read()
			if tt.expectReadErr && err == nil {
				t.Errorf("want error, got nil")
			}
			if !tt.expectReadErr && err != nil {
				t.Error(err)
			}
			if actualToken.Keyword != tt.expectReadKeyword {
				t.Errorf("want: %s, got: %s", tt.expectReadKeyword.String(), actualToken.Keyword)
			}
		})
	}
}
