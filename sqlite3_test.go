package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// http://www.sqlite.org/draft/tokenreq.html

type whitespace struct{}

func (me whitespace) Parse(c Context) {
	c.Parse(&Regexp{pattern: `\s*`})
}

type keyword string

func (me keyword) Parse(c Context) {
	re := Regexp{pattern: `([a-z]+)`}
	c.Parse(&re)
	if strings.ToLower(re.matches[1]) != strings.ToLower(string(me)) {
		c.Fail()
	}
}

type selectStmt struct {
	resultColumns []resultColumn
	table         table
}

func (me *selectStmt) Parse(c Context) {
	k := keyword("select")
	c.Parse(&k)
	for {
		c.Parse(whitespace{})
		var rc resultColumn
		c.Parse(&rc)
		me.resultColumns = append(me.resultColumns, rc)
		c.Parse(whitespace{})
		if !c.ParseOk(String(",")) {
			break
		}
	}
	c.Parse(whitespace{})
	c.Parse(keyword("from"))
	c.Parse(whitespace{})
	c.Parse(&me.table)
}

type resultColumn struct {
	ident
}

type table struct {
	ident
}

type ident string

func (me *ident) Parse(c Context) {
	re := Regexp{pattern: `([a-zA-Z]+)`}
	c.Parse(&re)
	*me = ident(re.matches[1])
}

func TestParseSQLite3Select(t *testing.T) {
	var ss selectStmt
	err := Parse(&ss, "select blah from herp")
	require.NoError(t, err)
	require.Len(t, ss.resultColumns, 1)
	assert.EqualValues(t, "blah", ss.resultColumns[0].ident)
	assert.EqualValues(t, "herp", ss.table.ident)
	ss = selectStmt{}
	err = Parse(&ss, "select herp, derp from blah")
	require.NoError(t, err)
	require.Len(t, ss.resultColumns, 2)
	assert.EqualValues(t, "herp", ss.resultColumns[0].ident)
	assert.EqualValues(t, "derp", ss.resultColumns[1].ident)
	assert.EqualValues(t, "blah", ss.table.ident)
	ss = selectStmt{}
	err = Parse(&ss, "select herp,, derp from blah")
	require.Error(t, err)
	t.Log(err)
}
