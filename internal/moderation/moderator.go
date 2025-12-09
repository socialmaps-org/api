package moderation

type Moderator interface {
	ID() string
	Moderate(review string) (*Decision, error)
}
