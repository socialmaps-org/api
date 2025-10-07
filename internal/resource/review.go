package resource

type Review struct {
	Object  string `json:"object"`
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Place   Place  `json:"place"`
	User    User   `json:"user"`
	Liked   bool   `json:"liked"`
	Comment string `json:"comment"`
	NLikes  uint64 `json:"n_likes"`
}
