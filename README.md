# TicTacToe Server Monitoring GUI Program 
틱택토 서버를 모니터링하기 위해 go로 만든 GUI 프로그램입니다.   
GUI 라이브러리로는 fyne v2를 사용함.  

주 기능   
- 플레이어 목록 가져오기 
- 방 목록 가져오기 
- 게임아이디를 통해서 실시간 게임 정보 받아오기 

# Main Logic
```go
go network.RecvPacket(conn, ch)
go controller.Controller(conn, w, ch)

defer conn.Close()
w.ShowAndRun()
```
main 코드이며 패킷을 받는 고루틴과 패킷의 프로토콜에 따라 함수를 호출하기 위한 컨틀로러 고루틴, UI 업데이트를 위한 로직으로 구현되어 있습니다.

```go
func RecvPacket(conn net.Conn, ch chan *model.Response) {
	for {
		var offset uint32 = 0
		b := make([]byte, 4096)

		n, err := conn.Read(b)
		if err != nil {
			if err == io.EOF {
				log.Println("disconnected from server")
				return
			}
			log.Fatal("read Error: " + err.Error())
			return
		}

		if n <= 0 {
			return
		}

		if offset < uint32(n) {
			res := parseResponse(b, &offset)
			ch <- res
		}
	}
}
```
지속적으로 패킷을 받으며 받은 패킷에 파싱을 한 뒤 컨트롤러 채널에 데이터를 넘깁니다.

```go
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
```
컨트롤러 고루틴은 받은 패킷의 프로토콜을 확인하여 각 각 알맞은 함수에게 넘겨줍니다. 

```go
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
```
컨트롤러에서 호출하는 방목록 가져오기 함수입니다. ui.RoomList.Set을 통해 목록들을 추가하는 것을 볼 수 있습니다. 이는 UI쪽과 `DataBinding`으로 연결되어 있기떄문에 UI쪽에 즉각적으로 반영됩니다.

컨트롤러 코드에서 UI 부분을 건드렸기 떄문에 좋지 못한 코드라는 것을 인지하고 있습니다. 하지만 단순히 방목록을 String형태로 간편하게 보여주고 싶었기 때문에 해당 방식을 선택했습니다. 

만약 그렇지 않다면 UI관련 고루틴을 하나 만들어서 select으로 대기하고 channel을 통해서 처리한다면 조금 더 깔끔한 코드가 나올 수 있을거 같습니다.

```go

```


# Architecture 
![1](./images/architecture.jpg)

# Imges 
### Login
![2](./images/login.png)
### Fetch Players
![3](./images/players.png)
### Fetch Rooms
![4](./images/rooms.png)
### Fetch Game Info
![5](./images/game.png)