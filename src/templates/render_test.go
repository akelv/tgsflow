package templates

import "testing"

func TestRender_SpecTemplate(t *testing.T) {
	out, err := Render("thought/10_spec.md.tmpl", map[string]any{"Role": "developer", "Outcome": "add features"})
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	if want := "developer"; !contains(out, want) {
		t.Fatalf("expected output to contain %q, got: %s", want, out)
	}
}

func contains(hay, needle string) bool { return len(hay) >= len(needle) && (index(hay, needle) >= 0) }

func index(hay, needle string) int {
	// small helper to avoid importing strings in tiny test
	for i := 0; i+len(needle) <= len(hay); i++ {
		if hay[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
