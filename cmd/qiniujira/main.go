package main

import (
	"fmt"
	"os"

	"github.com/eirture/qiniu-jira-cli/pkg/cmd/root"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
)

func main() {
	f, err := cmdutil.NewFactory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	cmd := root.NewCmd(f)

	expandedArgs := []string{}
	if len(os.Args) > 0 {
		expandedArgs = os.Args[1:]
	}

	cmd.SetArgs(expandedArgs)
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
