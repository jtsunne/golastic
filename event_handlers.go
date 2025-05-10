package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jtsunne/golastic/Structs"
	"github.com/jtsunne/golastic/Utils" // Added Utils
	"github.com/rivo/tview"
)

// Event handler for tvInfo (when Escape is pressed)
func tvInfoDoneHandler(k tcell.Key) {
	if k == tcell.KeyEscape {
		header.SetText("F1: Help | F2: Nodes | [yellow]F3: Indices[white]")
		pages.SwitchToPage("indices")
	}
}

// Event handler for tvDocsTable (when Escape is pressed)
func tvDocsTableDoneHandler(k tcell.Key) {
	if k == tcell.KeyEscape {
		header.SetText("F1: Help | F2: Nodes | [yellow]F3: Indices[white]")
		pages.SwitchToPage("indices")
	}
}

// Event handler for tvDocsTable (when a cell is selected)
func tvDocsTableSelectedHandler(row, column int) {
	r, _ := tvDocsTable.GetSelection()
	name := tvDocsTable.GetCell(r, 1)
	tvInfo.SetText(Utils.PrettyJson(name.Text)) // Utils will need to be accessible or passed
	pages.SwitchToPage("info")
}

// Input capture for tvIndices
func tvIndicesInputCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyCtrlE:
		r, _ := tvIndices.GetSelection()
		name := tvIndices.GetCell(r, 0)
		DeleteIndexMessageBox(name.Text)
		RefreshData()
		return nil
	case tcell.KeyCtrlP:
		r, _ := tvIndices.GetSelection()
		name := tvIndices.GetCell(r, 0)
		SetReplicasMessageBox(name.Text)
		RefreshData()
		return nil
	case tcell.KeyCtrlBackslash:
		r, _ := tvIndices.GetSelection()
		name := tvIndices.GetCell(r, 0)
		tvInfo.SetTitle(fmt.Sprintf(" Index [%s] Documents ", name.Text))
		GetDocsFromIndex(name.Text) // Fetch data
		FillDocsTable(name.Text)    // Populate table
		pages.SwitchToPage("docs")
		return nil
	}
	switch event.Rune() {
	case 'I', 'i': // Combined cases for simplicity
		SortData("index")
		return nil
	case 'O', 'o': // Combined cases for simplicity
		SortData("docCount")
		return nil
	case '?':
		app.SetFocus(filter)
		return nil
	}
	return event
}

// Input capture for tvNodes
func tvNodesInputCaptureHandler(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'r':
		pages.SwitchToPage("repos")
		return nil
	}
	return event
}

// Event handler for filter (when Done is triggered, e.g. Enter pressed)
func filterDoneHandler(key tcell.Key) {
	FilterData(filter.GetText())
}

// Handler for when a cell is selected in tvIndices
func tvIndicesSelectedHandler(row, column int) {
	// Note: column is provided by tview, but original selectedIndexFunc ignored it.
	// We operate on the global tvIndices table.
	selectedIndexName := tvIndices.GetCell(row, 0).Text
	settingsJSON, err := GetIndexSettings(selectedIndexName)
	if err != nil {
		tvInfo.SetText(fmt.Sprintf("Error fetching settings: %v", err))
		pages.SwitchToPage("info")
		return
	}
	tvInfo.SetTitle(fmt.Sprintf(" Index [%s] Settings ", selectedIndexName))
	tvInfo.SetText(Utils.PrettyJson(settingsJSON))
	pages.SwitchToPage("info")
	tvInfo.SetDoneFunc(tvInfoDoneHandler) // Re-assign standard done handler for tvInfo
}

// Done handler for tvIndices
func tvIndicesDoneHandler(k tcell.Key) {
	switch k {
	case tcell.KeyEscape:
		tvIndices.SetSelectable(false, false)
	case tcell.KeyEnter:
		tvIndices.SetSelectable(true, false)
	}
}

// Done handler for tvNodes
func tvNodesDoneHandler(k tcell.Key) {
	switch k {
	case tcell.KeyEscape:
		tvNodes.SetSelectable(false, false)
	case tcell.KeyEnter:
		tvNodes.SetSelectable(true, false)
	}
}

// DeleteIndexMessageBox, now in event_handlers.go
func DeleteIndexMessageBox(idx string) {
	mb := tview.NewModal().
		SetText("Do you want to DELETE index " + idx + "?").
		AddButtons([]string{"Cancel", "DELETE"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				app.SetRoot(pages, true)
			}
			if buttonLabel == "DELETE" {
				err := DeleteIndex(idx)
				if err != nil {
					footer.SetText(fmt.Sprintf("Error deleting index %s: %v", idx, err))
				}
				RefreshData()
				app.SetRoot(pages, true)
			}
		})
	app.SetRoot(mb, true)
}

// SetReplicasMessageBox, now in event_handlers.go
func SetReplicasMessageBox(idxName string) {
	var idx Structs.EsIndices // This might cause issues if 'indices' global is not accessible
	var repl string
	// Need to ensure 'indices' global from es_data.go is accessible here or passed.
	// For now, assuming it's accessible for the logic to remain similar.
	for _, item := range indices {
		if item.Index == idxName {
			idx = item
		}
	}
	formRep := tview.NewForm().
		AddInputField("Replicas", idx.Rep, 5, nil, func(text string) {
			repl = text
		}).
		AddButton("Save", func() {
			err := SetReplicas(idxName, repl)
			if err != nil {
				footer.SetText(fmt.Sprintf("Error setting replicas for %s: %v", idxName, err))
			}
			RefreshData()
			app.SetRoot(pages, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(pages, true)
		})
	formRep.SetBorder(true).SetTitle("Set Replicas [" + idxName + "]")
	app.SetRoot(formRep, true)
}

func repoTableSelectedHandler(row, column int) {
    repoName := repoTable.GetCell(row, column).Text // Assuming column 0 for Id
    if column == 0 { // ensure it's the ID column if others are added
        FillSnapshot(repoName, snapshotTable)
        app.SetFocus(snapshotTable)
    }
}

func repoTableDoneHandler(key tcell.Key) {
    header.SetText("F1: Help | [yellow]F2: Nodes[white] | F3: Indices")
    pages.SwitchToPage("nodes")
}

func snapshotTableDoneHandler(key tcell.Key) {
    snapshotTable.Clear() // Clear snapshot table specific to this handler
    app.SetFocus(repoTable)
}


func initEventHandlers() {
	tvInfo.SetDoneFunc(tvInfoDoneHandler)
	tvDocsTable.SetDoneFunc(tvDocsTableDoneHandler)
	tvDocsTable.SetSelectedFunc(tvDocsTableSelectedHandler)
	tvIndices.SetInputCapture(tvIndicesInputCaptureHandler)
	tvIndices.SetSelectedFunc(tvIndicesSelectedHandler) // Updated handler
	tvIndices.SetDoneFunc(tvIndicesDoneHandler)         // Updated handler
	tvNodes.SetInputCapture(tvNodesInputCaptureHandler)
    tvNodes.SetDoneFunc(tvNodesDoneHandler) // Updated handler
	filter.SetDoneFunc(filterDoneHandler)

    // Specific handlers for repoTable and snapshotTable
    repoTable.SetSelectedFunc(repoTableSelectedHandler)
    repoTable.SetDoneFunc(repoTableDoneHandler)

    snapshotTable.SetDoneFunc(snapshotTableDoneHandler)
    // snapshotTable.SetSelectedFunc if needed for future actions
}
