package resource

type List[E any] struct {
	Object        string  `json:"object"`
	Data          []E     `json:"data"`
	StartingAfter *string `json:"starting_after"`
	EndingBefore  *string `json:"ending_before"`
}
