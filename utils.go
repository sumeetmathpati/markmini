package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// func ReadFileData(r io.Reader) ([]byte, error) {
// 	data, err := ioutil.ReadAll(r)
// 	return data, err
// }

func isDir(path string) (bool, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open file:", path, err)
		return false, errors.New(fmt.Sprint("Couldn't open file:", path))
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't open file:", path, err)
		return false, errors.New(fmt.Sprint("Couldn't open file:", path))
	}

	return fileInfo.IsDir(), nil

}

func getHomeDirOrFail() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not get home direcltory")
	}

	return homeDir
}
