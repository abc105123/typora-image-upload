package utils

import (
	"os"
)

func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func ReadFile(filePath string) []byte {
	if IsFileExist(filePath) {
		content, err := os.ReadFile(filePath)
		if err != nil {
			os.Exit(-1)
		}
		return content
	}

	return nil
}

func WriteFile(filePath string, data []byte) bool {
	//path := strings.Replace(filePath, "\\", "/", -1)
	if !IsFileExist(filePath) {
		_, err := os.Create(filePath)
		if err != nil {
			os.Exit(-1)
		}
	}

	file, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm) // 0666?
	defer file.Close()

	//writer := bufio.NewWriter(file)
	//for i, _ := writer.Write(data); i > 0; {
	//	//...
	//}
	//writer.Flush()

	_, err := file.Write(data)
	if err != nil {
		return false
	}

	return true
}
