package model

import (
	"fmt"
)

const (
	LoginReq       uint32 = 0
	FetchPlayerReq        = 11
	FetchRoomReq          = 12
	FetchGameReq          = 13
)

const (
	LoginRes       uint32 = 0
	FetchPlayerRes        = 14
	FetchRoomRes          = 15
	FetchGameRes          = 16
)

const (
	ErrorNone int32 = 0
)

type Player struct {
	UserId     uint32 `msgpack:"user_id"`
	Win        int32  `msgpack:"win"`
	Lose       int32  `msgpack:"lose"`
	Draw       int32  `msgpack:"draw"`
	Point      int32  `msgpack:"point"`
	Name       string `msgpack:"name"`
	State      int32  `msgpack:"state"`
	RoomId     int32  `msgpack:"room_id"`
	Permission int32  `msgpack:"permission"`
}

type Room struct {
	RoomId    int32  `msgpack:"room_id"`
	RoomTitle string `msgpack:"title"`
	RoomState int32  `msgpack:"state"`
	IsStart   bool   `msgpack:"is_start"`
	MinPoint  int32  `msgpack:"min_point"`
	MaxPoint  int32  `msgpack:"max_point"`
	HostId    uint32 `msgpack:"host_id"`
	OtherId   uint32 `msgpack:"other_id"`
}

type Game struct {
	RoomId    uint32   `msgpack:"room_id"`
	HostId    uint32   `msgpack:"host_id"`
	OtherId   uint32   `msgpack:"other_id"`
	HostName  string   `msgpack:"host_name"`
	OtherName string   `msgpack:"other_name"`
	WhoIsTurn uint32   `msgpack:"who_is_turn"`
	Title     string   `msgpack:"title"`
	Block     []uint32 `msgpack:"board"`
}

type LoginRequest struct {
	Email    string `msgpack:"email"`
	Password string `msgpack:"password"`
}

type LoginResponse struct {
	Error  int32  `msgpack:"error"`
	UserId uint32 `msgpack:"user_id"`
	Token  string `msgpack:"token"`
}

type FetchPlayerReuqest struct {
	Token  string `msgpack:"token"`
	UserId uint32 `msgpack:"user_id"`
}

type FetchPlayerResponse struct {
	Error   int32    `msgpack:"error"`
	Players []Player `msgpack:"players"`
}

type FetchRoomRequest struct {
	Token  string `msgpack:"token"`
	UserId uint32 `msgpack:"user_id"`
}

type FetchRoomResponse struct {
	Error int32  `msgpack:"error"`
	Rooms []Room `msgpack:"rooms"`
}

type FetchGameRequest struct {
	Token  string `msgpack:"token"`
	UserId uint32 `msgpack:"user_id"`
	RoomId uint32 `msgpack:"room_id"`
}

type FetchGameResponse struct {
	Error int32 `msgpack:"error"`
	Info  Game  `msgpack:"game"`
}

type Response struct {
	Protocol uint32
	Length   uint32
	Data     []byte
}

func (res LoginResponse) String() string {
	return fmt.Sprintf("Error : %d\nUserId : %d\nToken: %s\n", res.Error, res.UserId, res.Token)
}

func (p Player) String() string {

	var sp string
	var ss string

	if p.State == Playing {
		ss = "In Room"
	} else {
		ss = "In Lobby"
	}

	if p.Permission == User {
		sp = "User"
	} else {
		sp = "Admin"
	}

	return fmt.Sprintf("UserId: %d\tName: %s\tPermission: %s\tState: %s\n\nWin: %d\tLose: %d\tDraw: %d\tPoint: %d\t",
		p.UserId, p.Name, sp, ss, p.Win, p.Lose, p.Draw, p.Point)
}

func (r Room) String() string {

	var is string
	var rs string

	if r.RoomState == RoomReady {
		rs = "Ready"
	} else {
		rs = "Full"
	}

	if r.IsStart == true {
		is = "Playing"
	} else {
		is = "Not Playing"
	}

	return fmt.Sprintf("Room Id : %d\tTitle: %s\tState :%s\tIs Start: %s\n\nHost Id: %d\tOther Id: %d\tMin Point: %d\tMax Point: %d\n",
		r.RoomId, r.RoomTitle, rs, is, r.HostId, r.OtherId, r.MinPoint, r.MaxPoint)
}
