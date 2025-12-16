package models

type CampusesList struct {
	Campuses []Campus `json:"campuses"`
}

type Campus struct {
	ID        string `json:"id"`
	ShortName string `json:"shortName"`
	FullName  string `json:"fullName"`
}
