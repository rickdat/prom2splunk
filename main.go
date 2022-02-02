package main

import (
	"fmt"
	"os"
	conv "promToSplunk/converter"
	"time"
)

func main() {
	go forever()
	select {} // block forever
}

func forever() {
	fmt.Println("Uploading events events")
	for {
		_, err := conv.UploadMetricJson()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		time.Sleep(time.Second * 20)
	}
}
