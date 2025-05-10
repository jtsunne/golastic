package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/jtsunne/golastic/Utils"
	"github.com/rivo/tview"
)

var () // Keep other global variables if any, otherwise this can be removed if empty after all refactoring

func init() {
	if os.Getenv("ESURL") == "" {
		if len(os.Args) <= 1 {
			fmt.Println("Elasticsearch cluster have to be set. Exiting...")
			os.Exit(0)
		}
		EsUrl = Utils.ParseEsUrl(os.Args[1])
	} else {
		EsUrl = Utils.ParseEsUrl(os.Getenv("ESURL"))
	}
	fmt.Println("ES_URL=" + EsUrl)

	RefreshData()
	initComponents()    // Initialize UI components
	initEventHandlers() // Initialize event handlers
}

func main() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlH:
			popup := tview.NewBox().SetTitle("This is a popup window").SetBorder(true)
			app.SetRoot(popup, true)
		case tcell.KeyCtrlR:
			tvIndices.Clear()
			RefreshData()
			FilterData(filter.GetText())
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyF1:
			header.SetText("[yellow]F1: Help[white] | F2: Nodes | F3: Indices")
			pages.SwitchToPage("help")
		case tcell.KeyF2:
			footer.SetText("Ctrl+I - Sort by IP | Ctrl+O - Sort by Node | r - Show Repositories")
			header.SetText("F1: Help | [yellow]F2: Nodes[white] | F3: Indices")
			pages.SwitchToPage("nodes")
			return nil
		case tcell.KeyF3:
			footer.SetText("i - Sort by Name | o - Sort by DocCount")
			header.SetText("F1: Help | F2: Nodes | [yellow]F3: Indices[white]")
			pages.SwitchToPage("indices")
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
