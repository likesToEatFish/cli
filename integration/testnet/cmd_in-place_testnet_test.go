package testnet

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/xos"
	envtest "github.com/ignite/cli/v29/integration"
)

func TestServeWithCustomConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	var (
		env = envtest.New(t)
		app = env.Scaffold("github.com/test/sgbloga")
	)
	// Move config
	newConfig := "new_config.yml"
	newConfigPath := filepath.Join(tmpDir, newConfig)
	err := xos.Rename(filepath.Join(app.SourcePath(), "config.yml"), newConfigPath)
	require.NoError(t, err)
	app.SetConfigPath(newConfigPath)

	servers := app.RandomizeServerPorts()

	var (
		ctx, cancel       = context.WithTimeout(env.Ctx(), envtest.ServeTimeout)
		isBackendAliveErr error
	)
	go func() {
		defer cancel()
		isBackendAliveErr = env.IsAppServed(ctx, servers.API)
	}()
	fmt.Println("pppppppp1")
	env.Must(app.Serve("should serve", envtest.ExecCtx(ctx)))

	cmd := exec.Command("sh", "-c", "killall sgbloga || true")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to run killall command: %s", err)
	}

	fmt.Println("pppppppp")
	require.NoError(t, isBackendAliveErr, "app cannot get online in time")
}
