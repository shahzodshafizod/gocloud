package notifications

type Message struct {
	AgentID string `json:"agent_id"`
	Token   string `json:"token"`
	Title   string `json:"title"`
	Body    string `json:"body"`
}

type agent struct {
	ID       string
	Name     string
	Priority int
}
