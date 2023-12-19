package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/cli/oauth"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	githubHost    = "https://github.com"
	oauthClientID = "0c3d984a91ff98827c08"
	// This value is safe to be embedded in version control
	oauthClientSecret = ""
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
				opts   []cmdutil.AskOption
			}{
				{"Jira Address", &cfg.Jira.BaseURL, nil},
				{"Jira Username", &cfg.Jira.Username, nil},
				{"Jira Password", &cfg.Jira.Password, []cmdutil.AskOption{cmdutil.WithAskPassword(true)}},
			}
			for _, ask := range asks {
				err = cmdutil.AskVar(ask.prompt, ask.value, ask.opts...)
				if err != nil {
					return
				}
			}
			err = cmdutil.AskVar("Github OAuth Token", &cfg.Github.OAuthToken, cmdutil.WithAskPassword(true))
			if err != nil {
				return
			}
			//createNewGithubToken := true
			//if cfg.Github.OAuthToken != "" {
			//	var answer string
			//	answer, err = cmdutil.Ask(fmt.Sprintf("\nGithub OAuth Token: %s\nIs this correct? (Y/n)", cfg.Github.OAuthToken), cmdutil.WithAskDefaultValue("Y"))
			//	if err != nil {
			//		return
			//	}
			//	createNewGithubToken = strings.ToUpper(answer) != "Y"
			//}
			//if createNewGithubToken {
			//	var token string
			//	token, err = githubOAuthFlow()
			//	if err != nil {
			//		return
			//	}
			//	cfg.Github.OAuthToken = token
			//}
			cfgbytes, err := yaml.Marshal(cfg)
			if err != nil {
				return
			}
			ok, err := cmdutil.Ask(fmt.Sprintf("---\nConfig:\n%s\nIs this correct? (Y/n)", cfgbytes), cmdutil.WithAskDefaultValue("Y"))
			if err != nil || strings.ToUpper(ok) != "Y" {
				return
			}

			err = f.UpdateConfig(&cfg)
			return
		},
	}

	return cmd
}

func githubOAuthFlow() (token string, err error) {
	flow := &oauth.Flow{
		Host:         oauth.GitHubHost(githubHost),
		Scopes:       []string{"repo"},
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		CallbackURI:  "http://127.0.0.1/callback",
		DisplayCode: func(code, verificationURL string) error {
			fmt.Fprintf(os.Stderr, "First copy your one-time code: %s\n", code)
			return nil
		},
	}

	accessToken, err := flow.DetectFlow()
	if err != nil {
		return
	}

	token = accessToken.Token
	return
}
