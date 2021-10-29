package cmd

import (
	"context"
	"strings"

	"github.com/k1LoW/octocov/config"
	"github.com/k1LoW/octocov/gh"
	"github.com/k1LoW/octocov/report"
)

func commentReport(ctx context.Context, c *config.Config, r, rOrig *report.Report) error {
	repo, err := gh.Parse(c.Repository)
	if err != nil {
		return err
	}
	g, err := gh.New()
	if err != nil {
		return err
	}
	n, err := g.DetectCurrentPullRequestNumber(ctx, repo.Owner, repo.Repo)
	if err != nil {
		return err
	}
	files, err := g.GetPullRequestFiles(ctx, repo.Owner, repo.Repo, n)
	if err != nil {
		return err
	}
	footer := "Reported by [octocov](https://github.com/k1LoW/octocov)"
	if c.Comment.HideFooterLink {
		footer = "Reported by octocov"
	}
	var table, fileTable string
	if rOrig != nil {
		d := rOrig.Compare(r)
		table = d.Table()
		fileTable = d.FileCoveagesTable(files)
	} else {
		table = r.Table()
		fileTable = r.FileCoveagesTable(files)
	}

	comment := strings.Join([]string{
		"## Code Metrics Report",
		table,
		"",
		fileTable,
		"---",
		footer,
	}, "\n")
	if err := g.PutComment(ctx, repo.Owner, repo.Repo, n, comment); err != nil {
		return err
	}
	return nil
}
