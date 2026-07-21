package result

type Item struct {
	Image string `json:"image"`

	Source string `json:"source"`

	Target string `json:"target"`

	Status string `json:"status"`

	Error string `json:"error,omitempty"`
}
