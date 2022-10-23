package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jtsunne/golastic/Structs"
	"github.com/rivo/tview"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	EsUrl     string
	nodes     []Structs.EsNode
	indices   []Structs.EsIndices
	app       = tview.NewApplication()
	pages     = tview.NewPages()
	tvNodes   = tview.NewTable()
	tvIndices = tview.NewTable()
	header    = tview.NewTextView()
	footer    = tview.NewTextView()
	filter    = tview.NewInputField()
	formRep   = tview.NewForm()
	c         = &http.Client{Timeout: 10 * time.Second}
)

func getJson(url string, target interface{}) error {
	r, err := c.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func init() {
	if os.Getenv("ESURL") == "" {
		if len(os.Args) <= 1 {
			fmt.Println("Elasticsearch cluster have to be set. Exiting...")
			os.Exit(0)
		}
		EsUrl = os.Args[1]
	} else {
		EsUrl = os.Getenv("ESURL")
	}
	fmt.Println("ES_URL=" + EsUrl)

	RefreshData()

	header.SetBorder(true)
	header.SetText("F1: nodes | F2: indices")

	footer.SetBorder(true).SetTitleAlign(tview.AlignRight).SetTitle(" Help ")
	footer.SetText("Ctrl+I - Sort by Name | Ctrl+O - Sort by DocCount")

	filter.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" Filter ")
	filter.SetLabel("Index Name Filter: ").SetFieldWidth(30)
	filter.SetDoneFunc(func(key tcell.Key) {
		FilterData(filter.GetText())
	})

	pages.AddPage("nodes",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(tvNodes, 0, 1, true).
			AddItem(footer, 3, 1, false),
		true, true)
	pages.AddPage("indices",
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 3, 1, false).
			AddItem(filter, 3, 1, false).
			AddItem(tvIndices, 0, 1, true).
			AddItem(footer, 3, 1, false),
		true, true)
}

func RefreshData() {
	getJson(fmt.Sprintf("%s/_cat/nodes?format=json", EsUrl), &nodes)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	getJson(fmt.Sprintf("%s/_cat/indices?format=json", EsUrl), &indices)
	sort.Slice(indices, func(i, j int) bool {
		return indices[i].Index < indices[j].Index
	})

	FillNodes(nodes, tvNodes)
	FillIndices(indices, tvIndices)
}

func SortData(sortBy string) {
	if sortBy == "docCount" {
		sort.Slice(indices, func(i, j int) bool {
			ii, _ := strconv.Atoi(indices[i].DocsCount)
			ij, _ := strconv.Atoi(indices[j].DocsCount)
			return ii < ij
		})
	}
	if sortBy == "index" {
		sort.Slice(indices, func(i, j int) bool {
			return indices[i].Index < indices[j].Index
		})
	}

	FillIndices(indices, tvIndices)
}

func FilterData(s string) {
	var idxs []Structs.EsIndices
	for _, item := range indices {
		if strings.Contains(item.Index, s) {
			idxs = append(idxs, item)
		}
	}
	tvIndices.Clear()
	FillIndices(idxs, tvIndices)
	app.SetFocus(tvIndices)
}

func FillNodes(n []Structs.EsNode, t *tview.Table) {
	t.Clear()
	t.SetBorder(true)
	t.SetCell(0, 0, tview.NewTableCell("IP").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("Name").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("Master").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("Role").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Heap %").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("RAM %").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("CPU").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("LA 1m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("LA 5m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 9, tview.NewTableCell("LA 15m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	for i, item := range n {
		t.SetCellSimple(i+1, 0, item.IP)
		t.SetCellSimple(i+1, 1, item.Name)
		t.SetCellSimple(i+1, 2, item.Master)
		t.SetCellSimple(i+1, 3, item.NodeRole)
		t.SetCellSimple(i+1, 4, item.HeapPercent)
		t.SetCellSimple(i+1, 5, item.RAMPercent)
		t.SetCellSimple(i+1, 6, item.CPU)
		t.SetCellSimple(i+1, 7, item.Load1M)
		t.SetCellSimple(i+1, 8, item.Load5M)
		t.SetCellSimple(i+1, 9, item.Load15M)
	}
	t.SetFixed(1, 1)
	t.SetSelectedFunc(func(row, column int) {
		selectedFunc(row, column, t)
	})
	t.SetDoneFunc(func(key tcell.Key) {
		tableDoneFunc(key, t)
	})
}

func FillIndices(idxs []Structs.EsIndices, t *tview.Table) {
	t.Clear()
	t.SetBorder(true)
	t.SetCell(0, 0, tview.NewTableCell("Index").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("Health").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("Primary").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Replicas").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("Doc Count").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("Doc Deleted").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("Primary Store Size").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Store Size").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	for i, item := range idxs {
		t.SetCellSimple(i+1, 0, item.Index)
		t.SetCellSimple(i+1, 1, item.Health)
		t.SetCellSimple(i+1, 2, item.Status)
		t.SetCellSimple(i+1, 3, item.Pri)
		t.SetCellSimple(i+1, 4, item.Rep)
		t.SetCellSimple(i+1, 5, item.DocsCount)
		t.SetCellSimple(i+1, 6, item.DocsDeleted)
		t.SetCellSimple(i+1, 7, item.PriStoreSize)
		t.SetCellSimple(i+1, 8, item.StoreSize)
	}
	t.SetFixed(2, 1)
	t.Select(2, 1)
	t.SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		selectedFunc(row, column, t)
	})
	t.SetDoneFunc(func(key tcell.Key) {
		tableDoneFunc(key, t)
	})
}

func selectedFunc(r int, c int, tbl *tview.Table) {
	//for i := 0; i < tbl.GetColumnCount(); i++ {
	//	tbl.GetCell(r, i).SetTextColor(tcell.ColorRed)
	//}
}

func tableDoneFunc(k tcell.Key, tbl *tview.Table) {
	switch k {
	case tcell.KeyEscape:
		tbl.SetSelectable(false, false)
	case tcell.KeyEnter:
		tbl.SetSelectable(true, false)
	}
}

func DeleteIndexMessageBox(idx string) {
	mb := tview.NewModal().
		SetText("Do you want to DELETE index " + idx + "?").
		AddButtons([]string{"Cancel", "DELETE"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Cancel" {
				app.SetRoot(pages, true)
			}
			if buttonLabel == "DELETE" {
				req, _ := http.NewRequest("DELETE", EsUrl+"/"+idx, nil)
				r, _ := c.Do(req)
				defer r.Body.Close()
				app.SetRoot(pages, true)
			}
		})
	app.SetRoot(mb, true)
}

func SetReplicasMessageBox(idxName string) {
	var idx Structs.EsIndices
	var repl string
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
			jsonBody := []byte(`{"index" : { "number_of_replicas":` + repl + ` }}`)
			bodyReader := bytes.NewReader(jsonBody)
			reqUrl := fmt.Sprintf(EsUrl+"/%s/_settings", idxName)
			req, _ := http.NewRequest(http.MethodPut, reqUrl, bodyReader)
			req.Header.Set("Content-Type", "application/json; charset=UTF-8")
			r, _ := c.Do(req)
			defer r.Body.Close()
			app.SetRoot(pages, true)
		}).
		AddButton("Cancel", func() {
			app.SetRoot(pages, true)
		})
	formRep.SetBorder(true).SetTitle("Set Replicas [" + idxName + "]")
	app.SetRoot(formRep, true)
}

func main() {
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlR:
			tvIndices.Clear()
			RefreshData()
		case tcell.KeyCtrlI:
			SortData("index")
		case tcell.KeyCtrlO:
			SortData("docCount")
		case tcell.KeyCtrlBackslash:
			footer.SetText("KeyBS pressed")
			app.SetFocus(filter)
		case tcell.KeyCtrlE:
			r, _ := tvIndices.GetSelection()
			name := tvIndices.GetCell(r, 0)
			DeleteIndexMessageBox(name.Text)
			RefreshData()
		case tcell.KeyCtrlP:
			r, _ := tvIndices.GetSelection()
			name := tvIndices.GetCell(r, 0)
			SetReplicasMessageBox(name.Text)
			RefreshData()
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyF1:
			footer.SetText("Ctrl+I - Sort by IP | Ctrl+O - Sort by Node")
			pages.SwitchToPage("nodes")
			return nil
		case tcell.KeyF2:
			footer.SetText("Ctrl+I - Sort by Name | Ctrl+O - Sort by DocCount")
			pages.SwitchToPage("indices")
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
