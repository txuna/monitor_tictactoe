package main

import (
	"tuuna/monitor/controller"
	"tuuna/monitor/model"
	"tuuna/monitor/network"
	"tuuna/monitor/ui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("TicTacToe Server Monitor")
	w.Resize(fyne.NewSize(800, 600))

	conn := network.StartConnect()
	if conn == nil {
		return
	}

	ch := make(chan *model.Response)

	loginContainer := ui.MakeLoginUI(conn)
	w.SetContent(loginContainer)

	go network.RecvPacket(conn, ch)
	go controller.Controller(conn, w, ch)

	defer conn.Close()
	w.ShowAndRun()
}
