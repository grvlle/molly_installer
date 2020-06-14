package install

import "github.com/wailsapp/wails"

// WailsInit will be called automatically when the binary runs.
func (i *Install) WailsInit(runtime *wails.Runtime) error {
	i.frontend = runtime
	return nil
}

func (i *Install) sendStatusMsg(msg string) {
	i.frontend.Events.Emit("status", msg)
	return
}

func (i *Install) sendProgress(percent string) {
	i.frontend.Events.Emit("progress", percent)
	return
}
