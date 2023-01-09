package utils

import "os"

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func ReadFile(filePaht string) []byte {
	if IsFileExist(filePaht) {
		content, err := os.ReadFile(filePaht)
		if err != nil {
			os.Exit(-1)
		}
		return content
	}

	return nil
}
