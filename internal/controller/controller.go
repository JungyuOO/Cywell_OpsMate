package controller

const (
	OperatorName = "cywell-opsmate"
	Namespace    = "cywell-opsmate"
)

func StartupSummary() string {
	return "cywell-opsmate operator scaffold: lightspeed-api, postgres, console, aiops, rag"
}
