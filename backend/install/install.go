package install

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/artdarek/go-unzip"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	log "github.com/sirupsen/logrus"
)

// Install type contains the Install processes mandatory data
type Install struct {
	downloadURL        string
	dagFolderPath      string
	tmpFolderPath      string
	newVersion         string
	OSSpecificSettings *settings
}

type settings struct {
	osBuild       string
	fileExt       string
	binaryPath    string
	startMenuPath string
	desktopPath   string
	shortcutPath  string
}

type unzippedContents struct {
	mollyBinaryPath  string
	updateBinaryPath string
}

// Init initializes the Install struct
func Init() (*Install, error) {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to locate users home directory: %v", err)
	}

	i := &Install{
		downloadURL:        "https://github.com/grvlle/constellation_wallet/releases/download",
		newVersion:         "1.1.9",
		dagFolderPath:      path.Join(userHomeDir, ".dag"),
		tmpFolderPath:      path.Join(userHomeDir, ".tmp"),
		OSSpecificSettings: getOSSpecificSettings(),
	}
	return i, err
}

// Run is the main method that runs the full install.
func (i *Install) Run() {
	var err error

	if !fileExists(i.dagFolderPath) {
		err := os.Mkdir(i.dagFolderPath, os.FileMode(777))
		if err != nil {
			log.Fatalf("Unable to prepare filesystem: %v", err)
		}
	}

	initLogger()

	err = i.PrepareFS()
	if err != nil {
		log.Fatalf("Unable to prepare filesystem: %v", err)
	}

	zippedArchive, err := i.DownloadAppBinary()
	if err != nil {
		log.Fatalf("Unable to download v%s of Molly Wallet: %v", i.newVersion, err)
	}

	ok, err := i.VerifyChecksum(zippedArchive)
	if err != nil || !ok {
		log.Fatalf("Checksum missmatch. Corrupted download: %v", err)
	}

	contents, err := unzipArchive(zippedArchive, i.tmpFolderPath)
	if err != nil {
		log.Fatalf("Unable to unzip contents: %v", err)
	}

	err = i.ReplaceAppBinary(contents)
	if err != nil {
		log.Errorf("Unable to overwrite old installation: %v", err)

	}

	err = i.LaunchAppBinary()
	if err != nil {
		log.Errorf("Unable to start up Molly after Install: %v", err)
	}

	// Clean up install artifacts
	err = i.CleanUp()
	if err != nil {
		log.Fatalf("Unable to clear previous local state: %v", err)
	}

}

func initLogger() {
	// initialize update.log file and set log output to file
	file, err := os.OpenFile(path.Join(getDefaultDagFolderPath(), "install.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// PrepareFS removes uneccesary artifacts from the installation process and creates .dag folder if missing
func (i *Install) PrepareFS() error {

	// TODO: Loop through the files instead. Too tired to fix this rn :D

	if fileExists(i.dagFolderPath + "/store.db") {
		err := os.Remove(i.dagFolderPath + "/store.db")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/cl-keytool.jar.tmp") {
		err := os.Remove(i.dagFolderPath + "/cl-keytool.jar.tmp")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/cl-keytool.jar") {
		err := os.Remove(i.dagFolderPath + "/cl-keytool.jar")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/cl-wallet.jar") {
		err := os.Remove(i.dagFolderPath + "/cl-wallet.jar")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/cl-wallet.jar.tmp") {
		err := os.Remove(i.dagFolderPath + "/cl-wallet.jar.tmp")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/mollywallet.zip") {
		err := os.Remove(i.dagFolderPath + "/mollywallet.zip")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/mollywallet.zip.tmp") {
		err := os.Remove(i.dagFolderPath + "/mollywallet.zip.tmp")
		if err != nil {
			return err
		}
	}

	if fileExists(i.tmpFolderPath + "/new_build") {
		err := os.RemoveAll(i.tmpFolderPath + "/new_build")
		if err != nil {
			return err
		}
	}
	return nil
}

// DownloadAppBinary downloads the latest Molly Wallet zip from github releases and returns the path to it
func (i *Install) DownloadAppBinary() (string, error) {

	filename := "mollywallet.zip"

	if i.OSSpecificSettings.osBuild == "unsupported" {
		return "", fmt.Errorf("the OS is not supported")
	}

	url := i.downloadURL + "/v" + i.newVersion + "-" + i.OSSpecificSettings.osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/mollywallet.zip
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(i.dagFolderPath, filename)
	err := downloadFile(url, filePath)
	if err != nil {
		return "", fmt.Errorf("unable to download remote checksum: %v", err)
	}

	return filePath, nil
}

// VerifyChecksum takes a file path and will check the file sha256 checksum against the checksum included
// in the downlaod. Returns false if there's a missmatch.
func (i *Install) VerifyChecksum(filePathZip string) (bool, error) {
	// Download checksum
	filename := "checksum.sha256"

	if i.OSSpecificSettings.osBuild == "unsupported" {
		return false, fmt.Errorf("the OS is not supported")
	}

	url := i.downloadURL + "/v" + i.newVersion + "-" + i.OSSpecificSettings.osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/checksum.sha256
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(i.dagFolderPath, filename)
	err := downloadFile(url, filePath)
	if err != nil {
		return false, fmt.Errorf("unable to download remote checksum: %v", err)
	}

	// Read the contents of the downloaded file (remoteChecksum)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(content), "\n")
	remoteChecksum := lines[0]
	log.Infof("Remote file checksum: %v", remoteChecksum)

	// Collect the checksum of the downloaded zip (localChecksum)
	f, err := os.Open(filePathZip)
	if err != nil {
		return false, err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return false, err
	}
	localChecksum := hex.EncodeToString(hasher.Sum(nil))
	log.Infof("Local file checksum: %v", localChecksum)

	return remoteChecksum == localChecksum, nil
}

// ReplaceAppBinary replaces the old Install module and molly binary with the newly downloaded ones.
func (i *Install) ReplaceAppBinary(contents *unzippedContents) error {
	// Replace old molly binary with the new one

	err := copy(contents.mollyBinaryPath, i.OSSpecificSettings.binaryPath)
	if err != nil {
		for n := 5; n > 0; n-- {
			time.Sleep(time.Duration(n) * time.Second)
			err = copy(contents.mollyBinaryPath, i.OSSpecificSettings.binaryPath)
			if err == nil {
				break
			} else if err != nil && n == 0 {
				return fmt.Errorf("unable to move the molly binary: %v", err)
			}
		}
	}
	// Replace old update binary with the new one
	if fileExists(contents.updateBinaryPath) {
		err = copy(contents.updateBinaryPath, i.dagFolderPath+"/update"+i.OSSpecificSettings.fileExt)
		if err != nil {
			return fmt.Errorf("unable to copy update binary to .dag folder: %v", err)
		}
	}

	if runtime.GOOS == "windows" {
		err = createWindowsShortcuts(i.OSSpecificSettings.binaryPath, i.OSSpecificSettings.shortcutPath)
		if err != nil {
			return fmt.Errorf("unable to create app shortcut: %v", err)
		}
		err = copy(i.OSSpecificSettings.shortcutPath, path.Join(i.OSSpecificSettings.startMenuPath, "Molly Wallet.lnk"))
		if err != nil {
			return fmt.Errorf("unable to copy app shortcut to start menu: %v", err)
		}
		err = copy(i.OSSpecificSettings.shortcutPath, path.Join(i.OSSpecificSettings.desktopPath, "Molly Wallet.lnk"))
		if err != nil {
			return fmt.Errorf("unable to copy app shortcut to desktop: %v", err)
		}
	}

	return nil
}

// LaunchAppBinary executes the new molly binary
func (i *Install) LaunchAppBinary() error {
	cmd := exec.Command(i.OSSpecificSettings.binaryPath)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to execute run command for Molly Wallet: %v", err)
	}
	return nil
}

// CleanUp removes uneccesary artifacts from the Install process
func (i *Install) CleanUp() error {

	if fileExists(i.dagFolderPath + "/mollywallet.zip") {
		err := os.Remove(i.dagFolderPath + "/mollywallet.zip")
		if err != nil {
			return err
		}
	}

	if fileExists(i.dagFolderPath + "/checksum.sha256") {
		err := os.Remove(i.dagFolderPath + "/checksum.sha256")
		if err != nil {
			return err
		}
	}

	if fileExists(i.tmpFolderPath + "/new_build") {
		err := os.RemoveAll(i.tmpFolderPath + "/new_build")
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(url, filePath string) error {

	tmpFilePath := filePath + ".tmp"
	out, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return err
	}

	out.Close() // Close file to rename

	if err = os.Rename(tmpFilePath, filePath); err != nil {
		return err
	}
	return nil
}

func getDefaultDagFolderPath() string {
	userDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("Unable to detect UserHomeDir: %v", err)
		return ""
	}
	return userDir + "/.dag"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !os.IsNotExist(err)
}

func copy(src string, dst string) error {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Unzips archive to dstPath, returns path to wallet binary
func unzipArchive(zippedArchive, dstPath string) (*unzippedContents, error) {

	uz := unzip.New(zippedArchive, path.Join(dstPath, "new_build"))
	err := uz.Extract()
	if err != nil {
		return nil, err
	}
	var fileExt string
	if runtime.GOOS == "windows" {
		fileExt = ".exe"
	}

	contents := &unzippedContents{
		mollyBinaryPath:  path.Join(dstPath, "new_build", "mollywallet"+fileExt),
		updateBinaryPath: path.Join(dstPath, "new_build", "update"+fileExt),
	}

	return contents, err
}

// getUserOS returns the users OS, the file extension of executables and path to put molly wallet binary for said OS
func getOSSpecificSettings() *settings {

	s := &settings{}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Errorf("unable to locate home dir: %v", err)
	}

	switch os := runtime.GOOS; os {

	case "darwin":
		s = &settings{
			osBuild:    "darwin",
			fileExt:    "",
			binaryPath: path.Join("usr", "local", "bin", "mollywallet"),
		}

	case "linux":
		s = &settings{
			osBuild:    "linux",
			fileExt:    "",
			binaryPath: path.Join("usr", "local", "bin", "mollywallet"),
		}

	case "windows":
		s = &settings{
			osBuild:       "windows",
			fileExt:       ".exe",
			binaryPath:    path.Join(getDefaultDagFolderPath(), "mollywallet.exe"),
			startMenuPath: path.Join(homeDir, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs"),
			desktopPath:   path.Join(homeDir, "Desktop"),
			shortcutPath:  getDefaultDagFolderPath() + "/Molly Wallet.lnk",
		}

	default:
		s = &settings{
			osBuild:    "unsupported",
			fileExt:    "",
			binaryPath: getDefaultDagFolderPath(),
		}

	}

	return s
}

func createWindowsShortcuts(src, dst string) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	defer ole.CoUninitialize()

	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer oleShellObject.Release()
	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer wshell.Release()
	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", dst)
	fmt.Println(dst + ".lnk")
	if err != nil {
		return err
	}
	idispatch := cs.ToIDispatch()
	oleutil.PutProperty(idispatch, "TargetPath", src)
	oleutil.CallMethod(idispatch, "Save")

	return nil
}

// setEnviornment sets env vars permenantly on Windows, but requires administrator access.
// func setEnvironment(key string, value string) {
// 	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\Session Manager\Environment`, registry.ALL_ACCESS)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer k.Close()

// 	err = k.SetStringValue(key, value)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
