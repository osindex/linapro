// This file implements the hidden __goframe command for embedded GoFrame CLI dispatch.

package main

import (
	"context"
	"fmt"

	"linactl/internal/goframecli"
)

// runEmbeddedGoFrame executes a whitelisted embedded GoFrame code generation
// command. It is intentionally hidden from help output and only supports the
// code generation paths used by linactl ctrl and linactl dao.
func runEmbeddedGoFrame(ctx context.Context, a *app, input commandInput) error {
	if len(input.Params) > 0 {
		return fmt.Errorf("embedded GoFrame only supports positional commands: gen ctrl or gen dao")
	}
	return goframecli.RunEmbedded(ctx, a.root, input.Args...)
}
