---
title: Module
---

모듈은 Go가 의존성을 관리하는 방법이다. Go는 1.11부터 모듈에 지원을 시작했다. 모듈은 파일 트리에 저장되는 **Go 패키지 모음**으로, 루트에 `go.mod` 파일이 있다. Go 1.11부터, 진행 중인 디렉토리가 `GOPATH` 변수의 `src` 아래가 아닌 경우, go 명령어는 현재 디렉토리나 또는 상위 디렉토리가 `go.mod` 파일을 가지고 있을 때 모듈을 사용하도록 했다. 그리고 1.13 이후 모든 개발 환경에서 기본 모드로 모듈을 사용했다.

{{< hint info >}}
Go 모듈은 굉장히 자주 변경점이 생기고 있다. 공식 블로그에서 업데이트 되는 걸 주기적으로 확인해보자.
{{< /hint >}}

{{< hint info >}}
이하 내용에서 `GOPATH/src`나 `GOPATH/pkg`는 Go 환경 변수에 `GOPATH`로 지정된 경로 아래 있는 `src` 디렉토리와 `pkg` 디렉토리를 의미한다.
{{< /hint >}}

## 새로운 모듈 만들기

GOPATH가 아닌 곳에 새로운 디렉토리를 만들고 패키지를 만들자.

```go
package diveIntoGo

func URL() string {
	return "https://changhoi.github.io/dive-into-go"
}
```

이 패키지를 모듈로 만들고, Github에 배포해서 다른 모듈에서도 사용할 수 있게 만들 수 있다. 우선, 모듈로 만들기 위해서 `go mod init` 명령어를 수행한다. `go mod init [저장소]`와 같이 일반적으로 사용하는데, 저장소 주소는 주소 형태여야 하는 것은 아니다. 그냥 모듈의 이름이라고 보면 된다. 그러나 다른 모듈에서 이를 다운로드 할 때 이 이름을 기준으로 다운로드 하기 때문에, 외부에 공개하기 위해서는 저장소 이름을 써야 한다.

```shell
$ go mod init github.com/changhoi/example-module
```

경과로 `go.mod` 파일을 얻을 수 있는데 처음 만들어진 `go.mod` 파일은 다음과 같이 생겼다.
```
module github.com/changhoi/example-module

go 1.17
```
`go.mod` 파일은 모듈의 루트에만 나타난다. 하위 디렉토리의 패키지에는 루트 패키지 이름에서 하위 디렉토리 패스가 붙은 경로로 임포트할 수 있다. 예를 들어서 `hello` 패키지 아래 `world` 패키지가 있다면, `hello/world` 패키지로 자동적으로 인식된다.

이렇게 만든 모듈에 태그를 붙여 커밋하고 저장소에 업로드한다.

```sh
$ git add .
$ git commit -m "v1.0.0"
$ git tag v1.0.0
$ git push origin master
$ git push origin v1.0.0
```

## 외부 모듈 사용하기

Go 모듈을 도입한 주된 이유는 다른 개발자들이 작성한 코드를 사용하는 경험을 개선하기 위해서이다. 방금 짠 모듈을 다른 모듈에서 사용해보자. 아무 곳에나 메인 함수를 만들고 `go mod init` 명령어로 모듈을 구성하고, 다음 코드를 쓴다. 모듈 이름은 임의로 설정해도 상관없다.

```go
package main

import (
	"fmt"

	"github.com/changhoi/example-module/diveIntoGo"
)

func main() {
	fmt.Println(diveIntoGo.URL())
}
```

이를 정상 동작 시키기 위해서는 `go get github.com/changhoi/example-module`을 해주어야 한다. 해당 명령어를 입력하면 `go.mod`에 변경이 생기고, `go.sum`이 생긴다. 그리고 `GOPATH/pkg/github.com`에도 `changhoi/example-module@v1.0.0`이라는 디렉토리와, 좀 전에 짰던 코드들이 다운로드 되어있는 것을 확인할 수 있다. 메인 함수를 실행해보면 정상적으로 url이 출력되는 것을 확인할 수 있다.

{{< expand "결괴" >}}
```sh
$ go run main.go
https://changhoi.github.io/dive-into-go
```
{{< /expand >}}

## 패키지를 모듈에 연결하는 과정

만약 go 명령어가 (`패키지 경로`)`package path`를 사용해서 패키지를 불러오고 있다면, go 명령어는 어떤 모듈이 해당 패키지를 제공하고 있는지를 결정해야 한다. `go` 명령어는 `빌드 목록`(`build list`)에서 패키지의 접두사 부분을 `모듈 경로`(`module path`)로 가지고 있는 모듈을 찾는다. 예를 들어서, `example.com/a/b`라는 패키지가 사용되고 있고, `example.com/a` 모듈이 빌드 리스트에 있다면 go 명령어는 `example.com/a/b` 디렉토리에 `b` 패키지가 있는지를 확인한다. 하나 이상의 `*.go` 파일이 있어야 해당 디렉토리가 패키지로서 존재한다고 판단한다. 만약 단 하나의 모듈만 해당 패키지를 제공하고 있는 경우, 그 모듈이 사용되고, 두 개 이상의 모듈이 발견되면 에러를 내보낸다.

go 명령어가 새로운 모듈을 찾아야 하는 경우, 즉 필요한 모듈이 설치되지 않은 상태인 경우, `GOPROXY` 환경 변수를 확인한다. `GOPROXY`는 콤마(`,`)로 연결된 리스트인데, 프록시 주소나 키워드 (`direct` or `off`)로 구성되어있다. 1.17 버전 기준으로 기본 값은 아래와 같다.

```sh
$ go env GOPROXY
https://proxy.golang.org,direct
```

프록시 주소는 go 명령어가 HTTP 통신을 하는 웹서버의 주소인데, 이 서버에서 `go.mod` 파일과 모듈 압축 파일을 받아온다. `Module proxy`에 대한 자세한 내용은 이 [링크](https://golang.org/ref/mod#module-proxy)에서 확인할 수 있다. `direct` 키워드는 go 명령어가 버전 컨트롤 시스템과 직접 커뮤니케이션 해야 한다는 것을 말한다. `off` 키워드는 프록시 주소들 외 커뮤니케이션을 직접 시도하지 않는다.

{{< hint info >}}
- `build list`: `go build`나 `go list`, `go test`와 같은 빌드 명령어에서 사용되는 모듈 버전들의 리스트
{{< /hint >}}

모든 `GOPROXY`의 엔트리에 따라서, go 명령어는 패키지를 제공할만한 모든 모듈 패스들에게 가장 최신 버전의 패키지를 요청한다. "패키지를 제공할만한 모든 모듈 패스"라는 뜻은, 예를 들어서 `example.com/a/b`라는 패키지를 찾고 있다면, `example.com/a`, `example.com`을 의미한다. 성공적으로 요청을 받은 모듈 패스들에 대해 go 명령어는 가장 길게 매칭된 경로의 모듈을 사용한다. 그러나 모듈은 찾았지만 찾고자 하는 패키지가 존재하지 않는다고 판단되면, 에러를 내보낸다. 만약 모듈조차 찾지 못 하면 `GOPROXY` 리스트의 다음 주소로 옮겨간다.

예를 들어서, go 명령어가 `golang.org/x/net/html`이라는 패키지를 담은 모듈을 찾고 있고, `GOPROXY`가 `https://corp.example.com,https://proxy.golang.org`로 설정되어있다면, 다음 순서로 요청을 보낸다.

- To https://corp.example.com/ (병렬 요청):
  - Request for latest version of golang.org/x/net/html
  - Request for latest version of golang.org/x/net
  - Request for latest version of golang.org/x
  - Request for latest version of golang.org
- To https://proxy.golang.org/ (모든 https://corp.example.com/ 에서의 요청이 404이거나 410라면):
  - Request for latest version of golang.org/x/net/html
  - Request for latest version of golang.org/x/net
  - Request for latest version of golang.org/x
  - Request for latest version of golang.org

이렇게 적절한 모듈을 찾게 되면, go 명령어는 `require` 키워드를 새로운 모듈 경로와 버전과 함께 메인 모듈의 `go.mod`에 추가한다. 이렇게 설정되면 미래에 같은 모듈은 항상 같은 버전의 패키지를 다운하게 된다. 만약, 연결된 패키지가 메인 모듈 안에 있는 패키지에서 임포트 되고 있지 않다면, 새로 만들어진 `require` 키워드 라인 마지막에 `//indirect` 주석이 붙게 된다.

## go.mod

하나의 모듈은 UTF-8로 인코딩 된 텍스트 파일에 정의된다. 이 파일의 이름은 `go.mod`이고, 모듈의 루트 디렉토리에 위치한다. `go.mod`는 라인 지향적인 정의를 하고 있다. 즉, 하나의 라인이 하나의 지시사항을 의미한다. 각 줄은 키워드와 인자값들로 이루어져있다.

```go.mod
module example.com/my/thing

go 1.12

require example.com/other/thing v1.0.2
require example.com/new/thing/v2 v2.3.4
exclude example.com/old/thing v1.2.3
replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
retract [v1.9.0, v1.9.5]
```
각 키워드들은 Go에서 임포트 하는 것처럼, 괄호를 사용해 인접한 블록을 하나의 블록으로 만들 수 있다.

```
require (
    example.com/new/thing/v2 v2.3.4
    example.com/old/thing v1.2.3
)
```

키워드는 다음과 같은 의미를 지닌다.

### `module`
메인 모듈의 경로를 정의한다. `go.mod` 파일은 반드시 하나의 `module` 키워드가 있어야 한다. 
```
module golang.org/x/net
```

### `go`
해당 모듈이 주어진 go 버전으로 작성되었음을 가리키는 키워드이다.
```
go 1.17
```

### `require`
모듈 의존성의 최소 요구되는 버전을 선언한다. go 명령어는 각 요구되는 모듈 버전에 대해, 다시 `go.mod` 파일을 로드해서 해당 파일이 필요한 의존성을 통합한다. 모든 요구 사항들이 로드되고 나면 [`MVS`(`Minimal Version Selelction`) 알고리즘](https://golang.org/ref/mod#minimal-version-selection)에 따라서 모듈을 선택하고 빌드 리스트를 만든다.

go 명령어는 요구되는 모듈을 직접적으로 메인 모듈에서 사용하지는 않는 경우, `//indirect` 주석을 붙인다.

```
require golang.org/x/net v1.2.3

require (
    golang.org/x/crypto v1.4.5 // indirect
    golang.org/x/text v1.6.7
)
```

### `exclude`

특정 모듈의 버전이 go 명령어에 의해 로드되는 것을 막는다. 이 키워드는 메인 모듈의 `go.mod`에만 적용되고 다른 모듈들에게는 무시된다.

```
exclude golang.org/x/net v1.2.3

exclude (
    golang.org/x/crypto v1.4.5
    golang.org/x/text v1.6.7
)
```

### `replace`

특정 버전의 모듈의 내용을 다른 모듈의 버전으로 대체하거나, 로컬 디렉토리로 대체하게 한다. go 명령어가 의존성을 해결할 때 `replacement path` 자리를 사용하게 된다.

```
replace module-path [module-version] => replacement-path [replacement-version]

require example.com/othermodule v1.2.3

replace example.com/othermodule => example.com/myfork/othermodule v1.2.3-fixed
```

### `retract`

`go.mod`에 의해 정의된 특정 모듈의 하나의 버전이나 버전의 범위를 지정한다. 이 키워드는 하나의 버전이 임시로 퍼블리싱 되어있거나 심각한 문제가 특정 버전이 배포된 다음 발견된 상황에 아주 유용하다.

```
retract version // rationale
retract [version-low,version-high] // rationale

retract v1.1.0 // Published accidentally.
retract [v1.0.0,v1.0.5] // Build broken on some platforms.
```

## Module-aware 명령어들

Module-aware라는 것은 패키지를 찾는 방법 중 하나이다. 두 가지 모드가 있는데, 하나는 Module-aware 모드이고, 하나는 GOPATH 모드이다. Module-aware 모드에서는 go 명령어들이 `go.mod`를 사용해서 의존성을 찾는다. GOPATH 모드에서는 의존성을 벤더 디렉토리(vendor directories)와 GOPATH에서 찾는다.

1.16 버전 이후로는 Module-aware이 기본 모드가 되었다. 더 낮은 버전에서는 조상 디렉토리나, 현재 디렉토리에 `go.mod`가 있어야 Module-aware 모드로 동작했다. 이 두 가지 방식은 Go 환경 변수 중에 `GO111MODULE`	 환경 변수에 의해 결정된다. 다음 세 가지 값 중 하나를 갖는다.

- `off`: `go.mod` 파일을 무시하고 GOPATH 모드를 사용한다.
- `on`이거나 세팅값이 없는 경우: Module-aware 모드로 동작한다. 심지어 `go.mod` 파일을 찾지 못했더라도 이 방법으로 동작한다. 즉, `go.mod`가 없으면 동작하지 않는 명령어들이 있다.
- `auto`: `go.mod` 파일이 있으면 Mode-aware 모드로 동작한다. 1.15 버전 이전에는 이것이 기본값이다. `go mod`의 하위 명령어들과 버전 쿼리 (version query)가 있는 `go install` 명령어의 경우는 `go.mod`가 없어도 Module-aware 모드로 동작한다.

### Build 명령어들

... 공부중

https://golang.org/ref/mod#build-commands