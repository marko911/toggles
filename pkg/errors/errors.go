package errors

import "errors"

var (

	// ErrJSONPayloadEmpty is returned when the JSON payload is empty.
	ErrJSONPayloadEmpty = errors.New("JSON payload is empty")

	// ErrJSONPayloadInvalidBody is returned when body cannot be parsed
	ErrJSONPayloadInvalidBody = errors.New("Cannot parse request body")

	// ErrJSONPayloadInvalidFormat is returned when the JSON payload is invalid
	ErrJSONPayloadInvalidFormat = errors.New("Invalid JSON format")

	//ErrJSONPayloadInvalidFlag is returned when the flag sent is in bad format
	ErrJSONPayloadInvalidFlag = errors.New("Invalid JSON flag format")

	//ErrFailedCreateFlag is returned when flag creation fails
	ErrFailedCreateFlag = errors.New("Error creating flag")

	//ErrEvalRequestMissingFlag is returned when eval request has no flag id
	ErrEvalRequestMissingFlag = errors.New("FlagID is required")

	//ErrCantCastUser returns when invalid user data is passed to evaluate
	ErrCantCastUser = errors.New("cannot cast user from request")

	//ErrEvalRequestMissingUser is returned when user field is missing from request
	ErrEvalRequestMissingUser = errors.New("User field is required")

	//ErrFlagNotFound is returned when flag is not in database
	ErrFlagNotFound = errors.New("Flag not found, invalid key")
	//SuccessFlagCreated is message returned on success of flag post
	SuccessFlagCreated = "Flag created successfully"
)
