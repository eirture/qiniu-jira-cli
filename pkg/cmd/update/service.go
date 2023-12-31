package update

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/eirture/qiniu-jira-cli/pkg/x/sets"
	"github.com/spf13/cobra"
)

func NewCmdUpdatePublishedServices(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update-published-services ISSUE SERVICE [SERVICE...]",
		Aliases: []string{"ups"},
		Args:    cobra.MinimumNArgs(2),
		Short:   "Update published services of all associated issues of the given deployment issue",
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
			if len(issueKeys) == 0 {
				fmt.Println("no linked issues")
				return
			}
			var svcs []string
			for _, issueKey := range issueKeys {
				svcs, err = UpdatePublishedService(cli, issueKey, args[1:]...)
				if err != nil {
					fmt.Printf("%s: %v\n", issueKey, err)
					continue
				}
				fmt.Printf("%s: [%s]\n", issueKey, strings.Join(svcs, ", "))
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
		var issue string
		if link.OutwardIssue != nil {
			issue = link.OutwardIssue.Key
		} else if link.InwardIssue != nil {
			issue = link.InwardIssue.Key
		} else {
			continue
		}
		issueKeys = append(issueKeys, issue)
	}
	return
}

func UpdatePublishedService(cli *jira.Client, issueKey string, publishedServices ...string) (updatedServices []string, err error) {
	issue, _, err := cli.Issue.Get(issueKey, nil)
	if err != nil {
		return
	}

	if issue.Fields.Status.Name != "发布" {
		return
	}

	pServices := sets.NewSet[string](publishedServices...)
	services, _ := issue.Fields.Unknowns[cmdutil.IssueFieldKeyServiceList].(string)
	var newServices strings.Builder
	for _, s := range strings.Split(services, "\n") {
		svc := strings.TrimSpace(s)
		if !cmdutil.IsPublishedService(svc) {
			psn := cmdutil.GetPureServiceName(svc)
			if pServices.Has(psn) {
				updatedServices = append(updatedServices, psn)
				// use the `svc`` to keep the original service name marks
				svc += cmdutil.ServicePublishedMark
			}
		}
		newServices.WriteString(svc + "\n")
	}

	if len(updatedServices) == 0 {
		return
	}
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
