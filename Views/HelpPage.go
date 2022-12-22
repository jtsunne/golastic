package Views

import "github.com/rivo/tview"

func MakeHelpPage(v *tview.TextView) {
	v.SetBorder(true).SetTitleAlign(tview.AlignCenter).SetTitle(" Help Page ")
	v.SetText(`
Use F1 - to see this Help Page
    F2 - to see the Nodes Page
    F3 - to see the Indices Page

Global HotKeys:
    Ctrl+R - refresh all data
	Ctrl+Q - Quit

Indices Page HotKeys:
	i      - Sort indices by Index Name
	o      - Sort indices by Documents Count
	?      - Set filter for the Indices view
	Ctrl+E - Remove selected Index
	Ctrl+P - Set Replicas Count for the selected index
	Enter  - Show Index Settings
`)

}
