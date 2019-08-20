package bookcache

import (
	"fmt"
	"io/ioutil"
	"time"
)

var FILE_PATH string = "d:\\test_data\\"

func ReadFile(book_id string) (contents string, err error) {
	fmt.Printf("Disk I/O ReadFile() book_id=%v time=%v\n", book_id, time.Now().Unix())
	buf, err := ioutil.ReadFile(FILE_PATH + book_id)
	if err != nil {
		return
	}
	contents = string(buf[:])
	return
}
func SaveFile(book_id, contents string) (err error) {
	fmt.Printf("Disk I/O SaveFile() book_id=%v time=%v\n", book_id, time.Now().Unix())
	err = ioutil.WriteFile(FILE_PATH+book_id, []byte(contents), 0644)
	return
}
