package install

import (
	"errors"
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
)

// CheckAndFetchWalletCLI will download the cl-wallet dependencies from
// the official Constellation Repo
func (i *Install) checkAndFetchWalletCLI() error {
	i.updateProgress(62, "Downloading the wallet SDK...")

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
		err := errors.New("Download failed")
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
