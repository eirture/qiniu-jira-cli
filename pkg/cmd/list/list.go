package list

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/eirture/qiniu-jira-cli/pkg/cmdutil"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
)

const (
	BotName = "qiniu-bot"
)

var (
	qiniuBotPRCommentPattern = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)/pull/([0-9]+)`)
	tomorrowNow              = time.Now().Add(24 * time.Hour)
)

func NewCmd(f cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list-deploying-issues SERVICE [SERVICE...]",
		Aliases: []string{"ldi"},
		Short:   "List all deploying issues of specified services",
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
			data := [][]string{}
			for _, issue := range issues {
				mergedAt := "Unknown"
				if !issue.MergedAt.Equal(tomorrowNow) {
					mergedAt = issue.MergedAt.Local().Format("2006-01-02 15:04:05 -07:00")
				}
				data = append(data, []string{
					issue.Key,
					mergedAt,
					strings.Join(issue.UnpublishedServices, ", "),
				})
				// fmt.Println(issue.Key, issue.MergedAt, issue.UnpublishedServices)
			}
			cmdutil.WriteTable(data, "Issue", "Merged At", "Unpublished Services")
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
	if err != nil || len(issues) == 0 {
		return
	}

	workers := len(issues)
	var wg sync.WaitGroup
	ch := make(chan string, len(issues))
	rch := make(chan *Issue, len(issues))
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for issueKey := range ch {
				issue, _, err := cli.Issue.Get(issueKey, &jira.GetQueryOptions{})
				if err != nil {
					fmt.Fprintf(os.Stderr, "get %s error: %v\n", issueKey, err)
					continue
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
					merged, mergedAt, err = cmdutil.IsPRMerged(githubCli, context.Background(), owner, repo, prNumber)
					if err != nil {
						fmt.Fprintf(os.Stderr, "check %s pr %s error: %v\n", issue.Key, parts[0], err)
						// assume it has been merged, and merged in the future
						merged = true
						mergedAt = tomorrowNow
					}
					if !merged {
						continue
					}
					if lastedMergedAt.Before(mergedAt) {
						lastedMergedAt = mergedAt
					}
				}
				unpublishedServices := getUnpublishedServices(issue, services)
				if len(unpublishedServices) == 0 {
					continue
				}
				if !lastedMergedAt.IsZero() {
					rch <- &Issue{
						Key:                 issue.Key,
						MergedAt:            lastedMergedAt,
						UnpublishedServices: unpublishedServices,
					}
				}
			}
		}()
	}
	for _, issue := range issues {
		ch <- issue.Key
	}
	close(ch)
	wg.Wait()
	close(rch)
	for r := range rch {
		results = append(results, *r)
	}
	return
}

func getUnpublishedServices(issue *jira.Issue, filterServices []string) (results []string) {
	services, _ := issue.Fields.Unknowns[cmdutil.IssueFieldKeyServiceList].(string)
LOOP:
	for _, s := range strings.Split(services, "\n") {
		svc := strings.TrimSpace(s)
		if cmdutil.IsPublishedService(svc) {
			continue
		}
		for _, ps := range filterServices {
			if cmdutil.MatchServiceName(svc, ps) {
				results = append(results, ps)
				// try to next service
				continue LOOP
			}
		}
	}
	return
}
