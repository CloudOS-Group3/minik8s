package prettyprint

import (
	"math/rand"
	"minik8s/util/log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func PrintTable(header []string, data [][]string) {
	rand.Seed(time.Now().UnixNano())
	colors := []tablewriter.Colors{
		tablewriter.Colors{tablewriter.FgBlueColor},
		tablewriter.Colors{tablewriter.FgGreenColor},
		tablewriter.Colors{tablewriter.FgCyanColor},
		tablewriter.Colors{tablewriter.FgYellowColor},
		tablewriter.Colors{tablewriter.FgMagentaColor},
		tablewriter.Colors{tablewriter.FgRedColor},
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	columnNum := 0
	headerColors := []tablewriter.Colors{}
	dataColors := []tablewriter.Colors{}
	for _, _ = range header {
		columnNum++
		randomInt := rand.Intn(len(colors))
		log.Info("%d", randomInt)
		randomColor := colors[randomInt]
		headerColors = append(headerColors, randomColor)
		dataColors = append(dataColors, randomColor)
	}

	table.SetHeaderColor(headerColors...)

	table.SetColumnColor(dataColors...)

	table.AppendBulk(data)
	table.Render()
}
