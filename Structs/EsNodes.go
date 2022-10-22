package Structs

type EsNode struct {
	IP          string `json:"ip"`
	HeapPercent string `json:"heap.percent"`
	RAMPercent  string `json:"ram.percent"`
	CPU         string `json:"cpu"`
	Load1M      string `json:"load_1m"`
	Load5M      string `json:"load_5m"`
	Load15M     string `json:"load_15m"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
	Name        string `json:"name"`
}
