package install

import (
	"fmt"
	"path"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// Uninstall will Uninstall the application by removing the files located
// in the .dag directory:
//   'Molly Wallet.lnk'   cl-wallet.jar   mollywallet.exe   tmp
//   cl-keytool.jar      install.log     store.db          wallet.log
// And also removing the shortcuts on Windows.
func (i *Install) Uninstall() {

	files := make([]string, 12)
	files = append(files, "update.log", "wallet.log", "install.log", "store.db", "cl-keytool.jar.tmp", "cl-keytool.jar", "cl-wallet.jar", "cl-wallet.jar.tmp", "mollywallet.zip", "mollywallet.zip.tmp", "Molly Wallet.lnk", "mollywallet.exe")

	err := removeFiles(i.dagFolderPath, files)
	if err != nil {
		i.sendErrorNotification("Unable to remove all files", fmt.Sprintf("%v", err))
	}

	folders := make([]string, 3)
	folders = append(folders, path.Join(i.dagFolderPath, "tmp"), i.tmpFolderPath, i.dagFolderPath)

	err = removeFolders(folders)
	if err != nil {
		i.sendErrorNotification("Unable to remove all folders", fmt.Sprintf("%v", err))
	}

	if runtime.GOOS == "windows" {
		err := removeFile(i.OSSpecificSettings.startMenuPath, "Molly Wallet.lnk")
		if err != nil {
			i.sendErrorNotification("Unable to remove shortcut from start menu", fmt.Sprintf("%v", err))
			log.Errorf("unable to remove shortcut from start menu: %v", err)
		}
		err = removeFile(i.OSSpecificSettings.desktopPath, "Molly Wallet.lnk")
		if err != nil {
			i.sendErrorNotification("Unable to remove shortcut from desktop", fmt.Sprintf("%v", err))
			log.Errorf("unable to remove shortcut from desktop: %v", err)
		}
	}

	i.sendSuccessNotification("Success!", "Molly wallet has been successfully uninstalled.")

	// i.frontend.Window.Close()

}
