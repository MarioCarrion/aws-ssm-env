// +build tools

package awsssmenv

// XXX When adding a new tool also update `scripts/tools.sh`.

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
)
