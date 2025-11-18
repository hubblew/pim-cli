package agents

type AgentTool interface {
	Descriptor() string
	ExecuteCommand(command string) (string, error)
}
