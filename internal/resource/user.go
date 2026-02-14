package resource

type User struct {
	ID int64 `json:"id" minimum:"1"`
}
