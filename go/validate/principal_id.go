//revive:disable:package-comments
package validate

import (
	"github.com/google/uuid"

	errors "github.com/pbrpc/connect-errors"
)

// PrincipalID validates a PrincipalID as an RFC 4122 variant UUIDv4.
// The check runs before any query so a malformed id is the caller's
// error rather than the database's.
func PrincipalID(principalID string) (violations []errors.FieldViolation) {
	if principalID == "" {
		return append(violations, errors.FieldViolation{
			Field:       "principal_id",
			Description: "principal_id is required",
		})
	}

	id, err := uuid.Parse(principalID)
	if err != nil || id.Version() != 4 || id.Variant() != uuid.RFC4122 {
		violations = append(violations, errors.FieldViolation{
			Field:       "principal_id",
			Description: "principal_id must be an RFC 4122 variant UUIDv4",
		})
	}

	return
}
