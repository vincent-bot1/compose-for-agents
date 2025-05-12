package main

import "strings"

type Flags struct {
	Project          string
	Image            string
	Network          string
	Tools            string
	LogCalls         string // Should be a bool but compose provider mechanism doesn't like that
	ScanSecrets      string // Should be a bool but compose provider mechanism doesn't like that
	VerifySignatures string // Should be a bool but compose provider mechanism doesn't like that
}

func (f *Flags) ContainerName(providerName string) string {
	return f.Project + "-" + providerName + "-" + f.Network
}

func (f *Flags) NetworkName() string {
	return f.Project + "_" + f.Network
}

func (f *Flags) LogCallsEnabled() bool {
	return strings.EqualFold(f.LogCalls, "yes")
}

func (f *Flags) ScanSecretsEnabled() bool {
	return strings.EqualFold(f.ScanSecrets, "yes")
}

func (f *Flags) VerifySignaturesEnabled() bool {
	return strings.EqualFold(f.VerifySignatures, "yes")
}
