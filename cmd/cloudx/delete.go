package cloudx

import (
	"github.com/spf13/cobra"

	"github.com/ory/cli/cmd/cloudx/oauth2"
	"github.com/ory/cli/cmd/cloudx/relationtuples"

	"github.com/ory/cli/cmd/cloudx/client"
	"github.com/ory/cli/cmd/cloudx/identity"
	"github.com/ory/x/cmdx"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
	}

	cmd.AddCommand(
		identity.NewDeleteIdentityCmd(),
		oauth2.NewDeleteOAuth2Client(),
		oauth2.NewDeleteJWKs(),
		oauth2.NewDeleteAccessTokens(),
		relationtuples.NewDeleteCmd(),
	)

	client.RegisterConfigFlag(cmd.PersistentFlags())
	client.RegisterYesFlag(cmd.PersistentFlags())
	cmdx.RegisterNoiseFlags(cmd.PersistentFlags())
	cmdx.RegisterJSONFormatFlags(cmd.PersistentFlags())
	return cmd
}
