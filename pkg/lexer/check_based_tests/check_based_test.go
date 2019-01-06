package tabledriven_selfmade_green

import (
	"fmt"
	"github.com/jensneuse/graphql-go-tools/pkg/lexer"
	. "github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	"testing"
)

func TestLexer_Peek_Read(t *testing.T) {

	type checkFunc func(lex *lexer.Lexer, i int)

	run := func(input string, checks ...checkFunc) {
		lex := lexer.NewLexer()
		lex.SetInput(input)
		for i := range checks {
			checks[i](lex, i+1)
		}
	}

	mustPeek := func(k Keyword) checkFunc {
		return func(lex *lexer.Lexer, i int) {
			peeked, err := lex.Peek(true)
			if err != nil {
				panic(err)
			}
			if peeked != k {
				panic(fmt.Errorf("mustPeek: want: %s, got: %s [check: %d]", k.String(), peeked.String(), i))
			}
		}
	}

	mustRead := func(k Keyword, literal string) checkFunc {
		return func(lex *lexer.Lexer, i int) {
			tok, err := lex.Read()
			if err != nil {
				panic(err)
			}
			if k != tok.Keyword {
				panic(fmt.Errorf("mustRead: want(keyword): %s, got: %s [check: %d]", k.String(), tok.String(), i))
			}
			if literal != tok.Literal {
				panic(fmt.Errorf("mustRead: want(literal): %s, got: %s [check: %d]", literal, tok.Literal, i))
			}
		}
	}

	mustPeekAndRead := func(k Keyword, literal string) checkFunc {
		return func(lex *lexer.Lexer, i int) {
			mustPeek(k)(lex, i)
			mustRead(k, literal)(lex, i)
		}
	}

	mustErrRead := func() checkFunc {
		return func(lex *lexer.Lexer, i int) {
			_, err := lex.Read()
			if err == nil {
				panic(fmt.Errorf("mustErrRead: want error, got nil [check: %d]", i))
			}
		}
	}

	t.Run("integer", func(t *testing.T) {
		run("1337", mustPeekAndRead(INTEGER, "1337"))
	})
	t.Run("integer whitespace", func(t *testing.T) {
		run(" 1337 ", mustPeekAndRead(INTEGER, "1337"))
	})
	t.Run("integer comma", func(t *testing.T) {
		run("1337,", mustPeekAndRead(INTEGER, "1337"))
	})
	t.Run("float", func(t *testing.T) {
		run("13.37", mustPeekAndRead(FLOAT, "13.37"))
	})
	t.Run("float whitespace", func(t *testing.T) {
		run(" 13.37 ", mustPeekAndRead(FLOAT, "13.37"))
	})
	t.Run("float comma", func(t *testing.T) {
		run("13.37,", mustPeekAndRead(FLOAT, "13.37"))
	})
	t.Run("float invalid", func(t *testing.T) {
		run("1.3.3.7,", mustErrRead())
	})
	t.Run("multiple integers and floats", func(t *testing.T) {
		run("1337 13.37,13337	113377 \n 1.337",
			mustPeekAndRead(INTEGER, "1337"),
			mustPeekAndRead(FLOAT, "13.37"),
			mustPeekAndRead(INTEGER, "13337"),
			mustPeekAndRead(INTEGER, "113377"),
			mustPeekAndRead(FLOAT, "1.337"))
	})
}
