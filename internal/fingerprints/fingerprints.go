package fingerprints

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/drsigned/substko/pkg/substko"
)

// Update is a
func Update(fingerprintsPath string) (bool, error) {
	fingerprintsURL := "https://raw.githubusercontent.com/drsigned/substko/main/static/fingerprints.json"

	fingerprintsFile, err := os.Create(fingerprintsPath)
	if err != nil {
		return false, err
	}

	defer fingerprintsFile.Close()

	resp, err := http.Get(fingerprintsURL)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != 200 {
		return false, errors.New("unexpected code")
	}

	defer resp.Body.Close()

	_, err = io.Copy(fingerprintsFile, resp.Body)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Load is a
func Load(fingerprintsPath string) (fingerprints []substko.Fingerprint, err error) {
	rawFingerprints, err := ioutil.ReadFile(fingerprintsPath)
	if err != nil {
		return fingerprints, err
	}

	err = json.Unmarshal(rawFingerprints, &fingerprints)
	if err != nil {
		return fingerprints, err
	}

	return fingerprints, nil
}
