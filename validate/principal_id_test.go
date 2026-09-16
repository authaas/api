package validate

import "testing"

func TestPrincipalID(t *testing.T) {
	t.Run("accepts an identifier", func(t *testing.T) {
		if violations := PrincipalID("01234567-89ab-cdef-0123-456789abcdef"); len(violations) > 0 {
			t.Errorf("expected no violations, got %v", violations)
		}
	})

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
