// This file verifies linactl's embedded GoFrame CLI dispatch boundary.

package goframecli

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"linactl/internal/toolrun"
)

func TestRunDispatchesHiddenCommandInCoreDir(t *testing.T) {
	root := t.TempDir()
	binary := filepath.Join(root, "bin", "linactl")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "ctrl", args: []string{"gen", "ctrl"}},
		{name: "dao", args: []string{"gen", "dao"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				gotOptions toolrun.Options
				gotName    string
				gotArgs    []string
				calls      int
			)
			runner := func(_ context.Context, options toolrun.Options, name string, args ...string) error {
				calls++
				gotOptions = options
				gotName = name
				gotArgs = append([]string(nil), args...)
				return nil
			}

			err := Run(context.Background(), root, func() (string, error) {
				return binary, nil
			}, runner, tc.args...)
			if err != nil {
				t.Fatalf("Run returned error: %v", err)
			}
			if calls != 1 {
				t.Fatalf("expected one runner call, got %d", calls)
			}
			if gotName != binary {
				t.Fatalf("runner name mismatch: got %q want %q", gotName, binary)
			}
			expectedArgs := append([]string{hiddenCommand}, tc.args...)
			if !reflect.DeepEqual(gotArgs, expectedArgs) {
				t.Fatalf("runner args mismatch: got %#v want %#v", gotArgs, expectedArgs)
			}
			expectedDir := filepath.Join(root, "apps", "lina-core")
			if gotOptions.Dir != expectedDir {
				t.Fatalf("runner dir mismatch: got %q want %q", gotOptions.Dir, expectedDir)
			}
		})
	}
}

func TestValidateArgsAllowsOnlyCodeGenerationCommands(t *testing.T) {
	for _, args := range [][]string{
		{"gen", "ctrl"},
		{"gen", "dao"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			if err := validateArgs(args); err != nil {
				t.Fatalf("validateArgs should allow %v: %v", args, err)
			}
		})
	}

	for _, args := range [][]string{
		{"install"},
		{"build"},
		{"gen", "service"},
		{"gen", "ctrl", "extra"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			err := validateArgs(args)
			if err == nil || !strings.Contains(err.Error(), "embedded GoFrame only supports gen ctrl or gen dao") {
				t.Fatalf("expected whitelist error for %v, got %v", args, err)
			}
		})
	}
}

func TestRunRejectsUnsupportedCommandsBeforeExecution(t *testing.T) {
	for _, args := range [][]string{
		{"install"},
		{"build"},
		{"gen", "service"},
		{"gen", "ctrl", "extra"},
		{"dao"},
	} {
		t.Run(strings.Join(args, "_"), func(t *testing.T) {
			err := Run(context.Background(), t.TempDir(), func() (string, error) {
				t.Fatalf("executable resolver should not run for invalid args")
				return "", nil
			}, func(context.Context, toolrun.Options, string, ...string) error {
				t.Fatalf("runner should not run for invalid args")
				return nil
			}, args...)
			if err == nil || !strings.Contains(err.Error(), "embedded GoFrame only supports gen ctrl or gen dao") {
				t.Fatalf("expected whitelist error, got %v", err)
			}
		})
	}
}

func TestRunReportsExecutableResolutionFailure(t *testing.T) {
	expectedErr := errors.New("no executable")
	err := Run(context.Background(), t.TempDir(), func() (string, error) {
		return "", expectedErr
	}, func(context.Context, toolrun.Options, string, ...string) error {
		t.Fatalf("runner should not run when executable resolution fails")
		return nil
	}, "gen", "ctrl")
	if err == nil || !strings.Contains(err.Error(), "resolve linactl executable") {
		t.Fatalf("expected executable resolution error, got %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped executable error, got %v", err)
	}
}

func TestRunEmbeddedRejectsUnsupportedCommandsBeforeChangingDirectory(t *testing.T) {
	original, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("resolve current directory: %v", err)
	}
	err = RunEmbedded(context.Background(), filepath.Join(t.TempDir(), "missing-root"), "install")
	if err == nil || !strings.Contains(err.Error(), "embedded GoFrame only supports gen ctrl or gen dao") {
		t.Fatalf("expected whitelist error, got %v", err)
	}
	current, err := filepath.Abs(".")
	if err != nil {
		t.Fatalf("resolve current directory after RunEmbedded: %v", err)
	}
	if current != original {
		t.Fatalf("RunEmbedded changed directory before validating args: got %q want %q", current, original)
	}
}
