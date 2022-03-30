package astnormalization

import "testing"

func TestFoo(t *testing.T) {
	t.Run("Single duplicate scalar is removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataOne, "scalar ScalarOne")
	})
	t.Run("Single duplicate scalar among other data is removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataTwo, "scalar ScalarOne type TypeOne scalar ScalarTwo")
	})
	t.Run("Several duplicate scalars among other data are removed", func(t *testing.T) {
		run(removeScalarDuplicates, "", testDataThree, "type TypeOne scalar ScalarOne type TypeTwo scalar ScalarTwo scalar ScalarThree type TypeThree scalar ScalarFour")
	})
}

const testDataOne = `
scalar ScalarOne
scalar ScalarOne
`
const testDataTwo = `
scalar ScalarOne
type TypeOne
scalar ScalarOne
scalar ScalarOne
scalar ScalarTwo
`
const testDataThree = `
type TypeOne
scalar ScalarOne
type TypeTwo
scalar ScalarOne
scalar ScalarTwo
scalar ScalarOne
scalar ScalarTwo
scalar ScalarThree
type TypeThree
scalar ScalarFour
scalar ScalarThree
scalar ScalarOne
`
