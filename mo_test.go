package gettext

import (
	"os"
	"testing"
)

func TestParseMO(t *testing.T) {
	r, err := os.Open("test.mo")
	if err != nil {
		t.Fatalf("Open failed with %+v", err)
	}
	if catalog, err := ParseMO(r); err != nil {
		t.Fatalf("ParseMO returned error: %+v", err)
	} else {
		if catalog.GetText("Untranslated String") != "Untranslated String" {
			t.Error("Catalog.GetText returned wrong result for untranslated string")
		}
		if catalog.GetText("Title") != "Título" {
			t.Error("Catalog.GetText returned wrong result for translated string")
		}

		if catalog.NGetText("Singular", "Plural", 1) != "Singular" {
			t.Error("Catalog.NGetText returned wrong result for untranslated singular string")
		}
		if catalog.NGetText("Singular", "Plural", 2) != "Plural" {
			t.Error("Catalog.NGetText returned wrong result for untranslated plural string")
		}
		if catalog.NGetText("%d topic", "%d topics", 1) != "%d tema" {
			t.Error("Catalog.NGetText returned wrong result for translated singular string")
		}
		if catalog.NGetText("%d topic", "%d topics", 4) != "%d temas" {
			t.Error("Catalog.NGetText returned wrong result for translated plural string")
		}
	}
	r.Close()
}
