package install

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/artdarek/go-unzip"
	log "github.com/sirupsen/logrus"
)

// getLatestRelease polls the github api for the latest release in the constellation_wallet repo
// and returns the sem ver and error
func (i *Install) getLatestRelease() (string, error) {

	const (
		url = "https://api.github.com/repos/grvlle/constellation_wallet/releases/latest"
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	if resp == nil {
		return "", fmt.Errorf("empty response from Github API: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return "", fmt.Errorf("unable to parse GitHub API response: %v", err)
	}

	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}

	release := result["tag_name"]
	bytes := []byte(release.(string))
	version := string(bytes[1:6])
	return version, err

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
