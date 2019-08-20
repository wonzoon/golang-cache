package main

import (
	"fmt"
	"strconv"

	// "io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func createFile(path string) {
	// check if file exists
	var _, err = os.Stat(path)
	// create if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}
	fmt.Println("file created : ", path)
}
func writeFile(path string) {
	// Open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	var builder strings.Builder
	builder.WriteString("{")
	const NUM = 100
	for i := 0; i < NUM; i++ {
		builder.WriteString("\"")
		builder.WriteString(strconv.Itoa(i))
		builder.WriteString("\"")
		builder.WriteString(":")
		builder.WriteString("\"")
		builder.WriteString(filepath.Base(file.Name()))
		builder.WriteString("\"")
		if i < NUM-1 {
			builder.WriteString(",\n")
		}

	}
	builder.WriteString("}")
	_, err = file.WriteString(builder.String())
	if err != nil {
		return
	}
	// Save
	err = file.Sync()
	if err != nil {
		return
	}

	fmt.Println("File Write")
}

//create cvs
func main() {

	files, err := ioutil.ReadDir("d:\\test_data")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f, err := os.OpenFile("d:\\go-work\\test_10k_rep3.cvs",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer f.Close()

	var count int = 10 * 10000
	for _, file := range files {
		if count < 1 {
			break
		}
		if _, err := f.WriteString(file.Name() + ",\n"); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(file.Name())
		count--
	}
	fmt.Println("Done!")
}

//create data file
// func main() {
// 	for i := 0; i < 300*10000; i++ {
// 		fileName := StringWithCharset(19, "0987654321")
// 		path := "d:\\test_data\\" + fileName
// 		createFile(path)
// 		writeFile(path)
// 	}
// 	fmt.Println("Done.")
// }
