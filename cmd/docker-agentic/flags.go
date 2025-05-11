package main

const (
	agentsContainerPrefix = "agents"
	uiContainerPrefix     = "ui"
)

type Flags struct {
	Project      string
	Network      string
	APIPort      string
	UIPort       string
	Config       string
	OpenAIAPIKey string
}

func (f *Flags) AgentsContainerName(providerName string) string {
	return f.Project + "-" + agentsContainerPrefix + "-" + providerName + "-" + f.Network
}

func (f *Flags) UIContainerName(providerName string) string {
	return f.Project + "-" + uiContainerPrefix + "-" + providerName + "-" + f.Network
}

func (f *Flags) NetworkName() string {
	return f.Project + "_" + f.Network
}
