package arweave

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TransactionConnectionResponse struct {
	Data *Data         `json:"data"`
	Err  *GraphqlError `json:"error"`
}

type Data struct {
	Transactions *TransactionConnection `json:"transactions"`
}

type TransactionConnection struct {
	PageInfo *PageInfo          `json:"pageInfo"`
	Edges    []*TransactionEdge `json:"edges"`
}

type TransactionEdge struct {
	Cursor *string      `json:"cursor"`
	Node   *Transaction `json:"node"`
}

type PageInfo struct {
	HasNextPage *bool `json:"hasNextPage"`
}

type Transaction struct {
	ID   *string `json:"id"`
	Tags []*Tag  `json:"tags"`
}

type Tag struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

type GraphqlError struct {
	Errors []*MessageError `json:"errors"`
}

func (g *GraphqlError) Error() string {
	if len(g.Errors) == 0 {
		return "arweave: no errors"
	}
	msgs := make([]string, 0, len(g.Errors))
	for _, i := range g.Errors {
		if i != nil {
			msgs = append(msgs, i.Error())
		}
	}
	if len(msgs) == 0 {
		return "arweave: no errors"
	}
	return fmt.Sprintf("arweave: %v", strings.Join(msgs, "; "))
}

type MessageError struct {
	Message *string `json:"message"`
}

func (e *MessageError) Error() string {
	if e.Message == nil {
		return "no error message"
	}
	return *e.Message
}

type TagFilter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type TagFilters []TagFilter

func (t *TagFilter) GraphqlInput() string {
	values, _ := json.Marshal(t.Values)
	return fmt.Sprintf(`{name:"%v" values:%v}`, t.Name, string(values))
}

func (tags TagFilters) GraphqlInput() string {
	if len(tags) == 0 {
		return "[]"
	}
	b := new(strings.Builder)
	for _, tag := range tags {
		if b.Len() != 0 {
			b.WriteString(",")
		}
		b.WriteString(tag.GraphqlInput())
	}
	return fmt.Sprintf("[%s]", b.String())
}

type GraphQLRequest struct {
	Query string `json:"query"`
}
