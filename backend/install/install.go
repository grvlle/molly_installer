package install

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails"
)

// Install type contains the Install processes mandatory data
type Install struct {
	downloadURL         string
	dagFolderPath       string
	tmpFolderPath       string
	incrementProgressCh chan int
	progressMessageCh   chan string
	OSSpecificSettings  *settings
	frontend            *wails.Runtime
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
	mollyMacOSApp    string
}

// Init initializes the Install struct
func Init() (*Install, error) {

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to locate users home directory: %v", err)
	}

	i := &Install{
		downloadURL:         "https://github.com/grvlle/constellation_wallet/releases/download",
		dagFolderPath:       path.Join(userHomeDir, ".dag"),
		tmpFolderPath:       path.Join(userHomeDir, ".tmp"),
		incrementProgressCh: make(chan int),
		progressMessageCh:   make(chan string),
		OSSpecificSettings:  getOSSpecificSettings(),
	}
	return i, err
}

// Run is the main method that runs the full install.
func (i *Install) Run() {
	var err error

	go i.startProgress() // Runs a go routine that increments the progress bar

	// Install Java on Windows if not detected
	i.updateProgress(8, "Checking Java Installation...")
	if runtime.GOOS == "windows" && !javaInstalled() {
		i.updateProgress(10, "Java not found. Installing Java (This may take some time)...")
		err = installJava()
		if err != nil {
			i.sendErrorNotification("Unable to install Java", fmt.Sprintf("%v", err))
			time.Sleep(10 * time.Second)
			log.Fatal("Unable to install Java: %v", err)
		}
	}

	// Remove old Molly Wallet artifacts
	i.updateProgress(33, "Preparing filesystem...")
	err = i.PrepareFS()
	if err != nil {
		i.sendErrorNotification("Unable to prepare filesystem", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Fatalf("Unable to prepare filesystem: %v", err)
	}

	// Download the mollywallet.zip from https://github.com/grvlle/constellation_wallet/
	i.updateProgress(35, "Downloading packages...")
	zippedArchive, err := i.DownloadAppBinary()
	if err != nil {
		i.sendErrorNotification("Unable to download Molly Wallet package", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Fatalf("Unable to download Molly Wallet package: %v", err)
	}

	i.updateProgress(42, "Downloading the wallet SDK...")
	err = i.checkAndFetchWalletCLI()
	if err != nil {
		i.sendErrorNotification("Unable to download CL files", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Errorf("Unable to download CL files: %v", err)
	}

	// Verify the integrity of the package
	i.updateProgress(86, "Verifying Checksum...")
	ok, err := i.VerifyChecksum(zippedArchive)
	if err != nil || !ok {
		i.sendErrorNotification("Checksum missmatch. Corrupted download", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Fatalf("Checksum missmatch. Corrupted download: %v", err)
	}

	// Extract the contents
	i.updateProgress(95, "Exctracting contents...")
	contents, err := unzipArchive(zippedArchive, i.tmpFolderPath)
	if err != nil {
		i.sendErrorNotification("Unable to unzip contents", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Fatalf("Unable to unzip contents: %v", err)
	}

	// Copy the contents (mollywallet and update) to the .dag folder
	i.updateProgress(98, "Copy binaries...")
	err = i.CopyAppBinaries(contents)
	if err != nil {
		i.sendErrorNotification("Unable to overwrite old installation", fmt.Sprintf("%v", err))
		log.Errorf("Unable to overwrite old installation: %v", err)

	}

	i.updateProgress(100, "Installation Complete! Launching Molly Wallet...")
	i.sendSuccessNotification("Success!", "Molly wallet has been successfully installed.")
	time.Sleep(5 * time.Second)

	// Lauch mollywallet
	err = i.LaunchAppBinary()
	if err != nil {
		i.sendErrorNotification("Unable to start up Molly after Install", fmt.Sprintf("%v", err))
		log.Errorf("Unable to start up Molly after Install: %v", err)
	}

	// Clean up install artifacts
	err = i.CleanUp()
	if err != nil {
		i.sendErrorNotification("Unable to clear previous local state", fmt.Sprintf("%v", err))
		time.Sleep(10 * time.Second)
		log.Fatalf("Unable to clear previous local state: %v", err)
	}

	i.frontend.Window.Close()

}

// PrepareFS removes uneccesary artifacts from the installation process and creates .dag folder if missing
func (i *Install) PrepareFS() error {
	// files slice will house the files that are to be removed before proceeding with installation.
	files := make([]string, 8)
	files = append(files, "cl-keytool.jar.tmp", "cl-keytool.jar", "cl-wallet.jar", "cl-wallet.jar.tmp", "mollywallet.zip", "mollywallet.zip.tmp", "Molly Wallet.lnk", "mollywallet.exe")

	removeFiles(i.dagFolderPath, files)

	// in case of a failed previous installation attempt, there may be extracted artifacts in .tmp
	if fileExists(i.tmpFolderPath) {
		err := os.RemoveAll(i.tmpFolderPath)
		if err != nil {
			return err
		}
	}

	// remove the old .dag folder
	folders := make([]string, 1)
	folders = append(folders, path.Join(i.dagFolderPath))

	if len(folders) != 0 {
		err := removeFolders(folders)
		if err != nil {
			i.sendErrorNotification("Error:", convertErrorToString(err))
			log.Errorf("Error: %v", err)
		}
	}

	// create a new .dag folder with the right permissions
	if !fileExists(i.dagFolderPath) {
		err := os.Mkdir(i.dagFolderPath, os.FileMode(774))
		if err != nil {
			i.sendErrorNotification("Unable to prepare filesystem", fmt.Sprintf("%v", err))
			time.Sleep(10 * time.Second)
			log.Fatalf("Unable to prepare filesystem: %v", err)
		}
	}

	// Remove the .app folder on MacOS
	if runtime.GOOS == "darwin" && fileExists(i.OSSpecificSettings.shortcutPath) {
		err := os.RemoveAll(i.OSSpecificSettings.shortcutPath)
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadAppBinary downloads the latest Molly Wallet zip from github releases and returns the path to it
func (i *Install) DownloadAppBinary() (string, error) {

	filename := "mollywallet.zip"
	version, err := i.getLatestRelease()
	if err != nil {
		return "", err
	}

	if i.OSSpecificSettings.osBuild == "unsupported" {
		return "", fmt.Errorf("the OS is not supported")
	}

	url := i.downloadURL + "/v" + version + "-" + i.OSSpecificSettings.osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/mollywallet.zip
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(i.dagFolderPath, filename)
	err = downloadFile(url, filePath)
	if err != nil {
		return "", fmt.Errorf("unable to download remote checksum: %v", err)
	}

	return filePath, nil
}

// CheckAndFetchWalletCLI will download the cl-wallet dependencies from
// the official Constellation Repo
func (i *Install) checkAndFetchWalletCLI() error {

	var downloadComplete bool

	keytoolPath := path.Join(i.dagFolderPath, "cl-keytool.jar")
	walletPath := path.Join(i.dagFolderPath, "cl-wallet.jar")

	err := i.fetchWalletJar("cl-keytool.jar", keytoolPath)
	if err != nil {
		log.Errorln("Unable to fetch or store cl-keytool.jar", err)
		return err
	}

	err = i.fetchWalletJar("cl-wallet.jar", walletPath)
	if err != nil {
		log.Errorln("Unable to fetch or store cl-wallet.jar", err)
		return err
	}

	if fileExists(keytoolPath) && fileExists(walletPath) {
		downloadComplete = true
	} else {
		downloadComplete = false
	}

	if !downloadComplete {
		err := errors.New("download failed")
		return err
	}
	return err

}

func (i *Install) fetchWalletJar(filename string, filePath string) error {
	url := "https://github.com/Constellation-Labs/constellation/releases/download/v2.6.0/" + filename
	log.Infof("Constructed the following URL: %s", url)

	filePath = path.Join(i.dagFolderPath, filename)
	err := downloadFile(url, filePath)
	if err != nil {
		return fmt.Errorf("unable to download remote checksum: %v", err)
	}

	return err
}

// VerifyChecksum takes a file path and will check the file sha256 checksum against the checksum included
// in the downlaod. Returns false if there's a missmatch.
func (i *Install) VerifyChecksum(filePathZip string) (bool, error) {

	// Download checksum
	filename := "checksum.sha256"
	version, err := i.getLatestRelease()
	if err != nil {
		return false, err
	}

	if i.OSSpecificSettings.osBuild == "unsupported" {
		return false, fmt.Errorf("the OS is not supported")
	}

	url := i.downloadURL + "/v" + version + "-" + i.OSSpecificSettings.osBuild + "/" + filename
	// e.g https://github.com/grvlle/constellation_wallet/releases/download/v1.1.9-linux/checksum.sha256
	log.Infof("Constructed the following URL: %s", url)

	filePath := path.Join(i.dagFolderPath, filename)
	err = downloadFile(url, filePath)
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

// CopyAppBinaries copies the update module and molly binary from the unzipped package to the .dag folder.
func (i *Install) CopyAppBinaries(contents *unzippedContents) error {
	err := copyFile(contents.mollyBinaryPath, i.OSSpecificSettings.binaryPath)
	if err != nil {
		for n := 5; n > 0; n-- {
			time.Sleep(time.Duration(n) * time.Second)
			err = copyFile(contents.mollyBinaryPath, i.OSSpecificSettings.binaryPath)
			if err == nil {
				break
			} else if err != nil && n == 0 {
				return fmt.Errorf("unable to move the molly binary: %v", err)
			}
		}
	}
	// Replace old update binary with the new one
	if fileExists(contents.updateBinaryPath) {
		err = copyFile(contents.updateBinaryPath, i.dagFolderPath+"/update"+i.OSSpecificSettings.fileExt)
		if err != nil {
			return fmt.Errorf("unable to copy update binary to .dag folder: %v", err)
		}
	}
	if runtime.GOOS == "darwin" {
		err := copy.Copy(contents.mollyMacOSApp, i.OSSpecificSettings.shortcutPath, copy.Options{AddPermission: 0774})
		if err != nil {
			return fmt.Errorf("unable to copy Molly - Constellation Desktop Wallet.app to Applications folder: %v", err)
		}
	}

	if runtime.GOOS == "windows" {
		err = createWindowsShortcuts(i.OSSpecificSettings.binaryPath, i.OSSpecificSettings.shortcutPath)
		if err != nil {
			return fmt.Errorf("unable to create app shortcut: %v", err)
		}
		err = copyFile(i.OSSpecificSettings.shortcutPath, path.Join(i.OSSpecificSettings.startMenuPath, "Molly Wallet.lnk"))
		if err != nil {
			return fmt.Errorf("unable to copy app shortcut to start menu: %v", err)
		}
		err = copyFile(i.OSSpecificSettings.shortcutPath, path.Join(i.OSSpecificSettings.desktopPath, "Molly Wallet.lnk"))
		if err != nil {
			return fmt.Errorf("unable to copy app shortcut to desktop: %v", err)
		}
	}

	return nil
}

// LaunchAppBinary executes the new molly binary
func (i *Install) LaunchAppBinary() error {
	cmd := exec.Command(i.OSSpecificSettings.binaryPath)

	if runtime.GOOS == "darwin" {
		cmd = exec.Command(path.Join(i.OSSpecificSettings.shortcutPath, "Contents", "MacOS", "mollywallet"))
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to execute run command for Molly Wallet: %v", err)
	}
	return nil
}

// CleanUp removes uneccesary artifacts from the Install process
func (i *Install) CleanUp() error {

	files := make([]string, 2)
	files = append(files, "mollywallet.zip", "checksum.sha256")

	removeFiles(i.dagFolderPath, files)

	if fileExists(i.tmpFolderPath) {
		err := os.RemoveAll(i.tmpFolderPath)
		if err != nil {
			return err
		}
	}
	return nil
}
