---
title: log
---

`log` 패키지는 로그 메시지를 유닉스 시스템 로그로 보내준다. 이번 장에서는 로그 수준과 종류를 설정해 로그를 시스템으로 보내는 방법을 정리하고 있다.

---

## Unix 시스템 로그

일반적으로 Unix 시스템의 로그 파일은 대부분 `/var/log` 디렉토리에 있다. 로깅 해야 하는 정보들을 커맨드라인에 표시하지 않고, 파일에 기록하는 것이 좋은데, 그 이유는 다음과 같다.

1. 영구 지속
2. 유닉스 도구의 지원을 받아 검색 및 처리 가능

### 로그 서버
유닉스 시스템에서는 로그 파일을 로깅하는 분리된 서버 프로세스가 존재한다. macOS의 경우 해당 프로세스 이름이 `syslogd`이다. 다른 대부분의 리눅스 시스템에서는 `rsyslogd`를 사용한다. `syslogd`보다 더 안정적이고 발전된 형태라고 한다.

### 로그 수준 (Log Level)
로그 내용의 응급도를 표현하는 값이다. `debug`, `init`, `notice`, `warning`, `err`, `crit`, `alert`, `emerg` 등이 있다. 쓰여진 순서대로 중요도가 높은 로그를 나타낸다.

### 로그 종류 (Log Facilities)
로그 종류는 로그가 포함되어있는 카테고리를 의미한다. `auth`, `authpriv`, `cron`, `daemon`, `kern`, `lpr`, `mail`, `user`, `local0` ~ `local7` 등이 있다. 유닉스 서버에 `/etc/rsyslog.conf` 또는 `/etc/rsyslog.d/*.conf` 형태의 파일에서 종류들과 어떤 로그파일로 향하는지를 설정값으로 정해두었다.

## log 패키지

Go에서 시스템 로그 파일에 메시지를 쓰는 방법을 확인해보자.

```go
package main

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
	"path/filepath"
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

	sysLog, e = syslog.New(syslog.LOG_ALERT|syslog.LOG_LOCAL7, "another " +  process) // 같은 프로그램 내에서 sysLog Writer를 여러번 만들 수 있다.
	if e != nil {
		log.Fatal(sysLog)
	}
	sysLog.Emerg("Emerg: Logging in Process!")

	fmt.Fprintf(sysLog, "log.Print: Logging in Process!") // Writer에 스트링을 씀으로써, New의 Writer에서 지정한 메시지 형식으로 로그 기록
}
```

`syslog.New` 함수를 사용해 주어진 인자의 로그 레벨과 종류로 설정된 Writer 타입을 만들어낸다. 위 예시에서는 `fmt.Fprintf`를 사용해 스트링을 적었지만, `log.SetOutput` 함수를 사용해 기본 로거를 설정해준 뒤, `log.Println`을 사용해 로그 서버로 내용을 보내는 방법도 있다.

```go 
    ...
    if e != nil {
        log.Fatal(e)
    } else {
        log.SetOutput(sysLog)
    }

    log.Println("log.Print: Logging in Go!")
```

결과는 브로드캐스팅 되기도 하고, 설정 파일에 지정된 파일들에 쌓이게 된다. 하나의 파일에 쌓일 수도 있지만 아닌 경우가 더 많다. 로그 레벨과 종류에 따라 여러 파일들에게 모두 중복해서 쓴다.

### log.Fatal

즉시 프로그램을 종료하고자 할 때 사용한다. 다만 종료시 설정된 로그 레벨과 로그 종류에 맞게 로그를 로그 서버에 보낸다.

```go
// src/packages/log/fatal.go

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
```

### log.Panic

프로세스가 심각한 에러로 인해 죽는 상황에서, 관련한 정보를 최대한 볼 수 있는 방식이다. 호출 스택과 유관한 정보들을 포함해서 로그를 보여준다.

```go
// src/packages/log/panic.go
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

	log.Panic("fatal error on process")
	fmt.Println("You can't see me") //마찬가지로 볼 수가 없다.
}
```

---

## 커스텀 로그

시스템에 설정된 로그 파일 외에 원하는 파일에 로그를 보내야 할 때가 있다. 예를 들어, 독립적인 로그 정보를 수집하려고 한다든지, 다른 포맷을 쓰고 싶다든지 등의 이유가 있다. 일반적인 파일 쓰기와 유사하지만, 로거를 만들어서 사용하는 부분에서 차이가 있다.

```go
// src/packages/log/custom.go
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
```

{{< expand "결과" >}}
```sh
$ vi /tmp/custom.log

CUSTOM:2021/09/18 00:56:43 custom log
```
{{< /expand >}}

{{< hint info>}}
**파일 만들기에 대해서**  
`os.OpenFile`은 첫 번째 인자로 파일 경로, 두 번째 인자로 파일 열기 모드를 설정하는 플래그, 마지막은 만약 생성한다면 적용되는 파일 권한이다.
{{< /hint >}}

`log.New`는 새로운 로거를 만들어준다. 첫 번째는 파일 쓰기 작업을 위한 `io.Writer` 타입이 필요하다. 두 번째는 로그 앞에 붙는 접두사를 설정하는 곳이고, 세 번째는 로그 파일의 특성(추가 정보) 속성이다. `log.LstdFlags`는 기본 로그의 로컬 타임과 로컬 데이트를 포함하는 플래그이다.

![다음과 같은 플래그들이 있다.](/dive-into-go/images/packages/log/custom/01.png)

위와 같은 플래그들이 있는데, 여러 플래그를 합성하기 위해서는 `log.LstdFlags | log.Lshortfile`과 같이 비트 OR 연산자로 합성할 수 있다. 예시와 같이 합성하면, 파일 이름과 파일 줄 수를 로그 파일에 추가해준다.