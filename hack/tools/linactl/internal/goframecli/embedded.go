// This file executes the embedded GoFrame CLI in the hidden child command.

package goframecli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/cmd/gf/v2/gfcmd"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// RunEmbedded executes a whitelisted GoFrame code generation command in the
// current process. Callers should use it only from linactl's hidden child
// command so GoFrame CLI exits cannot terminate the parent linactl process.
func RunEmbedded(ctx context.Context, root string, args ...string) error {
	if err := validateArgs(args); err != nil {
		return err
	}
	if err := os.Chdir(coreDir(root)); err != nil {
		return fmt.Errorf("enter GoFrame project directory: %w", err)
	}
	if err := configureGoFrameCLI(root); err != nil {
		return err
	}
	command, err := gfcmd.GetCommand(ctx)
	if err != nil {
		return fmt.Errorf("initialize embedded GoFrame CLI: %w", err)
	}
	if command == nil {
		return fmt.Errorf("initialize embedded GoFrame CLI: command is nil")
	}
	_, err = command.RunWithSpecificArgs(ctx, append([]string{"gf"}, args...))
	if err != nil {
		return fmt.Errorf("run embedded GoFrame %s %s: %w", args[0], args[1], err)
	}
	return nil
}

// configureGoFrameCLI mirrors the gf CLI's hack/config.yaml lookup without
// calling gfcmd.Command.Run, whose error path can exit the process.
func configureGoFrameCLI(root string) error {
	adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
	if !ok {
		return nil
	}
	configDir := filepath.Join(coreDir(root), "hack")
	if err := adapter.SetPath(configDir); err != nil {
		return fmt.Errorf("configure embedded GoFrame CLI: %w", err)
	}
	return nil
}
