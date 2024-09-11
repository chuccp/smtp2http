package unzip

import (
	"fmt"
	"strings"
	"testing"
)

func TestUnzip(t *testing.T) {

	xx := "d-mail-view.zip web"

	vs := strings.Split(xx, " ")

	println(vs[0])
	println(vs[1])
	// 指定ZIP文件路径和解压目标路径
	zipFilePath := vs[0]
	extractToPath := vs[1]

	err := Unzip(zipFilePath, extractToPath)
	if err != nil {
		fmt.Println("Error while unzipping:", err)
	} else {
		fmt.Println("Unzipped successfully!")
	}

}
