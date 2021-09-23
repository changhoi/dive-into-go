---
title: Module
---

Go는 1.11부터 모듈에 지원을 시작했다. Module은 Go의 새로운 의존성 관리 시스템으로, 의존성의 버전 정보를 명확하고 관리하기 쉽게 만든다. 모듈은 파일 트리에 저장되는 **Go 패키지 모음**으로, 루트에 `go.mod` 파일이 있다. Go 1.11부터, 진행 중인 디렉토리가 `GOPATH` 변수의 `src` 아래가 아닌 경우, go 명령어는 현재 디렉토리나 또는 상위 디렉토리가 `go.mod` 파일을 가지고 있을 때 모듈을 사용하도록 했다. 그리고 1.13 이후 모든 개발 환경에서 기본 모드로 모듈을 사용했다.

## 새로운 모듈 만들기

GOPATH가 아닌 곳에 새로운 디렉토리를 만들고 패키지를 만들자.

```go
package hello

func Hello() string {
	return "안녕, 세상."
}
```

위 상태에서 디렉토리 `hello`는 패키지를 담고 있지만, 모듈은 아니다. 왜냐면 `go.mod` 파일이 없기 때문이다. `hello` 디렉토리를 모듈의 루트로 만들기 위해서는 `go mod init` 명령어를 수행하면 된다. `go mod init` 명령어 인자로 패키지 이름이 들어가야 하는데, 일반적으로 패키지가 공유되는 URL 패스를 사용한다. `example.com/hello` 라는 이름으로 패키지를 만들었다.

```shell
$ go mod init example.com/hello
```

경과로 `go.mod` 파일을 얻을 수 있는데 처음 만들어진 `go.mod` 파일은 다음과 같이 생겼다.
``` go
module example.com/hello

go 1.17
```
`go.mod` 파일은 모듈의 루트에만 나타난다. 하위 디렉토리의 패키지에는 루트 패키지 이름에서 하위 디렉토리 패스가 붙은 경로로 임포트할 수 있다. 예를 들어서 `hello` 패키지 아래 `world` 패키지가 있다면, `hello/world` 패키지로 자동적으로 인식된다.

## 외부 모듈 사용하기

Go 모듈을 도입한 주된 이유는 다른 개발자들이 작성한 코드를 사용하는 경험을 개선하기 위해서이다. 방금 짠 `hello.go`에서 외부 패키지를 사용해보자.

```go
package hello

import "rsc.io/quote"

func Hello() string {
	return quote.Hello()
}
```

go 명령은 `go.mod`에 나열된 특정 의존성 모듈의 버전을 사용해 임포트를 해결한다. 그런데, 임포트한 패키지가 `go.mod` 안의 어떠한 모듈에서도 제공하지 않는 패키지라면 go 명령어는 자동적으로 그 패키지를 담고 있는 모듈을 찾아서 `go.mod`에 최신버전으로 추가한다. 예시에서는 `rsc.io/quote` 패키지를 새로 다운로드하고, 그 패키지가 사용하는 두 개의 종속성 패키지를 추가로 가져왔다.

```go
module example.com/hello

go 1.17

require rsc.io/quote v1.5.2

require (
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
	rsc.io/sampler v1.3.0 // indirect
)
```

... 공부 더 하는 중