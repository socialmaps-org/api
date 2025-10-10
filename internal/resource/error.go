package resource

type ErrorCode string

const (
	ErrorCodeTooLate ErrorCode = "too_late"
)

// Error Messages
const (
	ErrMsgInvalidPlaceID  string = "invalid place ID; a place ID starts with `plc_` prefix followed by numbers and lower- and upper-case letters"
	ErrMsgInvalidReviewID string = "invalid review ID; a review ID start with `rvw_` prefix followed by numbers and lower- and upper-case letters"
	ErrMsgNameTooLong     string = "name too long; a name must contain 256 characters or less"
	ErrMsgLatOutOfRange   string = "latitude is out of range; a latitude must be between -90 and +90 (both inclusive)"
	ErrMsgLonOutOfRange   string = "longitude is out of range; a longitude must be between -180 and +180 (both inclusive)"
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code,omitempty"`
}
