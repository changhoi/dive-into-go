---
title: context
---

`context` 패키지는 Go 개발을 하다 보면 심심치 않게 만날 수 있다. 패키지의 공식문서에 따르면, 이 패키지의 주된 목표는 데드라인, 취소 신호 및 기타 요청 범위 값을 API 경계와 프로세스들 사이에 전달하는 컨텍스트 유형을 만들기 위해 사용한다.

---

## `Context` 타입

`Context` 타입은 다섯 가지의 메소드를 가지고 있는 인터페이스이다. `Deadline`, `Err`, `Value`, `Done`이 있다. 

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}
```
### `Deadline`
컨텍스트를 대신해 수행된 작업을 취소해야 하는 시간을 반환한다. 데드라인은 설정되지 않으면 `ok`에 `false`담아 값을 반환한다.

### `Done`
이 함수는 `Context`가 취소되었거나, 타임아웃이 발생하면 닫힌 수신용 채널을 리턴한다. 만약 이 `Context`가 취소될 수가 없는 것이라면, `nil`을 반환한다.

### `Err`
`Done`이 아직 닫히지 않았다면, 이 함수는 `nil`을 반환한다. 만약 `Done`이 닫혀있다면, `Err`은 이유가 담긴 에러를 리턴한다. 만약 캔슬이 호출되었다면, `Canceled` 에러를 리턴하고, 데드라인이 지났다면, `DeadlineExceeded` 에러를 반환한다.

### `Value`
해당 `Context`와 연관된 `key`에 대한 값을 리턴하거나, 어떠한 값도 이 키와 연관이 없는 경우 `nil`을 리턴한다.

`Context` 타입을 사용해야 한다면, 직접 구현하기 보다는 `WithCancel`, `WithDeadline`, `WithTimeout` 같은 함수들을 사용해 만들 수 있다.

---

`WithCancel`, `WithDeadline`, `WithTimeout` 함수들은 `Context`를 받아서 각 이름에 맞게 설정된 새로운 자식 `Context`와 `CancelFunc` 형태의 함수를 내보낸다. `CallFunc`은 호출되면 해당 컨텍스트와 그 자식들을 취소 시키고, 컨텍스트의 부모가 해당 컨텍스트를 참조하고 있던 레퍼런스를 삭제한다. 그리고 연관된 타이머들을 모두 정지시킨다. `CallFunc`를 호출하는 것에 실패하면 부모가 취소되기 전까지 메모리 릭이 발생하는 것과 같다.

## 빈 `Context` 만들기
비어있는 `Context`라는 말은, 절대 취소되지도 않고, 값고 없고, 데드라인도 없는 `Context`라는 뜻이다. 이런 빈 `Context`는 다음 두 가지 함수로 만들 수 있다.

### `Background()`
빈 `Context`를 리턴한다. 일반적으로 메인 함수, 초기화, 테스트 과정 그리고 들어오는 요청들에 대한 루트 `Context`로 자주 사용된다.

### `TODO()`
어떤 `Context`를 사용할지 모르겠거나, 명확하지 않다면 이 함수를 사용해 빈 `Context`를 만든다.

## WithCancel

부모가 되는 `Context`를 인자로 받아 새로운 `Done` 채널과 함께 복사 한 값을 리턴한다. 반환된 `Context`의 `Done` 채널은 반환된 취소 함수가 호출되거나 부모의 `Done` 채널이 막히면 같이 막히게 된다.

```go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	ctx := context.Background()
	parent, parentCancel := context.WithCancel(ctx)
	child, childCancel := context.WithCancel(parent)

	if len(os.Args) != 3 {
		fmt.Printf("usage: %s <parent delay> <child delay>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	parentDelay, _ := strconv.Atoi(os.Args[1])
	childDelay, _ := strconv.Atoi(os.Args[2])

	go func() {
		time.Sleep(time.Duration(parentDelay) * time.Second)
		parentCancel()
	}()

	go func() {
		time.Sleep(time.Duration(childDelay) * time.Second)
		childCancel()
	}()

	select {
	case <-child.Done():
		fmt.Println("Child closed!:", child.Err())
		fmt.Println("Parent not closed!:", parent.Err())
	case <-parent.Done():
		fmt.Println("Parent closed!:", parent.Err())
		fmt.Println("Child closed!:", child.Err())
	}
}
```
{{< expand "결과" >}}
```sh
$ go run withCancel.go 1 2
Parent closed!: context canceled
Child closed!: context canceled

$ go run withCancel.go 2 1
Child closed!: context canceled
Parent not closed!: <nil>
```
{{< /expand >}}

위 예시에서처럼 부모가 닫히면 자동으로 자식의 `Done` 채널이 닫히면서 에러값을 갖게 되고, 반대는 성립하지 않는다.

---

취소하는 동작은 고루틴이 릭되는 현상을 막아줄 수 있다. [아래는 Go 패키지에서 제공되는 예시이다.](https://pkg.go.dev/context#example-WithCancel)

```go
import (
	"context"
	"fmt"
)

func main() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // 고루틴이 릭되는 것을 막기 위해 리턴
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 정수를 하나씩 소비하는 것을 끝내면 고루틴 종료

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
```
{{< expand "결과" >}}
```sh
1
2
3
4
5
```
{{< /expand >}}

## WithDeadline
사실 Cancel에 대해 이해 했다면, 그 뒤는 비슷하다. 
이 함수는 부모 `Context`에서 데드라인을 조정해서 내보낸다. 인자 값으로 들어오는 `d` 이후를 데드라인으로 설정한다. 만약 부모 `Context`의 데드라인이 설정되어있고, `d`보다 이르다면, 부모에게 맞춘다 (어짜피 부모가 취소 되면 자식되 취소 됨). 반환된 `Context`의 `Done` 채널은 데드라인이 지나거나, `cancel` 함수를 호출하면 닫히게 된다.

```go
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage %s <deadline>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	t, _ := strconv.Atoi(os.Args[1])
	deadline := time.Now().Add(time.Duration(t) * time.Second)

	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, deadline)
	go func() {
		select {
		case <-time.After(time.Duration(3) * time.Second):
			fmt.Println("TOO LONG DEADLINE!")
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("closed!:", ctx.Err())
	}
}
```
{{< expand "결과" >}}
```sh
$ go run withDeadline.go 2
closed!: context deadline exceeded

$ go run withDeadline.go 4
TOO LONG DEADLINE!
closed!: context canceled
```
{{< /expand >}}

데드라인을 지나서 `Done`의 채널이 닫히면, `DeadlineExceeded` 에러를 보내고, `cancel` 함수를 사용해서 닫으면 `Canceled` 에러를 보낸다.

## WithTimeout
이 함수는 데드라인 함수를 쉽게 쓰게 해준다. 인자로 들어온 `Duration`을 가지고 다음과 같이 리턴한다.
```go
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
```

## WithValue
이 함수는 다른 함수들과 마찬가지로, 부모 `Context`를 사용하지만, 키와 값을 함께 넣어서 저장한다. `Context`에 특정 값을 저장해둘 수 있다. 이 값은 추가 인자값을 다른 함수에 전달하는 용도로 쓰지말라고 공식 문서에 설명되어있다.

제공된 키는 반드시 비교가능해야 하고, 문자열이나 내장된 유형이면 안된다. 컨텍스트를 사용하는 패키지 사이에 충돌을 방지하기 위해서이다. `WithValue`를 사용하려면, 키 유형을 직접 정의해야한다.

아래는 [Go 패키지에 있는 예시](https://pkg.go.dev/context#example-WithValue)이다.

```go
package main

import (
	"context"
	"fmt"
)

func main() {
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))
}
```

{{< expand "결과" >}}
```sh
found value: Go
key not found: color
```
{{< /expand >}}

---

## 더 실용적인 예시들

Context는 언급한 대로 API의 사이 경계 부분이라든지, 고루틴의 릭되는 것을 막기 위해 Concurrency 패턴에서 자주 등장한다. 아래 블로그 글들을 보면, `context` 패키지의 실용적 예시들을 확인할 수 있다.


- <https://gobyexample.com/context>
- <https://go.dev/blog/context>