package cmdutil

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func WriteTable(data [][]string, cols ...string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(cols)
	table.SetBorder(false)
	table.AppendBulk(data)

	table.Render()
}
