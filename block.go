package dgman

import (
	"context"
	"strings"

	"github.com/dgraph-io/dgo/v200"
)

type Block struct {
	ctx context.Context
	tx  *dgo.Txn
	q   *Query
}

func (b *Block) Query(q *Query) *Block {
	b.q = q
	return b
}

func (b *Block) String() string {
	var queryBuf strings.Builder

	b.q.generateQuery(&queryBuf)

	return queryBuf.String()
}
