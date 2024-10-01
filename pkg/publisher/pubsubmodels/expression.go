package pubsubmodels

type Event struct {
	EventId    string `json:"eventId,omitempty"`
	Expression string `json:"expression,omitempty"`
}
