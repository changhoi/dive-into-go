`time` 패키지는 시간을 측정하거나 보여주기 위해 사용한다. 깊은 구현 방식까지 탐구하고자 했으나, 아주 복잡한 `wall clock`, `monotonic clock`에 대해 깔끔하게 정리할 수 있는 재주가 없어서 우선 실용적인 레벨에서 사용 방법을 정리해봤다.

---

## `type Time`

`Time`은 나노초 정밀도로 한 순간을 표현한다.

```go
type Time struct {
	wall uint64
	ext  int64
	loc *Location
}
```

`wall`값과 `ext` 값은 위에서 잠깐 언급한 `wall clock`의 초 단위, 나노초 단위로 변환하는 데 사용되고 선택적으로 `monotonic clock`을 나노초 단위에서 읽기 위해 사용된다. `loc` 값은 지역을 특정하기 위해 사용한다. `Location`이라고 하는 타입의 포인터 값을 갖는다. 만약 이 값이 `nil`이라면, `UTC`임을 의미한다.

> UTC는 협정 세계시를 의미하고 뉴질랜드 웰링턴을 **종점**으로 하는 기준시이다. 영국을 기준시로 하는 GMT와 유사하지만, 초의 소숫점에서 차이가 있어서 기술적 영역에서는 UTC를 사용한다고 한다.

> 간단히 말해서, `wall clock`은 시각과 관련된 정보이다. `monotonic clock`은 시간과 관련된 정보이다. 

## `time.Now`

`time.Now` 함수로 현재 로컬 시간을 만들수 있다.

```go
var t time.Time = time.Now()
```
