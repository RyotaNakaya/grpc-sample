package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	var filename string = "./helloworld_stream/client/list.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 1; i <= 1000000; i++ {
		s := "taro" + fmt.Sprint(i)
		_, err := file.WriteString(fmt.Sprintf("%s\n", s))
		if err != nil {
			log.Fatal(err)
		}
	}
}
