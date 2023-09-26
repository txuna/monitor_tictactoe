package model

var Token string = ""
var UserId uint32 = 0

const (
	Lobby   = 0
	Playing = 1

	RoomReady   = 0
	RoomPlaying = 1

	User  = 0
	Admin = 1

	NONE  = 0
	HOST  = 1
	OTHER = 2
)
