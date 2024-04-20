package pagegencore

import (
	"os"
)

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFileToString(filePath string) string {
	data, err := os.ReadFile(filePath)
	CheckErr(err)

	return string(data)
}

func WriteStringToFile(filepath, contents string) {
	err := os.WriteFile(filepath, []byte(contents), 0644)
	CheckErr(err)
}

// Check if directory exists, and if not, create it
func CreateDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}
}
