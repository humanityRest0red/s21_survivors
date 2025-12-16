package models

type MYJSON interface {
	Participant | ParticipantsList |
		Campus | CampusesList |
		Coalition | CoalitionsList |
		ParticipantsWorkstation
}
