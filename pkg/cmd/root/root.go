package root

import (
	cmdconfig "github.com/eirture/qiniu-jira-cli/pkg/cmd/config"
	"github.com/eirture/qiniu-jira-cli/pkg/cmd/list"
	"github.com/eirture/qiniu-jira-cli/pkg/cmd/update"
	"github.com/eirture/qiniu-jira-cli/pkg/cmd/version"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "qiniujira",
		Short:         "qiniujira is a tool for managing jira issues",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		cmdconfig.NewCmd(f),
		list.NewCmd(f),
		update.NewCmdUpdatePublishedServices(f),
		version.NewCmd(f),
	)

	return cmd
}
