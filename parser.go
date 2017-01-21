package parser

import (
	"fmt"
	"log"
	"regexp"
)

type parseError struct {
	Location Location
	Parser   Parser
	Child    *parseError
}

func (me parseError) Error() string {
	if me.Child == nil {
		return fmt.Sprintf("error parsing %v at %d", me.Parser, me.Location)
	} else {
		return fmt.Sprintf("%s while parsing %v at %d", me.Child.Error(), me.Parser, me.Location)
	}
}

type Location int64

type Parser interface {
	Parse(Context)
}

func Parse(p Parser, s string) (err error) {
	c := context{input: s}
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err = r.(error)
	}()
	c.Parse(p)
	return
}

type Regexp struct {
	pattern string
	matches []string
}

func (me *Regexp) Parse(c Context) {
	r := regexp.MustCompile("^" + me.pattern)
	log.Println(r, c.StreamString())
	me.matches = r.FindStringSubmatch(c.StreamString())
	if me.matches == nil {
		c.Fail()
	}
	c.Advance(len(me.matches[0]))
}

type String string

func (me String) Parse(c Context) {
	for i := range me {
		if c.Byte() != me[i] {
			c.Fail()
		}
		c.Advance(1)
	}
}
