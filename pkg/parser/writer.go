package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFile(out []byte, fileName string) error {
	shouldCreateFolder, folderPath := removeBaseFromPath(fileName)
	if shouldCreateFolder {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not create folder. Err: %v\n", err)
			return err
		}
	}

	err := ioutil.WriteFile(fileName, out, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to file. Err: %v\n", err)
		return err
	}

	fmt.Fprintf(os.Stdout, "Generated file with name %s\n", fileName)
	return nil
}

func removeBaseFromPath(f string) (shouldCreateFolder bool, folderPath string) {
	folderPath = ""
	shouldCreateFolder = true
	splitted := strings.Split(f, "/")
	if len(splitted) <= 1 {
		shouldCreateFolder = false
		return
	}
	folderPath = strings.Join(splitted[0:len(splitted)-1], "/")
	return
}
