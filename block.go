package dgman

import (
	"bytes"
	"context"
	"github.com/dgraph-io/dgo/v200"
)

type Block struct {
	ctx   context.Context
	tx    *dgo.Txn
	model interface{}
	q     *Query
	as    string
}

func (b *Block) As(as string) *Block {
	b.as = as
	return b
}

// Query defines the query portion other than the root function for deletion
func (b *Block) Query(query string, params ...interface{}) *Block {
	b.q.query = parseQueryWithParams(query, params)
	return b
}

// Filter defines a query filter, return predicates at the first depth
func (b *Block) Filter(filter string, params ...interface{}) *Block {
	b.q.filter = parseQueryWithParams(filter, params)
	return b
}

// UID returns the node with the specified uid
func (b *Block) UID(uid string) *Block {
	b.q.uid = uid
	return b
}

// RootFunc modifies the dgraph query root function, if not set,
// the default is "type(NodeType)"
func (b *Block) RootFunc(rootFunc string) *Block {
	b.q.rootFunc = rootFunc
	return b
}

func (b *Block) First(n int) *Block {
	b.q.first = n
	return b
}

func (b *Block) Offset(n int) *Block {
	b.q.offset = n
	return b
}

func (b *Block) After(uid string) *Block {
	b.q.after = uid
	return b
}

func (b *Block) OrderAsc(clause string) *Block {
	b.q.order = append(b.q.order, order{clause: clause})
	return b
}

func (b *Block) OrderDesc(clause string) *Block {
	b.q.order = append(b.q.order, order{descending: true, clause: clause})
	return b
}

func (b *Block) String() string {
	var queryBuf bytes.Buffer

	// START ROOT FUNCTION
	if b.as == "" {
		queryBuf.WriteString("\n\tvar(func: ")
	} else {
		queryBuf.WriteString("\n\t")
		queryBuf.WriteString(b.as)
		queryBuf.WriteString(" AS var(func: ")
	}

	if b.q.uid != "" {
		queryBuf.WriteString("uid(")
		queryBuf.WriteString(b.q.uid)
		queryBuf.WriteString(")")
	} else {
		if b.q.rootFunc == "" {
			// if root function is not defined, query from node type
			nodeType := GetNodeType(b.q.model)
			queryBuf.WriteString("type(")
			queryBuf.WriteString(nodeType)
			queryBuf.WriteByte(')')
		} else {
			queryBuf.WriteString(b.q.rootFunc)
		}

		if b.q.first != 0 {
			queryBuf.WriteString(", first: ")
			queryBuf.Write(intToBytes(b.q.first))
		}

		if b.q.offset != 0 {
			queryBuf.WriteString(", offset: ")
			queryBuf.Write(intToBytes(b.q.offset))
		}

		if b.q.after != "" {
			queryBuf.WriteString(", after: ")
			queryBuf.WriteString(b.q.after)
		}

		if len(b.q.order) > 0 {
			for _, order := range b.q.order {
				orderStr := ", orderasc: "
				if order.descending {
					orderStr = ", orderdesc: "
				}
				queryBuf.WriteString(orderStr)
				queryBuf.WriteString(order.clause)
			}
		}
	}
	queryBuf.WriteString(") ")
	// END ROOT FUNCTION

	if b.q.filter != "" {
		queryBuf.WriteString("@filter(")
		queryBuf.WriteString(b.q.filter)
		queryBuf.WriteByte(')')
	}

	if b.q.recurse > 0 {
		queryBuf.WriteString("@recurse(depth:")
		queryBuf.Write(intToBytes(b.q.recurse))
		queryBuf.WriteByte(')')
	}

	queryBuf.WriteString(b.q.query)

	return queryBuf.String()
}
