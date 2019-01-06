package plain_green

import (
	"github.com/jensneuse/graphql-go-tools/pkg/lexer"
	"github.com/jensneuse/graphql-go-tools/pkg/lexing/keyword"
	"testing"
)

func TestLexer_Peek_Integer(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1337")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), key.String())
	}
}

func TestLexer_Peek_IntegerWithSpace(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput(" 1337 ")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), key.String())
	}
}

func TestLexer_Peek_Integer_WithComma(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1337,")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), key.String())
	}
}

func TestLexer_Read_Integer(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1337")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), token.Keyword.String())
	}
	if token.Literal != "1337" {
		t.Fatalf("want literal: %s, got: %s", "1337", token.Literal)
	}
}

func TestLexer_Read_Integer_WithSpace(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput(" 1337 ")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), token.Keyword.String())
	}
	if token.Literal != "1337" {
		t.Fatalf("want literal: %s, got: %s", "1337", token.Literal)
	}
}

func TestLexer_Read_Integer_WithComma(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1337,")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.INTEGER {
		t.Fatalf("want keyword %s, got: %s", keyword.INTEGER.String(), token.Keyword.String())
	}
	if token.Literal != "1337" {
		t.Fatalf("want literal: %s, got: %s", "1337", token.Literal)
	}
}

func TestLexer_Peek_Float(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("13.37")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), key.String())
	}
}

func TestLexer_Peek_Float_WithWhitespace(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput(" 13.37 ")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), key.String())
	}
}

func TestLexer_Peek_Float_WithComma(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("13.37,")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), key.String())
	}
}

func TestLexer_Read_Float(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("13.37")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), token.Keyword.String())
	}
	if token.Literal != "13.37" {
		t.Fatalf("want literal: %s, got: %s", "13.37", token.Literal)
	}
}

func TestLexer_Read_Float_WithWhitespace(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput(" 13.37 ")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), token.Keyword.String())
	}
	if token.Literal != "13.37" {
		t.Fatalf("want literal: %s, got: %s", "13.37", token.Literal)
	}
}

func TestLexer_Read_Float_WithComma(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("13.37,")
	token, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if token.Keyword != keyword.FLOAT {
		t.Fatalf("want keyword %s, got: %s", keyword.FLOAT.String(), token.Keyword.String())
	}
	if token.Literal != "13.37" {
		t.Fatalf("want literal: %s, got: %s", "13.37", token.Literal)
	}
}

func TestLexer_Read_Float_Invalid(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1.3.3.7,")
	key, err := lex.Peek(true)
	if err != nil {
		t.Fatal(err)
	}
	if key != keyword.INTEGER {
		t.Fatalf("want integer, got: %s", key.String())
	}
	_, err = lex.Read()
	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestLexer_Read_Multiple(t *testing.T) {
	lex := lexer.NewLexer()
	lex.SetInput("1337 13.37,13337	113377 \n 1.337")
	first, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if first.Keyword != keyword.INTEGER {
		t.Fatalf("want integer, got: %s", first.Keyword.String())
	}
	if first.Literal != "1337" {
		t.Fatalf("want 1337, got: %s", first.Literal)
	}
	second, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if second.Keyword != keyword.FLOAT {
		t.Fatalf("want float, got: %s", second.Keyword.String())
	}
	if second.Literal != "13.37" {
		t.Fatalf("want 13.37, got: %s", second.Literal)
	}
	third, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if third.Keyword != keyword.INTEGER {
		t.Fatalf("want integer, got: %s", third.Keyword.String())
	}
	if third.Literal != "13337" {
		t.Fatalf("want 13337, got: %s", third.Literal)
	}
	fourth, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if fourth.Keyword != keyword.INTEGER {
		t.Fatalf("want integer, got: %s", fourth.Keyword.String())
	}
	if fourth.Literal != "113377" {
		t.Fatalf("want 113377, got: %s", fourth.Literal)
	}
	fifth, err := lex.Read()
	if err != nil {
		t.Fatal(err)
	}
	if fifth.Keyword != keyword.FLOAT {
		t.Fatalf("want integer, got: %s", fifth.Keyword.String())
	}
	if fifth.Literal != "1.337" {
		t.Fatalf("want 1.337, got: %s", fifth.Literal)
	}
}
