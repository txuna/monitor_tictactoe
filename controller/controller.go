package controller

import (
	"fmt"
	"net"
	"tuuna/monitor/model"
	"tuuna/monitor/parse"
	"tuuna/monitor/ui"

	"fyne.io/fyne/v2"
)

func Controller(conn net.Conn, w fyne.Window, ch chan *model.Response) {
	for {
		select {
		case res := <-ch:
			switch res.Protocol {
			case model.LoginRes:
				LoginController(conn, w, res)

			case model.FetchPlayerRes:
				FetchPlayerController(conn, w, res)

			case model.FetchRoomRes:
				FetchRoomController(conn, w, res)

			case model.FetchGameRes:
				FetchGameController(conn, w, res)
			}
		}
	}
}

func LoginController(conn net.Conn, w fyne.Window, res *model.Response) {
	login := parse.ParseLoginResponse(res.Data, res.Length)
	if login == nil || login.Error != model.ErrorNone {
		ui.AlertUI("Login Failed!")
		return
	}

	model.Token = login.Token
	model.UserId = login.UserId

	content := ui.MakeMainUI(conn, w)
	w.SetContent(content)
}

/*
w의 content를 덮어쓰면 안되고 현재 보고 있는 컨텐츠가 players 컨텐츠인지 확인하고
특정한 곳에 써야됨
*/
func FetchPlayerController(conn net.Conn, w fyne.Window, res *model.Response) {
	players := parse.ParseFetchPlayerResponse(res.Data, res.Length)
	if players == nil || players.Error != model.ErrorNone {
		ui.AlertUI("Failed Fetch Players")
		return
	}

	/* 초기화 */
	ui.PlayerList.Set([]string{})

	/* 추가 */
	for i := 0; i < len(players.Players); i++ {
		s := players.Players[i].String()
		ui.PlayerList.Append(s)
	}

}

func FetchRoomController(conn net.Conn, w fyne.Window, res *model.Response) {
	rooms := parse.ParseFetchRoomResponse(res.Data, res.Length)
	if rooms == nil || rooms.Error != model.ErrorNone {
		ui.AlertUI("Failed Fetch Rooms")
		return
	}

	ui.RoomList.Set([]string{})

	for i := 0; i < len(rooms.Rooms); i++ {
		s := rooms.Rooms[i].String()
		ui.RoomList.Append(s)
	}
}

func FetchGameController(conn net.Conn, w fyne.Window, res *model.Response) {
	game := parse.ParseFetchGameResponse(res.Data, res.Length)
	if game == nil || game.Error != model.ErrorNone {
		ui.AlertUI("Failed Fetch Rooms")
		return
	}

	str := fmt.Sprintf("Room Id: %d\tRoom Title: %s\tHost Name: %s(%d)\tOther Name: %s(%d)\nWho Is Turn: %d\n\n",
		game.Info.RoomId, game.Info.Title, game.Info.HostName, game.Info.HostId, game.Info.OtherName, game.Info.OtherId, game.Info.WhoIsTurn)

	for i := 0; i < 9; i++ {
		if game.Info.Block[i] == model.NONE {
			str += "*\t"
		} else if game.Info.Block[i] == model.HOST {
			str += "O\t"
		} else if game.Info.Block[i] == model.OTHER {
			str += "X\t"
		}
		if i%3 == 2 {
			str += "\n"
		}
	}

	ui.GameInfo.Set(str)
}
