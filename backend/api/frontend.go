package api

import (
	"github.com/wailsapp/wails"
)

// Frontend carries the FE connector
type Frontend struct {
	RT *wails.Runtime
}

// WailsInit will be called automatically when the binary runs.
func (fe *Frontend) WailsInit(RT *wails.Runtime) {
	fe.RT = RT
}

func (fe *Frontend) sendStatusMsg(msg string) {
	fe.RT.Events.Emit("progress", msg)
	return
}

func (fe *Frontend) sendProgress(percent int) {
	fe.RT.Events.Emit("progress", percent)
	return
}
