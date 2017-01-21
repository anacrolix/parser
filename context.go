package parser

type Context interface {
	Parse(Parser)
	Fail()
	StreamString() string
	Advance(int)
	ParseOk(Parser) bool
	Byte() byte
	Location() Location
}

var fail = new(byte)

type context struct {
	input string
	loc   Location
}

func (me context) Copy() Context {
	return &me
}

func (me *context) Location() Location {
	return me.loc
}

func (me *context) Advance(n int) {
	me.input = me.input[n:]
	me.loc += Location(n)
}

func (me *context) Parse(p Parser) {
	c := me.Copy()
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if r == fail {
			panic(parseError{
				Location: c.Location(),
				Parser:   p,
			})
		}
		if pe, ok := r.(parseError); ok {
			panic(parseError{
				Location: c.Location(),
				Parser:   p,
				Child:    &pe,
			})
		}
		panic(r)
	}()
	p.Parse(me)
}

func (me *context) Byte() byte {
	return me.input[0]
}

func (me *context) Fail() {
	panic(fail)
}

func (me *context) StreamString() string {
	return me.input
}

func (me *context) ParseOk(p Parser) bool {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if r == fail {
			return
		}
		panic(r)
	}()
	p.Parse(me)
	return true
}
