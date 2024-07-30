package exceptions

import (
	"fmt"
)

const prefix = "channels-api"

var (

	// Errors related to request validation
	ErrInvalidPayload       = fmt.Errorf("%s: invalid payload", prefix)
	ErrChannelAlreadyExists = fmt.Errorf("%s: channel already exists", prefix)
	ErrInvalidNameField     = fmt.Errorf("%s: invalid name field", prefix)
	ErrInvalidMembersField  = fmt.Errorf("%s: invalid members field", prefix)
	ErrInvalidAdminsField   = fmt.Errorf("%s: invalid admins field", prefix)
	ErrInvalidID            = fmt.Errorf("%s: invalid ID", prefix)
	ErrInvalidUserIdSent    = fmt.Errorf("%s: invalid user ID sent", prefix)

	// Database related errors
	ErrChannelNotFound = fmt.Errorf("%s: channel not found", prefix)
	ErrDatabaseFailure = fmt.Errorf("%s: database failure", prefix)
)
