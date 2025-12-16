package models

type MYJSON interface {
	Participant | ParticipantsWorkstation | ParticipantsResponse | Campus
}

type Participant struct {
	Login          string `json:"login"`
	ClassName      string `json:"className"`
	ParallelName   string `json:"parallelName"`
	ExpValue       int64  `json:"expValue"`
	Level          int32  `json:"level"`
	ExpToNextLevel int64  `json:"expToNextLevel"`
	Campus         Campus `json:"campus"`
	Status         string `json:"status"`
}

type ParticipantsWorkstation struct {
	ClusterId   int    `json:"clusterId"`
	ClusterName string `json:"clusterName"`
	Row         string `json:"row"`
	Number      int    `json:"number"`
}

type ParticipantsResponse struct {
	Participants []string `json:"participants"`
}

type Campus struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
}

type Coalition struct {
	CoalitionID string `json:"coalitionId"`
	Name        string `json:"name"`
}
