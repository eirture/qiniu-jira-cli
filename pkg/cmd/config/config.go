package config

import (
	"fmt"

	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "set the configs of jira and github",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cfg := f.Config()

			asks := []struct {
				prompt string
				value  *string
			}{
				{"Jira Address", &cfg.Jira.BaseURL},
				{"Jira Username", &cfg.Jira.Username},
				{"Jira Password", &cfg.Jira.Password},
			}
			for _, ask := range asks {
				err = cmdutil.AskVar(ask.prompt, ask.value)
				if err != nil {
					return
				}
			}
			cfgbytes, err := yaml.Marshal(cfg)
			if err != nil {
				return
			}
			ok, err := cmdutil.Ask(fmt.Sprintf("---\nConfig:\n%s\nIs this correct? (Y/n)", cfgbytes), cmdutil.WithAskDefaultValue("Y"))
			if err != nil || ok != "Y" {
				return
			}

			err = f.UpdateConfig(&cfg)
			return
		},
	}

	return cmd
}
