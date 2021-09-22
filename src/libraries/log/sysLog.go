package main

import (
	"fmt"
	"log"
	"log/syslog"
)

func main() {
	process := "example"
	sysLog, e := syslog.New(syslog.LOG_NOTICE|syslog.LOG_MAIL, process)
	// 로그 레벨을 notice로, 로그 종류는 mail로 설정해서 syslog.Writer를 만들어준다.
	// 두 번째 인자 값은, 로그 파일에서 필요한 정보를 얻기 위해 실행 중인 프로세스의 이름이 들어가는 것이 좋다.

	if e != nil {
		log.Fatal(e)
	}
	sysLog.Crit("Crit: Logging in Go!") // 로그 레벨을 Writer에 지정하긴 했으나, 다른 우선 순위로 메시지를 보내는 것도 허용함

	sysLog, e = syslog.New(syslog.LOG_ALERT|syslog.LOG_LOCAL7, "Some program!") // 같은 프로그램 내에서 sysLog Writer를 여러번 만들 수 있다.
	if e != nil {
		log.Fatal(sysLog)
	}
	sysLog.Emerg("Emerg: Logging in Go!")

	fmt.Fprintf(sysLog, "log.Print: Logging in Go!") // Writer에 스트링을 씀으로써, New의 Writer에서 지정한 메시지 형식으로 로그 기록
}
