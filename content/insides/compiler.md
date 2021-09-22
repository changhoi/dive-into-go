---
title: Compiler
---

이번 장은 컴파일러에 대해 아주 조금 더 깊게 공부해보자. 일반적으로 실행 파일을 만드는 명령어로 빌드까지 완성하기 때문에, 컴파일러만의 역할을 평소에 경험하기 어렵다. 컴파일러를 통해 오브젝트 파일을 만드는 것과 기타 옵션들에 대해 다룬다.

--- 

## 프로그램이 프로세스가 되기까지

## Go 컴파일러

컴파일러와 관련된 내용은 다음 


Go 코드를 컴파일 해서 오브젝트 파일을 만들기 위해서는 `go tool compile`을 사용한다. 컴파일을 하면, 마치 C언어 컴파일 한 것처럼 확장자가 `o`인 바이너리 오브젝트 파일을 내보낸다.

```sh
$ go tool compile compiled.go
```

{{< expand "결과" >}}
```sh
$ ls -al

drwxr-xr-x  4 changhoi  staff   128  9 18 13:45 .
drwxr-xr-x  3 changhoi  staff    96  9 18 13:44 ..
-rw-r--r--  1 changhoi  staff    73  9 18 13:45 compiled.go
-rw-r--r--  1 changhoi  staff  7086  9 18 13:45 compiled.o

$ file compiled.o

compiled.o: current ar archive
```
{{< /expand >}}

결과물이 아카이브 파일이 나오는데, 여러 파일을 한 파일로 묶는 바이너리 파일이다. 

... 공부하는 중