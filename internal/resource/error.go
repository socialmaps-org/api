package resource

type ErrorCode string

const (
	ErrorCodeTooLate ErrorCode = "too_late"
)

// Error Messages
const (
	ErrMsgInvalidPlaceID         string = "invalid place ID; a place ID starts with `plc_` prefix followed by numbers and lower- and upper-case letters"
	ErrMsgInvalidReviewID        string = "invalid review ID; a review ID start with `rvw_` prefix followed by numbers and lower- and upper-case letters"
	ErrMsgNameTooLong            string = "name too long; a name must contain 256 characters or less"
	ErrMsgLatOutOfRange          string = "latitude is out of range; a latitude must be between -90 and +90 (both inclusive)"
	ErrMsgLonOutOfRange          string = "longitude is out of range; a longitude must be between -180 and +180 (both inclusive)"
	ErrMsgLimitTooBig            string = "limit too big; a limit must be between 1 and 100 (both inclusive)"
	ErrMsgLimitZero              string = "limit is zero; a limit must be between 1 and 100 (both inclusive)"
	ErrMsgBeforeAfterBothPresent string = "starting_after and ending_before are both present; those are mutually exclusive"
)

type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code,omitempty"`
}
