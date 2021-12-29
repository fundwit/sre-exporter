package fail

import "errors"

var ErrUnexpectedDatabase = errors.New("unexpected database")
var ErrInvalidDatabaseUrl = errors.New("invalid mysql driver args")
