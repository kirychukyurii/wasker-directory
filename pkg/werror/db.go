package werror

import "errors"

var (
	ErrDatabaseRecordNotFound = errors.New("database: record not found")
	ErrDatabaseInternalError  = errors.New("database: internal error")
	ErrDatabaseQueryRow       = errors.New("database: querying rows")
	ErrDatabaseBuildSql       = errors.New("database: building SQL statement")
)

var (
	ErrBuildQueryReason = "BUILD_QUERY"
	ErrExecQueryReason  = "EXEC_QUERY"
)
