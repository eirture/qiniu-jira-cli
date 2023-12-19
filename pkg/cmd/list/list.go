package list

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

const (
	IssueFieldKeyServiceList = "customfield_12300"
	BotName                  = "qiniu-bot"
)

var (
	nameAlphaPattern         = regexp.MustCompile(`^[a-zA-Z0-9_].*$`)
	qiniuBotPRCommentPattern = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/pull/([0-9]+)`)
)

func NewCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-deploying-issues SERVICE [SERVICE...]",
		Aliases: []string{"ldi"},
		Short:   "List all issues",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			jiraCli, err := f.JiraClient()
			if err != nil {
				return
			}
			ctx := context.Background()
			ghCli, err := f.GithubClient(ctx)
			if err != nil {
				return
			}
			issues, err := SearchPublishingIssues(jiraCli, ghCli, args)
			if err != nil {
				return
			}

			slices.SortFunc(issues, func(a, b Issue) int {
				return a.MergedAt.Compare(b.MergedAt)
			})
			for _, issue := range issues {
				fmt.Println(issue.Key, issue.MergedAt, issue.UnpublishedServices)
			}
			return
		},
	}
	return cmd
}

type Issue struct {
	Key                 string
	MergedAt            time.Time
	UnpublishedServices []string
}

func SearchPublishingIssues(cli *jira.Client, githubCli *github.Client, services []string) (results []Issue, err error) {
	var svcFilters []string
	for _, svc := range services {
		svcFilters = append(svcFilters, fmt.Sprintf("服务列表 ~ %s", svc))
	}
	jql := fmt.Sprintf(
		"project = KODO AND status = 发布 AND (%s) ORDER BY created DESC",
		strings.Join(svcFilters, " OR "),
	)
	issues, _, err := cli.Issue.Search(jql, &jira.SearchOptions{
		Fields: []string{"*key"},
		//MaxResults: 2,
	})
	if err != nil {
		return
	}
	for _, si := range issues {
		var issue *jira.Issue
		issue, _, err = cli.Issue.Get(si.Key, &jira.GetQueryOptions{})
		if err != nil {
			return
		}
		var mergedAt, lastedMergedAt time.Time
		var merged bool
		for _, cmt := range issue.Fields.Comments.Comments {
			if cmt.Author.DisplayName != BotName {
				continue
			}
			parts := qiniuBotPRCommentPattern.FindStringSubmatch(cmt.Body)
			if len(parts) != 4 {
				continue
			}
			owner, repo, pr := parts[1], parts[2], parts[3]
			prNumber, _ := strconv.Atoi(pr)
			merged, mergedAt, err = isPRMerged(githubCli, context.Background(), owner, repo, prNumber)
			if err != nil {
				fmt.Fprintf(os.Stderr, "check %s pr %s error: %v\n", issue.Key, parts[0], err)
				continue
			}
			if !merged {
				continue
			}
			if lastedMergedAt.IsZero() || lastedMergedAt.Before(mergedAt) {
				lastedMergedAt = mergedAt
			}
		}
		unpublishedServices := getUnpublishedServices(issue, services)
		if len(unpublishedServices) == 0 {
			continue
		}
		if !lastedMergedAt.IsZero() {
			results = append(results, Issue{
				Key:                 issue.Key,
				MergedAt:            lastedMergedAt,
				UnpublishedServices: unpublishedServices,
			})
		}
	}
	return
}

func getUnpublishedServices(issue *jira.Issue, filterServices []string) (results []string) {
	services, _ := issue.Fields.Unknowns[IssueFieldKeyServiceList].(string)
LOOP:
	for _, s := range strings.Split(services, "\n") {
		svc := strings.TrimSpace(s)
		for _, ps := range filterServices {
			if serviceNameMatched(svc, ps) {
				if svc == ps {
					results = append(results, ps)
					continue LOOP
				}
			}
		}
	}
	return
}

func serviceNameMatched(src, target string) bool {
	if !strings.HasPrefix(src, target) {
		return false
	}
	if len(src) > len(target) && nameAlphaPattern.MatchString(src[len(target):]) {
		return false
	}
	return true
}

func isPRMerged(cli *github.Client, ctx context.Context, owner, repo string, number int) (merged bool, mergedAt time.Time, err error) {
	pr, _, err := cli.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		fmt.Println("get pr error:", err)
		return
	}
	fmt.Println("get pr:", number, pr.State)
	merged = pr.GetMerged()
	mergedAt = pr.GetMergedAt()
	return
}
