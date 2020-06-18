package install

import (
	"os"
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
		log.Errorf("unable to remove all files: %v", err)
	}

	folders := make([]string, 3)
	folders = append(folders, path.Join(i.dagFolderPath, "tmp"), i.tmpFolderPath, i.dagFolderPath)

	err = removeFolders(folders)
	if err != nil {
		log.Errorf("unable to remove all folders: %v", err)
	}

	if runtime.GOOS == "windows" {
		err := removeFile(i.OSSpecificSettings.startMenuPath, "Molly Wallet.lnk")
		if err != nil {
			log.Errorf("unable to remove shortcut from start menu: %v", err)
		}
		err = removeFile(i.OSSpecificSettings.desktopPath, "Molly Wallet.lnk")
		if err != nil {
			log.Errorf("unable to remove shortcut from desktop: %v", err)
		}
	}

	i.sendSuccessNotification("Success!", "Molly wallet has been successfully uninstalled.")

	// i.frontend.Window.Close()

}

func removeFile(filePath string, file string) error {
	if fileExists(path.Join(filePath, file)) && file != "" {
		err := os.Remove(path.Join(filePath, file))
		if err != nil {
			return err
		}
	}
	return nil
}

func removeFiles(filePath string, files []string) error {
	for _, file := range files {
		if fileExists(path.Join(filePath, file)) && file != "" {
			err := os.Remove(path.Join(filePath, file))
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func removeFolders(folders []string) error {
	for _, folder := range folders {
		if fileExists(folder) && folder != "" {
			if fileExists(folder) {
				err := os.Remove(folder)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
