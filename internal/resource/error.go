package resource

type ErrorCode string

const (
	ErrorCodeX ErrorCode = "x"
)

type Error struct {
	Message string     `json:"message"`
	Code    *ErrorCode `json:"code"`
}
