package provider

import (
	"strings"
)

// shaFileToMap converts the contents of a SHA checksums file into a map of file names to SHA checksums.
func shaFileToMap(contents []byte) map[string]string {
	var result = make(map[string]string)
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		shaSum := parts[0]
		fileName := parts[1]

		result[fileName] = shaSum
	}
	return result
}
