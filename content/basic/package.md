---
title: Package
---

Go는 크든 작든 패키지들로 구성된다. **패키지**는 Go로 작성된 코드를 말하고, 코드 시작점에 `package` 키워드를 사용해 이름을 지정한다. 그 중 `main` 패키지는 독립적인 프로그램으로서 동작하는 소스 코드임을 알리는 패키지이고, 이외 다른 패키지는 실행 파일을 만들 수 없다. 즉, 실행을 위해서는 `main` 패키지의 메인 함수에서 호출되어야 한다.


## 패키지 설치

Go는 풍부한 라이브러리를 가지고 있다. 그러나 기본 라이브러리 외 기능을 위해 패키지를 외부에서 설치해야 하는 상황도 있다.


## Go 패키지 만들기

일반적으로 패키지는 패키지 이름으로 된 디렉토리 아래 둔다. 다만, 예외적으로 `main` 패키지만 어디에서든 작성될 수 있다. `mypackage`라는 패키지를 구성해보자. `mypackage`는 같은 이름의 디렉토리 아래 있다.

```go
package mypackage

import "strings"

func GetName() Name {
	return Name("dive into go")
}

const VERSION = "1.0"

type Name string

func (n Name) String() string {
	return strings.ToUpper(string(n))
}
```

이제 이 패키지를 사용하려면, `main` 패키지에서 사용해야한다.

```go
package main

import (
	"fmt"
	"mypackage"
)

func main() {
	fmt.Print(mypackage.GetName())
}
```
{{< expand "결과" >}}
```sh
$ go run mypackage/main.go

mypackage/main.go:5:2: package mypackage is not in GOROOT (/usr/local/go/src/mypackage)
```
{{< /expand >}}

이제 실행하려 하면, `mypackage`라는 패키지를 찾을 수 없다고 에러를 보여줄 것이다. 이는 패키지가 