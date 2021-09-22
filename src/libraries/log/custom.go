package main

import (
	"log"
	"os"
)

const LOGFILE = "/tmp/custom.log"

func main() {
	f, err := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger := log.New(f, "CUSTOM:", log.LstdFlags)
	// New의 첫 번째 인자는 파일 Writer 타입, 두 번째는 접두사, 세 번째는 로그의 특성들을 추가해준다.
	// LstdFlags는 Local Date, Local Time이 있는 기본 로그 포맷이다.

	logger.Println("custom log")
}
