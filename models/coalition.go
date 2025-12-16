package models

type CoalitionsList struct {
	Coalitions []Coalition `json:"coalitions"`
}

type Coalition struct {
	CoalitionID int64  `json:"coalitionId"`
	Name        string `json:"name"`
}
