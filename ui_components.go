package main

import (
	"github.com/jtsunne/golastic/Views"
	"github.com/rivo/tview"
)

var (
	app           = tview.NewApplication()
	pages         = tview.NewPages()
	helpPage      = tview.NewTextView()
	tvNodes       = tview.NewTable()
	tvIndices     = tview.NewTable()
	tvDocsTable   = tview.NewTable()
	repoTable     = tview.NewTable()
	snapshotTable = tview.NewTable()
	header        = tview.NewTextView()
	footer        = tview.NewTextView()
	filter        = tview.NewInputField()
	tvInfo        = tview.NewTextView()
)

func initComponents() {
	header.SetBorder(true)
	header.SetDynamicColors(true)
	header.SetText("F1: Help | [yellow]F2: Nodes[white] | F3: Indices")

	footer.SetBorder(true).SetTitleAlign(tview.AlignRight).SetTitle(" Quick Help ")
	footer.SetText("i - Sort by Name | o - Sort by DocCount | Ctrl+E - Delete Index | Ctrl+P - Set Replicas")

	tvInfo.SetBorder(true)
	tvInfo.SetDynamicColors(true).SetRegions(true)
	// Event handlers like tvInfo.SetDoneFunc will be moved to event_handlers.go

	// Event handlers like tvDocsTable.SetDoneFunc will be moved to event_handlers.go

	// Input captures like tvIndices.SetInputCapture will be moved to event_handlers.go

	repoTable.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle("Repositories")
	snapshotTable.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle("Snapshots")
	// Input captures like tvNodes.SetInputCapture will be moved to event_handlers.go

	filter.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" Filter ")
	filter.SetLabel("Index Name Filter: ").SetFieldWidth(30)
	// Event handlers like filter.SetDoneFunc will be moved to event_handlers.go

	Views.MakeHelpPage(helpPage)

	pages.AddPage("info",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(tvInfo, 0, 1, true),
		true, true)
	pages.AddPage("repos",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(repoTable, 0, 1, true).
			AddItem(snapshotTable, 0, 1, false),
		true, true)
	pages.AddPage("help",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(helpPage, 0, 1, true),
		true, true)
	pages.AddPage("indices",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(filter, 3, 1, false).
			AddItem(tvIndices, 0, 1, true).
			AddItem(footer, 3, 1, false),
		true, true)
	pages.AddPage("docs",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(tvDocsTable, 0, 1, true).
			AddItem(footer, 3, 1, false),
		true, true)
	pages.AddPage("nodes",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(tvNodes, 0, 1, true).
			AddItem(footer, 3, 1, false),
		true, true)
}
