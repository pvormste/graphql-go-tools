package graph

import (
	"github.com/jensneuse/graphql-go-tools/examples/federation/accounts/graph/model"
)

var users = map[string]model.User{
	"1": {
		ID: "1",
		Name: &model.Name{
			First: stringPtr("Ada"),
			Last:  stringPtr("Lovelace"),
		},
		BirthDate: stringPtr("1815-12-10"),
		Username:  stringPtr("@ada"),
		Ssn:       stringPtr("123-45-6789"),
		// TODO: account: { __typename: 'LibraryAccount', id: '1' },
	},
	"2": {
		ID: "2",
		Name: &model.Name{
			First: stringPtr("Alan"),
			Last:  stringPtr("Turing"),
		},
		BirthDate: stringPtr("1912-06-23"),
		Username:  stringPtr("@complete"),
		Account: model.SMSAccount{
			Number: stringPtr("8675309"),
		},
		Ssn: stringPtr("987-65-4321"),
	},
}

func stringPtr(s string) *string {
	return &s
}
