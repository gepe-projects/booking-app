package constant

// Postgres error codes (SQLSTATE)
const (
	// Integrity constraint violations
	PgErrUniqueViolation     = "23505" // duplicate key
	PgErrForeignKeyViolation = "23503" // foreign key violation
	PgErrNotNullViolation    = "23502" // null value in column that disallows nulls
	// PgErrCheckViolation      = "23514" // check constraint violation

	// Transaction errors
	PgErrSerializationFailure = "40001" // serialization failure (e.g., concurrent update)
	PgErrDeadlockDetected     = "40P01" // deadlock detected

	// Authentication / authorization
	PgErrInvalidPassword       = "28P01" // invalid_password
	PgErrInsufficientPrivilege = "42501" // insufficient_privilege
)
