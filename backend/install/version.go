package install

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
