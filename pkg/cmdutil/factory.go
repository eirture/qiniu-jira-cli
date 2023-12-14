package cmdutil

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/andygrunwald/go-jira"
	"github.com/eirture/qiniu-jira-cli/pkg/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

var cfgPath = os.ExpandEnv("$HOME/.config/qiniujira/config.yaml")

type Factory interface {
	JiraClient() (*jira.Client, error)
	GithubClient(ctx context.Context) (*github.Client, error)

	UpdateConfig(cfg *config.Config) (err error)
	Config() config.Config
}

type factory struct {
	cfg *config.Config
}

func NewFactory() (f Factory, err error) {
	err = os.MkdirAll(filepath.Dir(cfgPath), 0755)
	if err != nil {
		err = fmt.Errorf("mkdir error: %w", err)
		return
	}
	file, err := os.OpenFile(cfgPath, os.O_CREATE, 0644)
	if err != nil {
		err = fmt.Errorf("open config file error: %w", err)
		return
	}
	defer file.Close()
	cfgbytes, err := io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("read config file error: %w", err)
		return
	}
	var cfg config.Config
	if len(cfgbytes) > 0 {
		err = yaml.Unmarshal(cfgbytes, &cfg)
		if err != nil {
			err = fmt.Errorf("decode config error: %w", err)
			return
		}
	}
	f = &factory{cfg: &cfg}
	return
}

func (f *factory) JiraClient() (cli *jira.Client, err error) {
	tr := jira.BasicAuthTransport{
		Username: f.cfg.Jira.Username,
		Password: f.cfg.Jira.Password,
	}
	return jira.NewClient(tr.Client(), f.cfg.Jira.BaseURL)
}

func (f *factory) GithubClient(ctx context.Context) (cli *github.Client, err error) {
	github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: f.cfg.Github.OAuthToken})))
	return
}

func (f *factory) Config() config.Config {
	return *f.cfg
}

func (f *factory) UpdateConfig(cfg *config.Config) (err error) {
	if reflect.DeepEqual(f.cfg, cfg) {
		return
	}
	file, err := os.OpenFile(cfgPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	err = yaml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return
	}
	f.cfg = cfg
	return
}
