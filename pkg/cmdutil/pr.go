package cmdutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/github"
)

func IsPRMerged(cli *github.Client, ctx context.Context, owner, repo string, number int) (merged bool, mergedAt time.Time, err error) {
	pr, _, err := cli.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		fmt.Fprintf(os.Stderr, "get pr %s/%s#%d error: %v\n", owner, repo, number, err)
		return
	}
	merged = pr.GetMerged()
	mergedAt = pr.GetMergedAt()
	return
}
