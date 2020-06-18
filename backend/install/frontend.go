package install

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
)

func (i *Install) startProgress() {

	var percent int

	go func() {
		for percent < 97 {
			percent++
			time.Sleep(2000 * time.Millisecond)
			i.incrementProgress(percent)
		}
	}()

	for {

		select {
		case percent = <-i.incrementProgressCh:
			i.incrementProgress(percent)
		case msg := <-i.progressMessageCh:
			i.sendStatusMsg(msg)
		}
	}
}

func (i *Install) updateProgress(progress int, progressMsg string) {
	i.incrementProgressCh <- progress
	i.progressMessageCh <- progressMsg
	log.Infoln(progressMsg)
}

func (i *Install) sendStatusMsg(msg string) {
	i.frontend.Events.Emit("status", msg)
	return
}

func (i *Install) incrementProgress(percent int) {
	i.frontend.Events.Emit("progress", percent)
	return
}

func (i *Install) sendErrorNotification(title, msg string) {
	i.frontend.Events.Emit("error", title, msg)
}

func (i *Install) sendSuccessNotification(title, msg string) {
	i.frontend.Events.Emit("success", title, msg)
}

// WailsInit will be called automatically when the binary runs.
func (i *Install) WailsInit(runtime *wails.Runtime) error {
	i.frontend = runtime
	return nil
}
