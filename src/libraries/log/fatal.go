package main

import (
	"fmt"
	"log"
	"log/syslog"
)

func main() {
	sysLog, err := syslog.New(syslog.LOG_ALERT|syslog.LOG_MAIL, "process")
	if err != nil {
		log.Fatal(err)
	} else {
		log.SetOutput(sysLog)
	}

	log.Fatal("Fatal error on process") // 이 메시지가 /var/log/mail.log에 찍혀있다.
	fmt.Println("You can't see me")     // 이 부분은 콘솔 위에서 확인할 수 없다.
}
