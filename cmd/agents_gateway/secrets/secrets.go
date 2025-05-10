package secrets

import (
	"errors"
	"io"
	"log/slog"
	"sync"

	"github.com/aquasecurity/trivy/pkg/fanal/secret"
	tlog "github.com/aquasecurity/trivy/pkg/log"
)

var ErrSecretsFound = errors.New("found secrets")

var scanner = sync.OnceValue(func() *secret.Scanner {
	// Disable trivy logging
	// Create a slog handler that discards all output
	discardHandler := slog.NewTextHandler(io.Discard, nil)
	logger := slog.New(discardHandler)
	tlog.SetDefault(logger)

	scanner := secret.NewScanner(nil)
	return &scanner
})

func ContainsSecrets(text string) bool {
	result := scanner().Scan(secret.ScanArgs{FilePath: "argument", Content: []byte(text)})
	return len(result.Findings) > 0
}
