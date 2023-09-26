package parse

import (
	"encoding/binary"
	"tuuna/monitor/model"

	"github.com/vmihailenco/msgpack"
)

func CreateLoginRequestPacket(email, password string) []byte {
	login := &model.LoginRequest{
		Email:    email,
		Password: password,
	}

	m, err := msgpack.Marshal(login)
	if err != nil {
		panic(err)
	}

	msg := appendProtocolHeader(m, model.LoginReq)
	return msg
}

func CreateFetchGameRequestPacket(roomId uint32) []byte {
	packet := &model.FetchGameRequest{
		UserId: model.UserId,
		Token:  model.Token,
		RoomId: roomId,
	}

	m, err := msgpack.Marshal(packet)
	if err != nil {
		panic(err)
	}

	msg := appendProtocolHeader(m, model.FetchGameReq)
	return msg
}

func CreateFetchRoomRequestPacket() []byte {
	packet := &model.FetchRoomRequest{
		UserId: model.UserId,
		Token:  model.Token,
	}

	m, err := msgpack.Marshal(packet)
	if err != nil {
		panic(err)
	}

	msg := appendProtocolHeader(m, model.FetchRoomReq)
	return msg

}

func CreateFetchPlayerRequestPacket() []byte {
	packet := &model.FetchPlayerReuqest{
		UserId: model.UserId,
		Token:  model.Token,
	}

	m, err := msgpack.Marshal(packet)
	if err != nil {
		panic(err)
	}

	msg := appendProtocolHeader(m, model.FetchPlayerReq)
	return msg
}

func appendProtocolHeader(b []byte, p uint32) []byte {
	protocolByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(protocolByte, p)

	length := len(b)
	lengthByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthByte, uint32(length))

	msg := append(protocolByte, lengthByte...)
	msg = append(msg, b...)

	return msg
}

func ParseLoginResponse(data []byte, length uint32) *model.LoginResponse {
	var login model.LoginResponse

	err := msgpack.Unmarshal(data[0:length], &login)
	if err != nil {
		return nil
	}

	return &login
}

func ParseFetchPlayerResponse(data []byte, length uint32) *model.FetchPlayerResponse {
	var players model.FetchPlayerResponse

	err := msgpack.Unmarshal(data[0:length], &players)
	if err != nil {
		return nil
	}

	return &players
}

func ParseFetchRoomResponse(data []byte, length uint32) *model.FetchRoomResponse {
	var rooms model.FetchRoomResponse

	err := msgpack.Unmarshal(data[0:length], &rooms)
	if err != nil {
		return nil
	}

	return &rooms
}

func ParseFetchGameResponse(data []byte, length uint32) *model.FetchGameResponse {
	var game model.FetchGameResponse

	err := msgpack.Unmarshal(data[0:length], &game)
	if err != nil {
		return nil
	}

	return &game
}
