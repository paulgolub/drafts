package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
  _ "image/gif"
	"os"
	"path/filepath"

	"github.com/liyue201/goqr"
)

func main() {
	dir := "./image"

  // read files from dir
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error, can not read a directory: ", err)
		return
	}

	// files only
	for _, file := range files {
		// skip dir
		if file.IsDir() {
			continue
		}

		// full path
		filePath := filepath.Join(dir, file.Name())

		f, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Error, can not open a file %s: %v\n", filePath, err)
			continue
		}
		defer f.Close()

		// decoding
		img, _, err := image.Decode(f)
		if err != nil {
			fmt.Printf("Error, can not decod %s: %v\n", filePath, err)
			continue
		}

		// QR recognizion
		codes, err := goqr.Recognize(img)
		if err != nil {
			fmt.Printf("QR-code nof found in file %s: %v\n", filePath, err)
			continue
		}

		for _, code := range codes {
			fmt.Printf("Finf QR-code in file %s: %s\n", filePath, string(code.Payload))
		}
	}
}
