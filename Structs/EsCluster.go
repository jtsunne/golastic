package Structs

import "time"

type EsClusterNode struct {
	Name        string `json:"name"`
	ClusterName string `json:"cluster_name"`
	ClusterUuid string `json:"cluster_uuid"`
	Version     struct {
		Number                           string    `json:"number"`
		BuildFlavor                      string    `json:"build_flavor"`
		BuildType                        string    `json:"build_type"`
		BuildHash                        string    `json:"build_hash"`
		BuildDate                        time.Time `json:"build_date"`
		BuildSnapshot                    bool      `json:"build_snapshot"`
		LuceneVersion                    string    `json:"lucene_version"`
		MinimumWireCompatibilityVersion  string    `json:"minimum_wire_compatibility_version"`
		MinimumIndexCompatibilityVersion string    `json:"minimum_index_compatibility_version"`
	} `json:"version"`
	Tagline string `json:"tagline"`
}

type EsClusterNodeTags struct {
	Node  string `json:"node"`
	Host  string `json:"host"`
	Ip    string `json:"ip"`
	Attr  string `json:"attr"`
	Value string `json:"value"`
}

type EsClusterRepository struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type EsSnapshot struct {
	Id               string `json:"id"`
	Status           string `json:"status"`
	StartEpoch       string `json:"start_epoch"`
	StartTime        string `json:"start_time"`
	EndEpoch         string `json:"end_epoch"`
	EndTime          string `json:"end_time"`
	Duration         string `json:"duration"`
	Indices          string `json:"indices"`
	SuccessfulShards string `json:"successful_shards"`
	FailedShards     string `json:"failed_shards"`
	TotalShards      string `json:"total_shards"`
}
