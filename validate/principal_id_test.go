package validate

import "testing"

func TestPrincipalID(t *testing.T) {
	t.Run("accepts an identifier", func(t *testing.T) {
		if violations := PrincipalID("01234567-89ab-4def-8123-456789abcdef"); len(violations) > 0 {
			t.Errorf("expected no violations, got %v", violations)
		}
	})

	for name, id := range map[string]string{
		"version 1":         "01234567-89ab-1def-8123-456789abcdef",
		"version 7":         "01234567-89ab-7def-8123-456789abcdef",
		"nil UUID":          "00000000-0000-0000-0000-000000000000",
		"NCS variant":       "01234567-89ab-4def-0123-456789abcdef",
		"Microsoft variant": "01234567-89ab-4def-c123-456789abcdef",
		"reserved variant":  "01234567-89ab-4def-e123-456789abcdef",
	} {
		t.Run("refuses "+name, func(t *testing.T) {
			violations := PrincipalID(id)
			if len(violations) != 1 {
				t.Fatalf("expected one violation, got %v", violations)
			}
			if violations[0].Field != "principal_id" {
				t.Errorf("expected the violation to name principal_id, got %q", violations[0].Field)
			}
		})
	}

	t.Run("refuses an empty identifier", func(t *testing.T) {
		violations := PrincipalID("")

		if len(violations) != 1 {
			t.Fatalf("expected one violation, got %v", violations)
		}

		if violations[0].Field != "principal_id" {
			t.Errorf("expected the violation to name principal_id, got %q", violations[0].Field)
		}
	})

	t.Run("refuses an identifier that is not a UUID", func(t *testing.T) {
		violations := PrincipalID("spaghetti")

		if len(violations) != 1 {
			t.Fatalf("expected one violation, got %v", violations)
		}

		if violations[0].Field != "principal_id" {
			t.Errorf("expected the violation to name principal_id, got %q", violations[0].Field)
		}
	})
}
