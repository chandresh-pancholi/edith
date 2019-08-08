package model

//Hit is
type Hit struct {
	ID     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
}

//HitsObject is a part of ES query response
type HitsObject struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
}

//ESResponse is ES search query response
type ESResponse struct {
	Took float32    `json:"took"`
	Hits HitsObject `json:"hits"`
}

//Total is
type Total struct {
	value int `json:"value"`
}
