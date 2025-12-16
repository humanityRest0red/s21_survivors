package models

type ParticipantsWorkstation struct {
	ClusterId   int    `json:"clusterId"`
	ClusterName string `json:"clusterName"`
	Row         string `json:"row"`
	Number      int    `json:"number"`
}
