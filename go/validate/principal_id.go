//revive:disable:package-comments
package validate

import (
	"github.com/google/uuid"

	"github.com/pbrpc/connect-foundation/errors"
)

// PrincipalID validates a PrincipalID, which is a UUID in its canonical string
// form. The check runs before any query so a malformed id is the caller's
// error rather than the database's.
func PrincipalID(principalID string) (violations []errors.FieldViolation) {
	if principalID == "" {
		return append(violations, errors.FieldViolation{
			Field:       "principal_id",
			Description: "principal_id is required",
		})
	}

	if _, err := uuid.Parse(principalID); err != nil {
		violations = append(violations, errors.FieldViolation{
			Field:       "principal_id",
			Description: "principal_id must be a UUID",
		})
	}

	return
}
