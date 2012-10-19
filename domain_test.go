package gettext

import "testing"

func TestDomain(t *testing.T) {
	d, err := NewDomain("test", "locale")
	if err != nil {
		t.Fatalf("NewDomain returned error %+v", err)
	}
	if c := d.GetCatalog("fallback"); c != NullCatalog {
		t.Fatalf("Domain.GetCatalog for unknown language didn't fallback to NullCatalog")
	}
	if c := d.GetCatalog("es"); c == NullCatalog {
		t.Fatalf("Domain.GetCatalog for known language returned NullCatalog")
	}
	if c := d.GetCatalog("en_US"); c == NullCatalog {
		t.Fatalf("Domain.GetCatalog for known language and country code returned NullCatalog")
	}
}
