package parser

import (
	"testing"

	"github.com/genjidb/genji/sql/query"
	"github.com/genjidb/genji/sql/query/expr"
	"github.com/stretchr/testify/require"
)

func TestParserInsert(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected query.Statement
		fails    bool
	}{
		{"Documents", `INSERT INTO test VALUES {a: 1, "b": "foo", c: 'bar', d: 1 = 1, e: {f: "baz"}}`,
			query.InsertStmt{
				TableName: "test",
				Values: expr.LiteralExprList{
					expr.KVPairs{
						expr.KVPair{K: "a", V: expr.IntValue(1)},
						expr.KVPair{K: "b", V: expr.TextValue("foo")},
						expr.KVPair{K: "c", V: expr.TextValue("bar")},
						expr.KVPair{K: "d", V: expr.Eq(expr.IntValue(1), expr.IntValue(1))},
						expr.KVPair{K: "e", V: expr.KVPairs{
							expr.KVPair{K: "f", V: expr.TextValue("baz")},
						}},
					},
				}}, false},
		{"Documents / Multiple", `INSERT INTO test VALUES {"a": 'a', b: -2.3}, {a: 1, d: true}`,
			query.InsertStmt{
				TableName: "test",
				Values: expr.LiteralExprList{
					expr.KVPairs{
						expr.KVPair{K: "a", V: expr.TextValue("a")},
						expr.KVPair{K: "b", V: expr.Float64Value(-2.3)},
					},
					expr.KVPairs{expr.KVPair{K: "a", V: expr.IntValue(1)}, expr.KVPair{K: "d", V: expr.BoolValue(true)}},
				},
			}, false},
		{"Documents / Positional Param", "INSERT INTO test VALUES ?, ?",
			query.InsertStmt{
				TableName: "test",
				Values:    expr.LiteralExprList{expr.PositionalParam(1), expr.PositionalParam(2)},
			},
			false},
		{"Documents / Named Param", "INSERT INTO test VALUES $foo, $bar",
			query.InsertStmt{
				TableName: "test",
				Values:    expr.LiteralExprList{expr.NamedParam("foo"), expr.NamedParam("bar")},
			},
			false},
		{"Values / With columns", "INSERT INTO test (a, b) VALUES ('c', 'd', 'e')",
			query.InsertStmt{
				TableName:  "test",
				FieldNames: []string{"a", "b"},
				Values: expr.LiteralExprList{
					expr.LiteralExprList{expr.TextValue("c"), expr.TextValue("d"), expr.TextValue("e")},
				},
			}, false},
		{"Values / Multiple", "INSERT INTO test (a, b) VALUES ('c', 'd'), ('e', 'f')",
			query.InsertStmt{
				TableName:  "test",
				FieldNames: []string{"a", "b"},
				Values: expr.LiteralExprList{
					expr.LiteralExprList{expr.TextValue("c"), expr.TextValue("d")},
					expr.LiteralExprList{expr.TextValue("e"), expr.TextValue("f")},
				},
			}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, err := ParseQuery(test.s)
			if test.fails {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, q.Statements, 1)
			require.EqualValues(t, test.expected, q.Statements[0])
		})
	}
}
