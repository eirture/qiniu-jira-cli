package version

import (
	"fmt"

	"github.com/eirture/qiniu-jira-cli/pkg/build"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

const (
	format = `%s
  commit: %s
  go version: %s
`
)

func NewCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf(format, build.Version, build.Commit, build.GoVersion)
			return nil
		},
	}
	return cmd
}
