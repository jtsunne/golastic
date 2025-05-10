package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jtsunne/golastic/Structs"
	"github.com/rivo/tview"
)

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

	t.SetCell(0, 7, tview.NewTableCell("Shards").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 8, tview.NewTableCell("Disk.i").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 9, tview.NewTableCell("Disk %").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))

	t.SetCell(0, 10, tview.NewTableCell("LA 1m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 11, tview.NewTableCell("LA 5m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 12, tview.NewTableCell("LA 15m").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 13, tview.NewTableCell("Version").
		SetTextColor(tcell.ColorYellow).
		SetAlign(tview.AlignCenter))
	t.SetCell(0, 14, tview.NewTableCell("Tags").
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

		for _, itm := range nodesAllocation {
			if itm.Ip == item.IP {
				t.SetCellSimple(i+1, 7, itm.Shards)
				t.SetCellSimple(i+1, 8, itm.DiskIndices)
				t.SetCellSimple(i+1, 9, itm.DiskPercent)
			}
		}

		t.SetCellSimple(i+1, 10, item.Load1M)
		t.SetCellSimple(i+1, 11, item.Load5M)
		t.SetCellSimple(i+1, 12, item.Load15M)
		for _, itm := range clusterNodes {
			if itm.Name == item.Name {
				t.SetCellSimple(i+1, 13, itm.Version.Number)
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
		t.SetCellSimple(i+1, 14, s)
	}
	t.SetFixed(1, 1)
	// Event handlers like t.SetDoneFunc will be moved to event_handlers.go
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
	// Event handlers like t.SetSelectedFunc and t.SetDoneFunc will be moved to event_handlers.go
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
  // Event handlers like t.SetSelectedFunc and t.SetDoneFunc will be moved to event_handlers.go
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
	// Event handlers like t.SetSelectedFunc and t.SetDoneFunc will be moved to event_handlers.go
}

func FillDocsTable(idxName string) {
	tvDocsTable.Clear()
	tvDocsTable.SetTitle(fmt.Sprintf(" Index [%s] Documents ", idxName))
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
	// Event handlers like tvDocsTable.SetSelectedFunc will be moved to event_handlers.go
}
