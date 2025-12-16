package models

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

type ParticipantsList struct {
	Participants []string `json:"participants"`
}
