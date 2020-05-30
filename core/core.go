package core

type Feed interface {
	Fetch() error
	Stream() error
}

type GraphQL struct {
	ID    string   `json:"uid" dgraph:"uid,omitempty"`
	Types []string `json:"dgraph.type" dgraph:"dgraph.type,omitempty"`
}
