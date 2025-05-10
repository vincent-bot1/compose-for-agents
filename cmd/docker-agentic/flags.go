package main

type Flags struct {
	Project      string
	Network      string
	APIPort      string
	UIPort       string
	Config       string
	OpenAIAPIKey string
}

func (f *Flags) ContainerName(providerName string) string {
	return f.Project + "-" + providerName + "-" + f.Network
}

func (f *Flags) NetworkName() string {
	return f.Project + "_" + f.Network
}
