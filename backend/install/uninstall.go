package install

import (
	"path"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Uninstall will Uninstall the application by removing the files located
// in the .dag directory:
//   'Molly Wallet.lnk'   cl-wallet.jar   mollywallet.exe   tmp
//   cl-keytool.jar      install.log     store.db          wallet.log
//   update.exe
// And also removing the shortcuts on Windows.
func (i *Install) Uninstall() {

	updateBinary := "update" + i.OSSpecificSettings.fileExt

	files := make([]string, 13)
	files = append(files, "update.log", updateBinary, "wallet.log", "store.db", "cl-keytool.jar.tmp", "cl-keytool.jar", "cl-wallet.jar", "cl-wallet.jar.tmp", "mollywallet.zip", "mollywallet.zip.tmp", "Molly Wallet.lnk", "mollywallet.exe")

	log.Infoln("Removing dependencies...")
	err := removeFiles(i.dagFolderPath, files)
	if err != nil {
		i.sendErrorNotification("Error:", convertErrorToString(err))
		log.Errorf("Error: %v", err)
	}

	folders := make([]string, 3)
	folders = append(folders, path.Join(i.dagFolderPath, "tmp"), i.tmpFolderPath, i.dagFolderPath)

	err = removeFolders(folders)
	if err != nil {
		i.sendErrorNotification("Error:", convertErrorToString(err))
		log.Errorf("Error: %v", err)
	}

	log.Infoln("Removing shortcuts on Windows...")
	if runtime.GOOS == "windows" {
		err := removeFile(i.OSSpecificSettings.startMenuPath, "Molly Wallet.lnk")
		if err != nil {
			i.sendErrorNotification("Unable to remove shortcut from start menu", convertErrorToString(err))
			log.Errorf("Error: %v", err)
		}
		err = removeFile(i.OSSpecificSettings.desktopPath, "Molly Wallet.lnk")
		if err != nil {
			i.sendErrorNotification("Unable to remove shortcut from desktop", convertErrorToString(err))
			log.Errorf("Error: %v", err)
		}
	}

	i.sendSuccessNotification("Success!", "Molly wallet has been successfully uninstalled.")

	// i.frontend.Window.Close()

}

// strip non-regex complient chars and return clean error string
func convertErrorToString(err error) string {
	if err == nil {
		return ""
	}

	chars := []string{"]", "^", "\\\\", "[", ".", "(", ")", "-"}
	r := strings.Join(chars, "")
	errString := err.Error()
	re := regexp.MustCompile("[" + r + "]+")
	errString = re.ReplaceAllString(errString, "")

	if len(errString) >= 15 {
		bytes := []byte(errString)
		return string(bytes[0:15])
	}

	return errString
}
