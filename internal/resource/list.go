package resource

type List[E any] struct {
	Object string `json:"object"`
	Data   []E    `json:"data"`
}
