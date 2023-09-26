package ui

import (
	"image/color"
	"net"
	"strconv"
	"tuuna/monitor/network"
	"tuuna/monitor/parse"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type numericalEntry struct {
	widget.Entry
}

func newNumericalEntry() *numericalEntry {
	entry := &numericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *numericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *numericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func AlertUI(msg string) {
	w := fyne.CurrentApp().NewWindow("Alert")
	w.SetContent(container.NewCenter(widget.NewLabel(msg)))
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(200, 100))
	w.Show()
}

func MakeLoginUI(conn net.Conn) fyne.CanvasObject {

	title := canvas.NewText("TicTacToe Monitoring System", color.White)
	title.TextSize = 32.0

	emailField := widget.NewEntry()
	emailField.SetPlaceHolder("Email...")

	passwordField := widget.NewEntry()
	passwordField.Password = true
	passwordField.SetPlaceHolder("Password...")

	loginButton := widget.NewButton("Login", func() {
		msg := parse.CreateLoginRequestPacket(emailField.Text, passwordField.Text)
		go network.SendPacket(conn, msg)
	})

	boxContainer := container.NewVBox(title,
		widget.NewLabel(""),
		emailField,
		passwordField,
		loginButton)

	centerContainer := container.NewCenter(boxContainer)

	return centerContainer
}

func MakeMainUI(conn net.Conn, w fyne.Window) fyne.CanvasObject {
	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("Navbar")
	intro.Wrapping = fyne.TextWrapWord

	setNav := func(e NavElement, conn net.Conn) {
		title.SetText(e.Title)
		intro.SetText(e.Intro)

		/* 해당 레이아웃 컨테이너 반환 */
		/* 세팅된 content에 값을 넣어 동적 변화 */
		content.Objects = []fyne.CanvasObject{e.View(w, conn)}
	}
	/* content 세팅 */
	element := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)

	split := container.NewHSplit(MakeNavUI(setNav, conn), element)
	split.Offset = 0.2
	return split
}

func MakeNavUI(SetNav func(element NavElement, conn net.Conn), conn net.Conn) fyne.CanvasObject {
	a := fyne.CurrentApp()
	list := widget.NewList(
		func() int {
			return len(NavElements)
		},

		func() fyne.CanvasObject {
			return widget.NewButton("HEllo World", nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			element := NavElements[i]
			o.(*widget.Button).SetText(element.Title)
			o.(*widget.Button).OnTapped = func() {
				SetNav(element, conn)
			}
		})

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, list)
}

/*
어떻게 네트워크를 호출할 것인가
PlayerList를 Set으로 새롭게 정의해야할듯 컨트롤러에서 불러와야함
*/
func PlayerScreen(_ fyne.Window, conn net.Conn) fyne.CanvasObject {

	fetchButton := widget.NewButton("Fetch Player List", func() {
		// 패킷을 보내면 컨트롤러에서는 PlayerList.Set을 통해서 값 넣음
		packet := parse.CreateFetchPlayerRequestPacket()
		go network.SendPacket(conn, packet)
		/*
			PlayerList.Set([]string{
				Players[0].String(),
			})
		*/
	})

	playerList := widget.NewListWithData(
		PlayerList,
		func() fyne.CanvasObject {
			return widget.NewLabel("\ntemplate\n\n")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	content := container.NewBorder(nil, fetchButton, nil, nil, playerList)
	return content
}

func LogoutScreen(_ fyne.Window, conn net.Conn) fyne.CanvasObject {
	return container.NewCenter(widget.NewLabel("Hello Logout"))
}

func RoomScreen(_ fyne.Window, conn net.Conn) fyne.CanvasObject {
	fetchButton := widget.NewButton("Fetch Room List", func() {
		packet := parse.CreateFetchRoomRequestPacket()
		go network.SendPacket(conn, packet)
	})

	roomList := widget.NewListWithData(
		RoomList,
		func() fyne.CanvasObject {
			return widget.NewLabel("\ntemplate\n\n")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)
	content := container.NewBorder(nil, fetchButton, nil, nil, roomList)
	return content
}

func GameScreen(_ fyne.Window, conn net.Conn) fyne.CanvasObject {
	roomIdField := newNumericalEntry()
	roomIdField.SetPlaceHolder("Write RoomId...")
	fetchButton := widget.NewButton("Fetch Game Infom", func() {
		if roomIdField.Text == "" {
			return
		}
		u32, err := strconv.ParseUint(roomIdField.Text, 10, 32)
		if err != nil {
			return
		}
		packet := parse.CreateFetchGameRequestPacket(uint32(u32))
		go network.SendPacket(conn, packet)
	})

	GameInfo.Set("Game Info Panel!")

	label := widget.NewLabelWithData(GameInfo)
	center := container.NewCenter(label)
	content := container.NewBorder(roomIdField, fetchButton, center, nil, nil)
	return content
}
