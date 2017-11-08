package gettext

import "testing"

func TestModuloParser(t *testing.T) {
	if f := (&parser{exp: []byte("n")}).modulo(); f(111) != 111 {
		t.Fatalf("n modulo parsing failed. Expected 111 got %d", f(111))
	}
	if f := (&parser{exp: []byte("n%10")}).modulo(); f(12) != 2 {
		t.Fatalf("n modulo parsing failed. Expected 2 got %d", f(12))
	}
}

func TestFactorParser(t *testing.T) {
	if f := (&parser{exp: []byte("123")}).factor(); f(0) != 123 {
		t.Fatalf("Number factor parsing failed. Expected 123 got %d", f(0))
	}
	if f := (&parser{exp: []byte("n")}).factor(); f(111) != 111 {
		t.Fatalf("n factor parsing failed. Expected 111 got %d", f(111))
	}
}

func TestNumberParser(t *testing.T) {
	p := &parser{exp: []byte("123")}
	f := p.number()
	if f(0) != 123 {
		t.Fatalf("Number formula parsing failed. Expected 123 got %d", f(0))
	}
}
