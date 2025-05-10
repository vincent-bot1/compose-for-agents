package compose

var LabelNames = struct {
	Project         string
	Service         string
	OneOff          string
	ContainerNumber string
	ConfigHash      string
}{
	Project:         "com.docker.compose.project",
	Service:         "com.docker.compose.service",
	OneOff:          "com.docker.compose.oneoff",
	ContainerNumber: "com.docker.compose.container-number",
	ConfigHash:      "com.docker.compose.config-hash",
}
