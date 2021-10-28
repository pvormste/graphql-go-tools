# graphql-go-tools test tasks

### Example graphql schema for task 1 & task 2

```graphql
union SearchResult = Human | Droid | Starship

schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

type Query {
    hero: Character
    droid(id: ID!): Droid
    search(name: String!): SearchResult
}

type Mutation {
    createReview(episode: Episode!, review: ReviewInput!): Review
}

type Subscription {
    remainingJedis: Int!
}

input ReviewInput {
    stars: Int!
    commentary: String
}

type Review {
    id: ID!
    stars: Int!
    commentary: String
}

enum Episode {
    NEWHOPE
    EMPIRE
    JEDI
}

interface Character {
    name: String!
    friends: [Character]
}

type Human implements Character {
    name: String!
    height: String!
    friends: [Character]
}

type Droid implements Character {
    name: String!
    primaryFunction: String!
    friends: [Character]
}

type Starship {
    name: String!
    length: Float!
}
```


## Tasks

### [1. Task 1 - Fix failing tests (Mandatory) ](https://github.com/TykTechnologies/graphql-go-tools/blob/test-task/pkg/testtask/task_1.go)

### [2. Task 2 - Collect document stats (Try to do)](https://github.com/TykTechnologies/graphql-go-tools/blob/test-task/pkg/testtask/task_2.go)

#### Task Description

Goal of this task is to collect some stats about graphl document:

- Collect names of ObjectTypeDefinitions
- Collect uniq field names of ObjectTypeDefinitions

#### Task Tips:


### [3. Task 3 - Create AST Programmatically (Optional)](https://github.com/TykTechnologies/graphql-go-tools/blob/test-task/pkg/testtask/task_3.go)

In this task we have following example graphql schema

```graphql
schema {
    query: Query
}

type Query {
    droid: Droid!
    hero(id: ID!): Character
}

interface Character {
    name: String!
}

type Droid implements Character {
    name: String!
}
```

We need to recreate ast representing this schema programatically using import helpers from [ast package](https://github.com/TykTechnologies/graphql-go-tools/tree/master/pkg/ast)

Task is already bootstrapped and has a corresponding test

[Introspection Converter](https://github.com/TykTechnologies/graphql-go-tools/blob/master/pkg/introspection/converter.go) contains an examples of how to use import helpers


