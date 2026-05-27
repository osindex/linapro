// This file implements the ctrl command for GoFrame controller generation.

package main

import (
	"context"

	"linactl/internal/goframecli"
)

// runCtrl runs the embedded GoFrame gen ctrl command in the core application
// directory without requiring an external gf binary.
func runCtrl(ctx context.Context, a *app, _ commandInput) error {
	return goframecli.Run(ctx, a.root, a.executable, a.runCommand, "gen", "ctrl")
}
