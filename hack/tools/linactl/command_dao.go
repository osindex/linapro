// This file implements the dao command for GoFrame DAO generation.

package main

import (
	"context"

	"linactl/internal/goframecli"
)

// runDao runs the embedded GoFrame gen dao command in the core application
// directory without requiring an external gf binary.
func runDao(ctx context.Context, a *app, _ commandInput) error {
	return goframecli.Run(ctx, a.root, a.executable, a.runCommand, "gen", "dao")
}
