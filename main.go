package main

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jtsunne/golastic/Structs"
	"github.com/jtsunne/golastic/Utils"
	"github.com/jtsunne/golastic/Views"
	"github.com/rivo/tview"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	EsUrl            string
	nodes            []Structs.EsNode
	indices          []Structs.EsIndices
	idxAliases       []Structs.EsIndexAlias
	clusterNodes     []Structs.EsClusterNode
	clusterNodesTags []Structs.EsClusterNodeTags
	clusterRepos     []Structs.EsClusterRepository
	clusterSnapshots map[string][]Structs.EsSnapshot
	docs             Structs.EsDocs
	app              = tview.NewApplication()
	pages            = tview.NewPages()
	helpPage         = tview.NewTextView()
	tvNodes          = tview.NewTable()
	tvIndices        = tview.NewTable()
	tvDocsTable      = tview.NewTable()
	repoTable        = tview.NewTable()
	snapshotTable    = tview.NewTable()
	header           = tview.NewTextView()
	footer           = tview.NewTextView()
	filter           = tview.NewInputField()
	tvInfo           = tview.NewTextView()
	sortIndexAsc     = true
	sortDocCountAsc  = true
	c                = &http.Client{Timeout: 10 * time.Second}
)

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

	header.SetBorder(true)
	header.SetText("F1: Help | F2: Nodes | F3: Indices")

	footer.SetBorder(true).SetTitleAlign(tview.AlignRight).SetTitle(" Quick Help ")
	footer.SetText("i - Sort by Name | o - Sort by DocCount | Ctrl+E - Delete Index | Ctrl+P - Set Replicas")

	tvInfo.SetBorder(true)
	tvInfo.SetDynamicColors(true).SetRegions(true)
	tvInfo.SetDoneFunc(func(k tcell.Key) {
		if k == tcell.KeyEscape {
			pages.SwitchToPage("indices")
		}
	})
	tvDocsTable.SetDoneFunc(func(k tcell.Key) {
		if k == tcell.KeyEscape {
			pages.SwitchToPage("indices")
		}
	})

	tvIndices.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			GetDocsFromIndex(name.Text)
			pages.SwitchToPage("docs")
			return nil
		}
		switch event.Rune() {
		case 'I':
			SortData("index")
			return nil
		case 'i':
			SortData("index")
			return nil
		case 'O':
			SortData("docCount")
			return nil
		case 'o':
			SortData("docCount")
			return nil
		case '?':
			app.SetFocus(filter)
			return nil
		}
		return event
	})

	repoTable.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle("Repositories")
	snapshotTable.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle("Snapshots")
	tvNodes.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		}
		switch event.Rune() {
		case 'r':
			pages.SwitchToPage("repos")
			return nil
		}
		return event
	})

	filter.SetBorder(true).
		SetTitleAlign(tview.AlignCenter).
		SetTitle(" Filter ")
	filter.SetLabel("Index Name Filter: ").SetFieldWidth(30)
	filter.SetDoneFunc(func(key tcell.Key) {
		FilterData(filter.GetText())
	})

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

func RefreshData() {
	Utils.GetJson(fmt.Sprintf("%s/_cat/nodes?format=json", EsUrl), &nodes)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	Utils.GetJson(fmt.Sprintf("%s/_cat/nodeattrs?format=json", EsUrl), &clusterNodesTags)

	var clusterNode Structs.EsClusterNode
	clusterNodes = []Structs.EsClusterNode{}
	for _, itm := range nodes {
		// TODO: do not hardcode http schema and port
		err := Utils.GetJson(fmt.Sprintf("http://%s:9200/", itm.IP), &clusterNode)
		if err != nil {
			clusterNode.Name = itm.Name
			clusterNode.Version.Number = "N/A"
		}
		clusterNodes = append(clusterNodes, clusterNode)

	}

	Utils.GetJson(fmt.Sprintf("%s/_cat/indices?format=json", EsUrl), &indices)
	sort.Slice(indices, func(i, j int) bool {
		return indices[i].Index < indices[j].Index
	})
	Utils.GetJson(fmt.Sprintf("%s/_cat/aliases?format=json", EsUrl), &idxAliases)
	Utils.GetJson(fmt.Sprintf("%s/_cat/repositories?format=json", EsUrl), &clusterRepos)
	clusterSnapshots = make(map[string][]Structs.EsSnapshot)
	for _, itm := range clusterRepos {
		var s []Structs.EsSnapshot
		Utils.GetJson(fmt.Sprintf("%s/_cat/snapshots/%s?format=json", EsUrl, itm.Id), &s)
		clusterSnapshots[itm.Id] = s
	}

	FillNodes(nodes, tvNodes)
	FillIndices(indices, tvIndices)
	FillRepos(clusterRepos, repoTable)
	dt := time.Now()
	footer.SetText("Data refreshed @ " + dt.Format(time.ANSIC))
}

func GetDocsFromIndex(idxName string) {
	b := []byte(`{"query": { "match_all": {} }, "size": 100}`)
	Utils.PostJson(fmt.Sprintf("%s/%s/_search", EsUrl, idxName), string(b), &docs)
	tvDocsTable.Clear()
	tvDocsTable.SetBorder(true)
	tvDocsTable.SetCell(0, 0, tview.NewTableCell("_id").
		SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	tvDocsTable.SetCell(0, 1, tview.NewTableCell("_source").
		SetTextColor(tcell.ColorYellow).SetAlign(tview.AlignCenter))
	for i, item := range docs.Hits.Hits {
		tvDocsTable.SetCellSimple(i+1, 0, item.Id)
		tvDocsTable.SetCellSimple(i+1, 1, string(item.Source))
	}
	tvDocsTable.SetFixed(2, 1)
	tvDocsTable.Select(2, 1)
	tvDocsTable.SetSelectable(true, false)
	tvDocsTable.SetSelectedFunc(func(row, column int) {
		r, _ := tvDocsTable.GetSelection()
		name := tvDocsTable.GetCell(r, 1)
		tvInfo.SetText(Utils.ColorizeJson(name.Text))
		pages.SwitchToPage("info")
	})
	tvInfo.SetDoneFunc(func(k tcell.Key) {
		if k == tcell.KeyEscape {
			pages.SwitchToPage("docs")
		}
	})

}

func SortData(sortBy string) {

	if sortBy == "docCount" {
		sort.Slice(indices, func(i, j int) bool {
			ii, _ := strconv.Atoi(indices[i].DocsCount)
			ij, _ := strconv.Atoi(indices[j].DocsCount)
			if sortDocCountAsc {
				return ii > ij
			}
			return ii < ij
		})
		sortDocCountAsc = !sortDocCountAsc
	}
	if sortBy == "index" {
		sort.Slice(indices, func(i, j int) bool {
			if sortIndexAsc {
				return indices[i].Index > indices[j].Index
			}
			return indices[i].Index < indices[j].Index
		})
		sortIndexAsc = !sortIndexAsc
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
	t.SetTitle(fmt.Sprintf(" Cluster Name: %s ", clusterNodes[0].ClusterName))
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
	t.SetCell(0, 10, tview.NewTableCell("Version").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 11, tview.NewTableCell("Tags").
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
		for _, itm := range clusterNodes {
			if itm.Name == item.Name {
				t.SetCellSimple(i+1, 10, itm.Version.Number)
			}
		}
		s := ""
		for _, itm := range clusterNodesTags {
			if itm.Ip == item.IP {
				if s == "" {
					s = itm.Attr + "=" + itm.Value
				} else {
					s = s + "," + itm.Attr + "=" + itm.Value
				}
			}
		}
		if s == "" {
			s = "N/A"
		}
		t.SetCellSimple(i+1, 11, s)
	}
	t.SetFixed(1, 1)
	//t.SetSelectedFunc(func(row, column int) {
	//	selectedFunc(row, column, t)
	//})
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
	t.SetCell(0, 1, tview.NewTableCell("Alias").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("Health").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Primary").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("Replicas").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("Doc Count").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("Doc Deleted").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Primary Store Size").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 9, tview.NewTableCell("Store Size").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	for i, item := range idxs {
		t.SetCellSimple(i+1, 0, item.Index)
		s := "-"
		for _, itm := range idxAliases {
			if itm.Index == item.Index {
				s = itm.Alias
			}
		}
		t.SetCellSimple(i+1, 1, s)
		t.SetCellSimple(i+1, 2, item.Health)
		t.SetCellSimple(i+1, 3, item.Status)
		t.SetCellSimple(i+1, 4, item.Pri)
		t.SetCellSimple(i+1, 5, item.Rep)
		t.SetCellSimple(i+1, 6, item.DocsCount)
		t.SetCellSimple(i+1, 7, item.DocsDeleted)
		t.SetCellSimple(i+1, 8, item.PriStoreSize)
		t.SetCellSimple(i+1, 9, item.StoreSize)
	}
	t.SetFixed(2, 1)
	t.Select(2, 1)
	t.SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		selectedIndexFunc(row, column, t)
	})
	t.SetDoneFunc(func(key tcell.Key) {
		tableDoneFunc(key, t)
	})
}

func FillRepos(r []Structs.EsClusterRepository, t *tview.Table) {
	t.Clear()
	t.SetBorder(true)
	t.SetCell(0, 0, tview.NewTableCell("Id").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("Type").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	for i, itm := range r {
		t.SetCellSimple(i+1, 0, itm.Id)
		t.SetCellSimple(i+1, 1, itm.Type)
	}
	t.SetFixed(2, 1)
	t.Select(2, 1)
	t.SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {
		repo := t.GetCell(row, 0).Text
		FillSnapshot(repo, snapshotTable)
		app.SetFocus(snapshotTable)
	})
	t.SetDoneFunc(func(key tcell.Key) {
		pages.SwitchToPage("nodes")
	})
}

func FillSnapshot(repoName string, t *tview.Table) {
	t.Clear()
	t.SetBorder(true)
	t.SetTitle("Snapshots [" + repoName + "]")
	t.SetCell(0, 0, tview.NewTableCell("Id").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 1, tview.NewTableCell("Status").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 2, tview.NewTableCell("Start Time").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 3, tview.NewTableCell("End Time").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 4, tview.NewTableCell("Duration").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 5, tview.NewTableCell("Indices").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 6, tview.NewTableCell("Successful Shards").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 7, tview.NewTableCell("Failed Shards").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Total Shards").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	for i, itm := range clusterSnapshots[repoName] {
		t.SetCellSimple(i+1, 0, itm.Id)
		t.SetCellSimple(i+1, 1, itm.Status)
		t.SetCellSimple(i+1, 2, itm.StartTime)
		t.SetCellSimple(i+1, 3, itm.EndTime)
		t.SetCellSimple(i+1, 4, itm.Duration)
		t.SetCellSimple(i+1, 5, itm.Indices)
		t.SetCellSimple(i+1, 6, itm.SuccessfulShards)
		t.SetCellSimple(i+1, 7, itm.FailedShards)
		t.SetCellSimple(i+1, 8, itm.TotalShards)
	}
	t.SetFixed(2, 1)
	t.Select(2, 1)
	t.SetSelectable(true, false)
	t.SetSelectedFunc(func(row, column int) {

	})
	t.SetDoneFunc(func(key tcell.Key) {
		t.Clear()
		app.SetFocus(repoTable)
	})
}

func selectedIndexFunc(row int, _ int, tbl *tview.Table) {
	var selectedIndexName string
	selectedIndexName = tbl.GetCell(row, 0).Text
	r, _ := c.Get(fmt.Sprintf("%s/%s/_settings?pretty", EsUrl, selectedIndexName))
	defer r.Body.Close()

	tvInfo.SetTitle(fmt.Sprintf(" Index [%s] Settings ", selectedIndexName))
	body, _ := io.ReadAll(r.Body)
	tvInfo.SetText(Utils.ColorizeJson(string(body)))
	pages.SwitchToPage("info")
	tvInfo.SetDoneFunc(func(k tcell.Key) {
		if k == tcell.KeyEscape {
			pages.SwitchToPage("indices")
		}
	})
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
			FilterData(filter.GetText())
		case tcell.KeyCtrlQ:
			app.Stop()
		case tcell.KeyF1:
			pages.SwitchToPage("help")
		case tcell.KeyF2:
			footer.SetText("Ctrl+I - Sort by IP | Ctrl+O - Sort by Node | r - Show Repositories")
			pages.SwitchToPage("nodes")
			return nil
		case tcell.KeyF3:
			footer.SetText("i - Sort by Name | o - Sort by DocCount")
			pages.SwitchToPage("indices")
			return nil
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
