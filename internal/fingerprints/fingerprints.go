package fingerprints

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/drsigned/substko/pkg/substko"
)

// Update is a
func Update(fingerprintsPath string) (bool, error) {
	fingerprintsURL := "https://raw.githubusercontent.com/drsigned/substko/main/static/fingerprints.json"

	if _, err := os.Stat(fingerprintsPath); os.IsNotExist(err) {
		directory, _ := path.Split(fingerprintsPath)

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			if directory != "" {
				err = os.MkdirAll(directory, os.ModePerm)
				if err != nil {
					return false, err
				}
			}
		}
	}

	fingerprintsFile, err := os.Create(fingerprintsPath)
	if err != nil {
		return false, err
	}

	defer fingerprintsFile.Close()

	res, err := http.Get(fingerprintsURL)
	if err != nil {
		return false, err
	}

	if res.StatusCode != 200 {
		return false, errors.New("unexpected code")
	}

	defer res.Body.Close()

	_, err = io.Copy(fingerprintsFile, res.Body)
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
