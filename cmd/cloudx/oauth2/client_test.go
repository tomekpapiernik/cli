package oauth2_test

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid/v3"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/testhelpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCreateClient(t *testing.T) {
	t.Run("is not able to create client if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "create", "client", "--quiet", "--project", defaultProject)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to create client", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "create", "client", "--format", "json", "--project", defaultProject)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		t.Logf("Created client: %s", stdout)
	})
}

func TestDeleteClient(t *testing.T) {
	t.Run("is not able to delete oauth2 client if not authenticated and quiet flag", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "delete", "oauth2-client", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to delete oauth2 client", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		stdout, stderr, err := defaultCmd.Exec(nil, "delete", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, userID, out.String(), "stdout: %s", stdout)
	})

	t.Run("is able to delete oauth2 client after authenticating", func(t *testing.T) {
		userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "delete", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Equal(t, userID, out.String(), stdout)
	})
}

func TestGetClient(t *testing.T) {
	userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()

	t.Run("is not able to get oauth2 if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "get", "oauth2-client", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to get oauth2", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "get", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})

	t.Run("is able to get oauth2 after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "get", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})
}

func TestImportClient(t *testing.T) {
	t.Run("is not able to import oauth2-client if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "import", "oauth2-client", "--quiet", "--project", defaultProject)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to import oauth2-client", func(t *testing.T) {
		name := uuid.Must(uuid.NewV4()).String()
		stdout, stderr, err := defaultCmd.Exec(nil, "import", "oauth2-client", "--format", "json", "--project", defaultProject, testhelpers.MakeRandomClient(t, name))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, name, out.Get("client_name").String())
	})

	t.Run("is able to import oauth2-client after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		name := uuid.Must(uuid.NewV4()).String()
		stdout, stderr, err := cmd.Exec(r, "import", "oauth2-client", "--format", "json", "--project", defaultProject, testhelpers.MakeRandomClient(t, name))
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Equal(t, name, out.Get("client_name").String())
	})
}

func TestListClients(t *testing.T) {
	project := testhelpers.CreateProject(t, defaultConfig)

	userID := testhelpers.CreateClient(t, defaultCmd, project).Get("client_id").String()

	t.Run("is not able to list oauth2 clients if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "list", "oauth2-clients", "--quiet", "--project", project)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	for _, proc := range []string{"list", "ls"} {
		t.Run(fmt.Sprintf("is able to %s oauth2 clients", proc), func(t *testing.T) {
			stdout, stderr, err := defaultCmd.Exec(nil, proc, "oauth2-clients", "--format", "json", "--project", project)
			require.NoError(t, err, stderr)
			out := gjson.Parse(stdout).Get("items")
			assert.True(t, gjson.Valid(stdout))
			assert.Len(t, out.Array(), 1)
			assert.Equal(t, userID, out.Array()[0].Get("client_id").String(), "%s", out)
		})
	}

	t.Run("is able to list oauth2 clients after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "ls", "oauth2-clients", "--format", "json", "--project", project)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout).Get("items")
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String(), "%s", out)
	})
}

func TestUpdateOAuth2(t *testing.T) {
	userID := testhelpers.CreateClient(t, defaultCmd, defaultProject).Get("client_id").String()

	t.Run("is not able to update oauth2 if not authenticated and quiet flag", func(t *testing.T) {
		configDir := testhelpers.NewConfigDir(t)
		cmd := testhelpers.ConfigAwareCmd(configDir)
		_, _, err := cmd.Exec(nil, "update", "oauth2-client", "--quiet", "--project", defaultProject, userID)
		require.ErrorIs(t, err, client.ErrNoConfigQuiet)
	})

	t.Run("is able to update oauth2", func(t *testing.T) {
		stdout, stderr, err := defaultCmd.Exec(nil, "update", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		out := gjson.Parse(stdout)
		assert.True(t, gjson.Valid(stdout))
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})

	t.Run("is able to update oauth2 after authenticating", func(t *testing.T) {
		cmd, r := testhelpers.WithReAuth(t, defaultEmail, defaultPassword)
		stdout, stderr, err := cmd.Exec(r, "update", "oauth2-client", "--format", "json", "--project", defaultProject, userID)
		require.NoError(t, err, stderr)
		assert.True(t, gjson.Valid(stdout))
		out := gjson.Parse(stdout)
		assert.Len(t, out.Array(), 1)
		assert.Equal(t, userID, out.Array()[0].Get("client_id").String())
	})
}
