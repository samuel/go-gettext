package gettext

import (
	"bytes"
	"fmt"
)

type parser struct {
	exp []byte
}

type parseError string

func (e parseError) Error() string {
	return string(e)
}

func parse(exp []byte) (pf PluralFormula, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(parseError); ok {
				pf = nil
				err = e
				return
			}
			panic(r)
		}
	}()

	p := &parser{
		exp: exp,
	}
	pf = p.factor()
	return pf, nil
}

func (p *parser) error(msg string, a ...interface{}) {
	panic(parseError(fmt.Sprintf(msg, a...)))
}

func (p *parser) trim() {
	for len(p.exp) > 0 && p.exp[0] == ' ' {
		p.exp = p.exp[1:]
	}
}

func (p *parser) popChar(c byte) bool {
	p.trim()
	if len(p.exp) > 0 && p.exp[0] == c {
		p.exp = p.exp[1:]
		return true
	}
	return false
}

func (p *parser) pop(b []byte) bool {
	p.trim()
	if len(p.exp) >= len(b) && bytes.Equal(b, p.exp[:len(b)]) {
		p.exp = p.exp[len(b):]
		return true
	}
	return false
}

func (p *parser) mustPopChar(c byte) {
	if !p.popChar(c) {
		p.error("Expected '%s'", string(c))
	}
}

func (p *parser) expression() PluralFormula {
	return nil
}

func (p *parser) inEquality() PluralFormula {
	f := p.modulo()
	switch {
	case p.pop([]byte("<=")):
		m := p.modulo()
		return func(n int) int {
			if f(n) <= m(n) {
				return 1
			}
			return 0
		}
	case p.pop([]byte(">=")):
		m := p.modulo()
		return func(n int) int {
			if f(n) >= m(n) {
				return 1
			}
			return 0
		}
	case p.pop([]byte("<")):
		m := p.modulo()
		return func(n int) int {
			if f(n) < m(n) {
				return 1
			}
			return 0
		}
	case p.pop([]byte(">")):
		m := p.modulo()
		return func(n int) int {
			if f(n) > m(n) {
				return 1
			}
			return 0
		}
	}
	return f
}

func (p *parser) modulo() PluralFormula {
	f := p.factor()
	if p.popChar('%') {
		mod := p.factor()
		return func(n int) int {
			return f(n) % mod(n)
		}
	}
	return f
}

func (p *parser) factor() PluralFormula {
	switch {
	case p.popChar('n'):
		return func(n int) int { return n }
	case p.popChar('('):
		e := p.expression()
		p.mustPopChar(')')
		return e
	}
	return p.number()
}

func (p *parser) number() PluralFormula {
	i := 0
	v := 0
	for len(p.exp) > i && (p.exp[i] >= '0' && p.exp[i] <= '9') {
		v2 := (v * 10) + int(p.exp[i]-'0')
		if v2 < v {
			p.error("Number overflow")
		}
		v = v2
		i++
	}
	if i == 0 {
		p.error("Expected a number")
	}
	p.exp = p.exp[i:]
	return func(n int) int { return v }
}
