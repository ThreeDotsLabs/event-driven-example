package pkg

type CommitPushed struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	Author     string `json:"author"`
	OccurredOn string `json:"occurred_on"`
}

type CommitDeployed struct {
	ID         string `json:"id"`
	Env        string `json:"env"`
	OccurredOn string `json:"occurred_on"`
}
