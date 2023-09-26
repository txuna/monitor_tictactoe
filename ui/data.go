package ui

import (
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

type NavElement struct {
	Title string
	Intro string
	View  func(w fyne.Window, conn net.Conn) fyne.CanvasObject
}

var NavElements = []NavElement{
	{
		"Logout", "Desc: Logout Layout", LogoutScreen,
	},
	{
		"Players", "Desc: Fetch All Players", PlayerScreen,
	},
	{
		"Rooms", "Desc: Fetch All Rooms", RoomScreen,
	},
	{
		"Game Info", "Desc: Fetch Game Info using room Id", GameScreen,
	},
}

var PlayerList = binding.BindStringList(
	&[]string{},
)

var RoomList = binding.BindStringList(
	&[]string{},
)

var GameInfo = binding.NewString()
