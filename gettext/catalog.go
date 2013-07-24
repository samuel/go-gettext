package gettext

import (
	"net/textproto"
)

type PluralFormula func(int) int

type Translation struct {
	Plural      string
	Translation []string
}

type Catalog struct {
	Header  textproto.MIMEHeader
	Strings map[string]*Translation

	// TODO: Parsed from "Plural-Forms" header (e.g. 'Plural-Forms: nplurals=2; plural=n != 1;')
	// See: http://www.gnu.org/savannah-checkouts/gnu/gettext/manual/html_node/Plural-forms.html
	// NumPlurals    int
	PluralFormula PluralFormula
}

var (
	// The NullCatalog has no translations so can be used as a fallback
	NullCatalog = &Catalog{
		Header:        make(map[string][]string),
		Strings:       map[string]*Translation{},
		PluralFormula: GermanicPluralFormula,
	}
)

func GermanicPluralFormula(n int) int {
	if n == 1 {
		return 0
	}
	return 1
}

func (c *Catalog) GetText(text string) string {
	if t := c.Strings[text]; t != nil && len(t.Translation) > 0 {
		return t.Translation[0]
	}
	return text
}

func (c *Catalog) NGetText(singular string, plural string, n int) string {
	fallback := singular
	if n != 1 {
		fallback = plural
	}

	t := c.Strings[singular]
	if t == nil {
		return fallback
	}

	pn := c.PluralFormula(n)
	if pn < len(t.Translation) {
		return t.Translation[pn]
	}

	return fallback
}
