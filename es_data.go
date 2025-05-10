package main

import (
	"bytes"
	"fmt"
	"github.com/jtsunne/golastic/Structs"
	"github.com/jtsunne/golastic/Utils"
	"net/http"
	"sort"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	EsUrl            string
	nodes            []Structs.EsNode
	nodesAllocation  []Structs.EsNodeAllocation
	indices          []Structs.EsIndices
	idxAliases       []Structs.EsIndexAlias
	clusterNodes     []Structs.EsClusterNode
	clusterNodesTags []Structs.EsClusterNodeTags
	clusterRepos     []Structs.EsClusterRepository
	clusterSnapshots map[string][]Structs.EsSnapshot
	docs             Structs.EsDocs
	c                = &http.Client{Timeout: 10 * time.Second}
	sortIndexAsc     = true
	sortDocCountAsc  = true
)

func RefreshData() {
	Utils.GetJson(fmt.Sprintf("%s/_cat/nodes?format=json", EsUrl), &nodes)
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})

	Utils.GetJson(fmt.Sprintf("%s/_cat/allocation?format=json", EsUrl), &nodesAllocation)

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
	// UI filling part will be moved to ui_fillers.go
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

func DeleteIndex(idxName string) error {
	req, _ := http.NewRequest("DELETE", EsUrl+"/"+idxName, nil)
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func SetReplicas(idxName string, replicaCount string) error {
	jsonBody := []byte(`{"index" : { "number_of_replicas":` + replicaCount + ` }}`)
	bodyReader := bytes.NewReader(jsonBody)
	reqUrl := fmt.Sprintf(EsUrl+"/%s/_settings", idxName)
	req, _ := http.NewRequest(http.MethodPut, reqUrl, bodyReader)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r, err := c.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func GetIndexSettings(idxName string) (string, error) {
    r, err := c.Get(fmt.Sprintf("%s/%s/_settings?pretty", EsUrl, idxName))
    if err != nil {
        return "", err
    }
    defer r.Body.Close()

    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        return "", err
    }
    return string(bodyBytes), nil
}
