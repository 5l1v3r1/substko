package targets

import (
	"bufio"
	"os"
)

// Load is a
func Load(targetsPath string) (targets []string, err error) {
	var scanner *bufio.Scanner

	if targetsPath != "" {
		targetsFile, err := os.Open(targetsPath)
		if err != nil {
			return targets, err
		}

		defer targetsFile.Close()

		scanner = bufio.NewScanner(targetsFile)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		targets = append(targets, scanner.Text())
	}

	if scanner.Err() != nil {
		return targets, scanner.Err()
	}

	return targets, nil
}
