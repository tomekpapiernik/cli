package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/cli/cmd/cloudx/relationtuples"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/project"
	"github.com/ory/x/cmdx"
)

func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List resources",
	}

	cmd.AddCommand(
		project.NewListProjectsCmd(),
		identity.NewListIdentityCmd(),
		oauth2.NewListOAuth2Clients(),
		relationtuples.NewListCmd(),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
