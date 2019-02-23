package pkg

type commitPushed struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	Author     string `json:"author"`
	OccurredOn string `json:"occurred_on"`
}

type commitDeployed struct {
	ID         string `json:"id"`
	Env        string `json:"env"`
	OccurredOn string `json:"occurred_on"`
}
