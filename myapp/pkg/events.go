package pkg

// * CommitPushed is the event that is sent when a commit is pushed to GitHub
type CommitPushed struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	Author     string `json:"author"`
	OccurredOn string `json:"occurred_on"`
}

// * CommitDeployed is the event that is sent when a commit is deployed to an environment
type CommitDeployed struct {
	ID         string `json:"id"`
	Env        string `json:"env"`
	OccurredOn string `json:"occurred_on"`
}
