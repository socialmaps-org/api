package resource

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type Review struct {
	Object  string    `json:"object" enum:"review"`
	ID      int64     `json:"id" minimum:"1" doc:"Unique identifier for this **Review**."`
	Created int64     `json:"created" minimum:"0" doc:"The [UNIX timestamp](https://en.wikipedia.org/wiki/Unix_time) of when this **Review** was created."`
	Place   PlaceStub `json:"place" doc:"The **Place** which this **Review** is about. This is a stub object, meaning it contains only an ID."`
	User    UserStub  `json:"user" doc:"The **User** who created this **Review**. This is a stub object, meaning it contains only an ID."`
	Liked   bool      `json:"liked" doc:"Whether the **User** liked this **Place** or not."`
	Comment string    `json:"comment" doc:"The comment written by the **User** about this **Place**, if written. Otherwise can be an empty string."`
	NLikes  uint64    `json:"n_likes" minimum:"0" doc:"The number of users who liked this **Review**."`
}

type ReviewWithUser struct {
	Review
	User User `json:"user" doc:"The **User** who created this **Review**. This is an expanded object, meaning it contains all fields."`
}

// Inline ReviewWithUser instead of using refs to avoid listing it under Schemas.
func (ReviewWithUser) Schema(r huma.Registry) *huma.Schema {
	type raw ReviewWithUser
	return huma.SchemaFromType(r, reflect.TypeOf(raw{}))
}
