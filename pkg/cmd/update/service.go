package update

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdatePublishedServices(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-published-services ISSUE SERVICE [SERVICE...]",
		Aliases: []string{"ups"},
		Args:    cobra.MinimumNArgs(2),
		Short:   "Update published services of all associated issues",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cli, err := f.JiraClient()
			if err != nil {
				return
			}
			issue := args[0]
			issueKeys, err := GetLinkIssues(cli, issue)
			if err != nil {
				return
			}
			for _, issueKey := range issueKeys {
				err = UpdatePublishedService(cli, issueKey, args[1:]...)
				if err != nil {
					fmt.Println(issueKey, "failed: ", err)
				}
			}
			return
		},
	}
	return cmd
}

func GetLinkIssues(cli *jira.Client, issueID string) (issueKeys []string, err error) {
	issue, _, err := cli.Issue.Get(issueID, nil)
	if err != nil {
		return
	}
	for _, link := range issue.Fields.IssueLinks {
		if link.OutwardIssue == nil {
			continue
		}
		issueKeys = append(issueKeys, link.OutwardIssue.Key)
	}
	return
}

func UpdatePublishedService(cli *jira.Client, issueKey string, publishedServices ...string) (err error) {
	issue, _, err := cli.Issue.Get(issueKey, nil)
	if err != nil {
		return
	}

	if issue.Fields.Status.Name != "发布" {
		return
	}

	services, _ := issue.Fields.Unknowns[cmdutil.IssueFieldKeyServiceList].(string)
	var newServices strings.Builder
	var updated bool
LOOP:
	for _, s := range strings.Split(services, "\n") {
		svc := strings.TrimSpace(s)
		for _, ps := range publishedServices {
			if cmdutil.MatchServiceName(svc, ps) {
				if svc == ps {
					updated = true
					newServices.WriteString(svc + "（已发布）\n")
					continue LOOP
				}
			}
		}
		newServices.WriteString(svc + "\n")
	}

	if !updated {
		fmt.Println(issueKey, "nothing changed")
		return
	}
	fmt.Println(issueKey, "updated")
	err = UpdateIssueServiceList(cli, issueKey, newServices.String())
	return
}

func UpdateIssueServiceList(cli *jira.Client, issueKey string, newValue string) (err error) {
	update := map[string]any{
		"fields": map[string]string{
			cmdutil.IssueFieldKeyServiceList: newValue,
		},
	}
	resp, err := cli.Issue.UpdateIssue(issueKey, update)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}
