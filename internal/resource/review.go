package resource

type Review struct {
	Object  string `json:"object" enum:"review"`
	ID      int64  `json:"id" minimum:"1"`
	Created int64  `json:"created" minimum:"0"`
	Place   Place  `json:"place"`
	User    User   `json:"user"`
	Liked   bool   `json:"liked"`
	Comment string `json:"comment"`
	NLikes  uint64 `json:"n_likes" minimum:"0"`
}
