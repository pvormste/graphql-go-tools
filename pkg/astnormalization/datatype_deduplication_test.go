package astnormalization

import "testing"

const foo = `
scalar FooBar
type Bananas
scalar FooBar
scalar FooBar
scalar BarFoo
`

const bar = `
scalar Foo
scalar BarFoo
type Bananas
scalar FooBar
scalar FooBar
scalar BarFoo
scalar Foo
scalar Bar
scalar Bar
`

func TestFoo(t *testing.T) {
	t.Run("", func(t *testing.T) {
		run(dataTypeDeduplication, "", foo, "scalar FooBar type Bananas scalar BarFoo")
	})
	t.Run("", func(t *testing.T) {
		run(dataTypeDeduplication, "", bar, "scalar Foo scalar BarFoo type Bananas scalar FooBar scalar Bar")
	})
}
