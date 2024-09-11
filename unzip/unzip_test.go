package unzip

import (
	"fmt"
	"testing"
)

func TestUnzip(t *testing.T) {

	// 指定ZIP文件路径和解压目标路径
	zipFilePath := "d-mail-view.zip"
	extractToPath := "web"

	err := Unzip(zipFilePath, extractToPath)
	if err != nil {
		fmt.Println("Error while unzipping:", err)
	} else {
		fmt.Println("Unzipped successfully!")
	}

}
