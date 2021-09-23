---
title: Package
---

Go는 크든 작든 패키지들로 구성된다. **패키지**는 Go로 작성된 코드를 말하고, 코드 시작점에 `package` 키워드를 사용해 이름을 지정한다. 그 중 `main` 패키지는 독립적인 프로그램으로서 동작하는 소스 코드임을 알리는 패키지이고, 이외 다른 패키지는 실행 파일을 만들 수 없다. 즉, 실행을 위해서는 `main` 패키지의 메인 함수에서 호출되어야 한다. `main` 패키지는 다른 패키지가 공존하는 곳만 아니라면 어디에 짜든 상관 없다.

프로그램을 짜다 보면, 어떤 시점에서는 결국 코드를 조직화 하고 분산시키기 위해 패키지를 만들어 사용해야 하는 시점이 있다. 패키지는 연관된 소스코드를 하나의 디렉토리 아래 둠으로써, 코드를 분산시킨다. 즉, 패키지는 소스 코드들로 구성되어있다.

# 패키지 만들기
과거 Go는 패키지 임포트를 할 때, `GOPATH`에 설정된 작업 디렉토리의 `src` 디렉토리 아래나 `GOROOT`의 `src`에 패키지가 존재 했어야만 했다. 외부 모듈에 있는 패키지를 사용할 때에도 원래는 `GOPATH/src` 아래에 다운받아지는 방식이었지만, 현재는 모듈을 사용해 모듈과 유관한 패키지들을 버전과 함께 편하게 관리하고 있다.

{{< hint info >}}
`go env GOPATH` 명령어로 Go 환경 변수에 설정된 작업 디렉토리를 확인할 수 있다. 일반적으로 `~/go`로 설정되어있다.
{{< /hint >}}

{{< hint info >}}
과거엔 버전 관리 방식이 아주 신기하게, Go 패키지들이 `src` 디렉토리에 모이면 그 전체 파일들을 하나의 stable 버전으로 관리하는 방식이라고 한다. 즉, 디테일한 버전 관리는 git을 통하게 하고, `GOPATH` 아래 모인 파일들이 하나의 모듈을 이룬다고 보는 것이다.
{{< /hint >}}

예를 들어, `GOPATH/src`에 아래와 같은 패키지가 선언되어 있다.

```
bin/
pkg/
src/
  diveintogo/
    diveInto.go
```

```go
// diveintogo/diveInto.go
package diveintogo

import "fmt"

func Book() {
	fmt.Println("https://changhoi.github.io/dive-into-go")
}
```

{{< hint info >}}
패키지 이름은 되도록 소문자로 구성하는 것이 컨벤션이다.
{{< /hint >}}

현재 버전(1.17) 기준으로 이 패키지를 사용하려고 하면 `GOROOT`에 있는 것이 아니기 때문에 사용할 수 없다는 명령어를 보게 된다. 이는 패키지를 관리하는 방식이 기본적으로 모듈을 선택하도록 되어있기 때문인데, 과거 `GOPATH`를 사용하는 방식으로 바꾸면 정상 동작한다. Go는 패키지를 선택하는 방식을 관리하는 플래그를 Go 환경 변수에 담아 두었다. Go의 환경 변수에 있는 `GO111MODULE`이라는 환경 변수인데, `on` 상태인 경우 모듈 사용하는 방식으로 패키지를 가져오고, `off` 상태인 경우 `GOPATH`에서 패키지를 가져온다.

{{< hint info >}}
`GO111MODULE` 환경 변수는 임시 환경 변수로, 버전이 높아지고 모듈을 사용하는 쪽으로 완전히 정착되고 나면 사라질 예정이라고 한다. 1.17 버전에 사라질 것이라는 말도 있었는데, 일단 현재도 보인다.
{{< /hint >}}

```sh
$ GO111MODULE=off go run main.go
https://changhoi.github.io/dive-into-go
```
{{< hint info >}}
`main` 패키지는 다른 패키지와 독립적인 곳이면 아무 공간에 두어도 된다.
{{< /hint >}}

위 방식은 1회 한정으로 환경 변수 값을 설정하는 방식이다. Go 모듈이 표준화 되고 있기 때문에, 기본값인 `on`으로 두고, 모듈을 사용하도록 하자. 모듈은 별도로 정리된 페이지가 있기 때문에, 자세한 내용은 생략하고, 패키지를 사용하기 위한 기본적인 방식을 보자.

---

다음 명령어로 현재 디렉토리를 하나의 모듈로서 정의한다.

```sh
go mod init [module name]
```
그러면, 이 디렉토리를 `GOPATH/src`와 같이 취급한다. 이 루트 디렉토리를 벗어나 다른 공간에서는 이 안에 있는 패키지를 인식할 수 없다.

패키지를 구성하는 모듈로서 동작하게 하기 위해 위 명령어를 사용한 상태로 그 아래 패키지를 작성해야 한다.

{{< hint info >}}
Go 모듈도 복잡한 역사와 구현 방식이 있기 때문에, 자세한 내용은 별도의 페이지에 정리했다.
{{< /hint >}}

# 패키지 구성
패키지는 일반적으로 디렉토리 기준으로 이름을 삼는다. 물론 같아야 하는 것은 아니다. 달라도 되지만, 임포트는 디렉토리 기준으로 임포트한다. 사용은 선언한 패키지 이름으로 사용한다. 다음과 같은 구성을 가진 모듈이 있다고 해보자.

```
operation/
    math.go
go.mod
main.go
```

`main.go`는 `math.go`에서 두 수를 더하는 함수를 호출하고 있다. `math.go`는 `operation`이라는 이름의 디렉토리 아래 있지만, `mathematics`라는 이름의 패키지를 구성하고 있다.

```go
// main.go
package main

import (
	"declare/operation"
	"fmt"
)

func main() {
	fmt.Println(mathematics.Add(1, 1))
}
```

```go
// operation/math.go
package mathematics

func Add(a, b int) int {
	return a + b
}
```
위 코드는 문제 없이 돌아간다. Go가 기본적으로는 패키지를 임포트 하는 방식이 디렉토리를 기준으로 하고 있음을 알 수 있다. 이를 통해 한 가지 더 알 수 있는 점은, 하나의 패키지 디렉토리 아래, 같은 Depth에서는 여러 패키지를 가질 수 없다는 것이다. 만약 `operation/operation.go`가 있어서, 그 파일에서는 `package operation`으로 선언하면, Go에서 `declare/operation`으로 임포트 했을 때, 어떤 패키지를 가져올지 알 수 없다. 따라서, 패키지 이름은 디렉토리에 종속적이지는 않지만, 임포트는 디렉토리 기준으로 하고, 사용할 때는 선언한 패키지 이름으로 사용한다는 것을 알 수 있다.

{{< hint info>}}
디렉토리와 다른 이름으로 패키지를 만들었다면, 개발환경에 따라 앞에 alias가 붙기도 한다. `import mathematics "declare/operation"`와 같이 자동 완성 해준다. 가독성을 위해 웬만하면 같게 하는 것이 좋다.
{{< /hint >}}


# `init()` 함수

각 패키지들은 선택적으로 `init`이라는 프라이빗 함수를 만들 수 있다. 이 함수는 패키지가 초기화 될 때 자동적으로 실행되는 함수이다. 다음과 같은 특징을 가지고 있다.

- 인자 값을 갖지 않고, 리턴 값도 없다.
- 메인 함수에서도 `init` 호출이 Go 내부적으로 발생한다. 그런 경우, `main` 함수보다 앞서서 호출된다. 실제로 모든 `init` 함수들은 `main` 함수보다 먼저 호출된다.
- 여러 개의 `init` 함수가 있을 수 있고, 이는 선언된 순서로 실행된다.
- `init` 함수는 프로세스의 라이프타임 중 일 회 호출되고 이후 얼마나 임포트 되는지와 무관하게 호출되지 않는다.
- 패키지는 여러 개의 파일을 담을 수 있다. 각 소스 파일은 하나 이상의 `init` 함수들을 가질 수 있다.

`init` 함수를 사용할만한 상황은, 패키지 함수들을 실행하기 전에 필요한 시간이 걸리는 작업을 미리 처리하는 상황이나, 존재해야 하는 파일을 미리 만들어둔다든지, 프로그램이 실행하기 위해 필요한 자원들이 모두 사용 가능한지 확인하는 등의 용도가 있다.

## `main` 함수 실행을 위해 거쳐오는 길

다음과 같은 프로그램으로 확인해보자. `test` 라는 모듈을 만들어서, 그 아래 `hello/hello.go`, `variable/variable.go`, `world/world.go`를 만들어 각 디렉토리 이름으로 패키지를 만들었다. 아래와 같은 구조이다.

```
hello/
  hello.go
  greeting.go
variable/
  variable.go
world/
  world.go
go.mod
main.go
```

우선 각 패키지들은 모두 `[구분자] + INIT`이라는 문자열을 호출하는 `init` 함수가 있다. `variable.go`에서만 `init` 함수가 두 개 선언 되어있고, 각 `init` 함수는 마지막에 호출 순서에 따라 번호가 붙어있다. `hello.go`에서는 `world.go`에 있는 `World` 함수를 호출해서 `Hello` 함수를 구성하고 있다. `greeting.go`는 `main.go`에서 호출 중인 함수는 없지만, `hello` 패키지의 일부이고, 내부에 `init` 함수를 선언했다. `main.go`의 메인 함수에서는 `hello.World` 함수를 사용해 문자열을 출력한다. 그리고, 전역 변수들을 `variable.go`의 `OverVar` 함수와 `UnderVar` 함수를 이용해, 두 개를 초기화 하고 있다. 코드는 다음과 같다.

```go
// main.go
package main

import (
	"fmt"

	"test/hello"
	"test/variable"
)

var overVar = variable.OverVar()

func init() {
	fmt.Println("MAIN INIT")
}

func main() {
	fmt.Print("MAIN FUNC / ")
	fmt.Println(hello.Hello())
}

var underVar = variable.UnderVar()
```

```go
// hello/hello.go
package hello

import (
	"fmt"
	"test/world"
)

func init() {
	fmt.Println("HELLO INIT")
}

func Hello() string {
	return "Hello, " + world.World()
}
```

```go
// hello/greeting.go
package hello

import "fmt"

func init() {
	fmt.Println("GREETING INIT")
}
```

```go
// variable/variable.go
package variable

import "fmt"

func init() {
	fmt.Println("VAR INIT1")
}

func OverVar() string {
	fmt.Println("OverVar Call")
	return "OverVar"
}

func UnderVar() string {
	fmt.Println("UnderVar Call")
	return "UnderVar"
}

func init() {
	fmt.Println("VAR INIT2")
}
```

```go
// world/world.go
package world

import "fmt"

func init() {
	fmt.Println("WORLD INIT")
}

func World() string {
	return "World!"
}
```

이미지로 보자면 다음과 같이 의존 관계가 있다.

![](/dive-into-go/images/basic/packages/init-package-graph.png)

실행 결과는 다음과 같다.

```sh
$ go run main.go
WORLD INIT
GREETING INIT
HELLO INIT
VAR INIT1
VAR INIT2
OverVar Call
UnderVar Call
MAIN INIT
MAIN FUNC / Hello, World!
```

위 결과를 통해, 재귀적으로 같은 순서로 패키지를 실행하는 것을 알 수 있는데 그 순서는 다음과 같다.

1. 패키지 도입
2. 패키지에서 의존하는 패키지가 있다면, 해당 패키지 초기화 (Import된 순서대로, 다시 1번부터 재귀적으로 동작)
2. 전역 변수 초기화 (선언된 순서대로)
3. 패키지의 `init` 함수 호출 (선언된 순서대로) 

위 과정으로 결과를 나눠보면, 다음과 같다.

1. `main` 패키지 도입
   1. `hello` 패키지 도입
      1. `world` 패키지 도입
         1. `world.go`의 `init` 함수 호출
      2. `greeting.go`와 `hello.go`의 `init` 함수 호출 (알파벳 순)
   2. `variable` 패키지 도입
      1. `variable.go`의 `init` 함수 호출 (선언 순)
2. 전역 변수 초기화
3. 메인 함수 실행