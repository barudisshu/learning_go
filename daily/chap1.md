# 第一章

## 特性

1. 语法简单

类型和规则与C相似，却没有C的隐晦。规则严谨，条例简单。

2. 并发模型

GO本身将所有一切纳入Goroutine作为并发，包括`main.main`入口函数。搭配channel，实现CSP模型。

3. 内存分配

Go选择了`tcmalloc`的内存分配架构。使用cache为当前执行线程提供无锁分配，多个central在不同线程间平衡内存单元复用。

4. 垃圾回收

降低STW时间，在Go的1.5版本后实现并发标记，逐步引入三色标记和写屏障。

5. 静态链路

无需安装运行环境和下载第三方库即可简单部署和发布操作。

6. 标准库

7. 工具链

## 功能预览

### 源文件

```go
package main

func main() {
  println("hello, world!")
}
```

以".go"作为文件扩展名，语句结束分号会被默认省略，支持C样式注释。入口函数main没有参数，且必须放在main包中。

用import导入标准库或第三方包。

```go
package main

import (
  "fmt"
)

func main() {
  fmt.Println("hello, world!")
}
```

> 请删除未使用的导入，否则编译器会将其当作错误。

可直接运行，或编译为可执行文件。

```shell script
go run main.go

hello, world!
```

### 变量

使用`var`定义变量，支持类型推断。基础数据类型划分清晰明确，有助于编写跨平台应用。编译器确保变量总是被初始化为零值，避免出现意外状况。

```go
package main

func main() {
	var x int32
    var s = "hello, world!"
  
    println(x, s)
}
```

在函数内部，可以省略`var`关键字，使用更简单的定义模式。

```go
package main

func main() {
    x := 100        // 注意，赋值符号不同
    println(x)
}
```

> 编译器将未使用的局部变量定义当做错误。

### 表达式

Go仅有三种流控制语句，

if

```go
package main

func main() {
    x := 100
    
    if x > 0 {
        println("x")
    } else if x < 0 {
        println("-x")
    } else {
        println("0")
    }
}
```

switch 

```go
package main

func main() {
    x := 100

    switch {
     case x > 0:
        println("x")
    case x < 0:
        println("-x")
    default:
        println("0")
    }
}
```

for

```go
package main

func main() {
 for i := 0; i < 5; i++ {
  println(i)
 }
 for i := 4; i >= 0; i-- {
  println(i)
 }
}
```

```go
package main

func main() {
 x := 0

 for x < 5 {
  println(x)    // 相当于 while (x < 5) { ... }
  x++
 }
}
```

```go
package main

func main() {
 x := 4

 for {          // 相当于while (true) { ... }
  println(x)
  x--

  if x < 0 {
   break
  }
 }
}
```

在迭代遍历时，`for...range`除元素外，还可以返回索引。

```go
package main

func main() {
 x := []int{100, 101, 102}
 
 for i, n := range x {
  println(i, ":", n)
 }
}
```

输出：

```shell script
0 : 100
1 : 101
2 : 102
```

### 函数

函数可以定义多个返回值，甚至对其命名

```go
package main

import (
 "errors"
 "fmt"
)

func div(a, b int) (int, error) {
 if b == 0 {
  return 0, errors.New("division by zero")
 }
 
 return a / b, nil
}

func main() {
 a, b := 10, 2  // 定义多个变量
 c, err := div(a, b)    // 接收多个返回值
 
 fmt.Println(c, err)
}
```

函数是第一类型，可作为参数或返回值。

```go
package main

func test(x int) func() {       // 得知返回值是个函数
 return func() {                // 匿名函数
  println(x)                    // 闭包
 }
}

func main() {
 x := 100
 f := test(x)
 f()
}
```

用`defer`定义延迟调用，无论函数是否出错，它都确保结束前被调用。

```go
package main

func test(a, b int) {
 defer println("dispose...")    // 常用来释放资源、解除锁定，或执行一些清理操作
                                // 可定义多个defer，按FILO顺序执行
 println(a / b)
}

func main() {
 test(10, 0)
}
```

输出：

```shell script
go run test.go

dispose...
panic: runtime error: integer divide by zero
```


### 数据

切片(slice)可实现类似动态数组的功能。

```go
package main

import "fmt"

func main() {
 x := make([]int, 0, 5)         // 创建容量为5的切片

 for i := 0; i < 8; i++ {       // 追加数据。当超过容量限制时，自动分配更大的存储空间
  x = append(x, i)
 }
fmt.Println(x)

}
```

输出：

```shell script
[0 1 2 3 4 5 6 7]
```

将字典(map)类型内置，可直接从运行时层面获得性能优化。

```go
package main

import "fmt"

func main() {
	m := make(map[string]int) // 创建字典类型对象
	m["a"] = 1                // 添加或设置
	x, ok := m["b"]           // 使用ok-idiom获取值，可知道key/value是否存在
	fmt.Println(x, ok)
	delete(m, "a") // 删除
}
```

> 所谓`ok-idiom`模式，是指在多返回值中用一个名为ok的布尔值来标示操作是否成功。因为很多操作默认返回零值，所以须额外说明。


结构体(struct)可匿名嵌入其他类型。

```go
package main

import "fmt"

type user struct { // 结构体类型
	name string
	age  byte
}

type manager struct {
	user
	title string // 匿名嵌入其他类型
}

func main() {
	var m manager

	m.name = "Tom" // 直接访问匿名字段的成员
	m.age = 29
	m.title = "CTO"

	fmt.Println(m)
}
```

输出：

```shell script
{{Tom 29} CTO}
```

### 方法

可以为当前包内的任意类型定义方法。

```go
package main

type X int

func (x *X) inc() { // 名称前的参数称作 receiver，作用类似python self
	*x++
}

func main() {
	var x X
	x.inc()
	println(x)
}
```

还可以直接调用匿名字段的方法，这种方式可实现与继承类似的功能。

```go
package main

import "fmt"

type user struct {
	name string
	age  byte
}

func (u user) ToString() string {
	return fmt.Sprintf("%+v", u)
}

type manager struct {
	user
	title string
}

func main() {
	var m manager
	m.name = "Tom"
	m.age = 29

	println(m.ToString()) // 调用 user.ToString()
}
```

### 接口

接口采用了duck type方式，也就是说无须在实现类型上添加显式声明。

```go
package main

import "fmt"

type user struct {
	name string
	age  byte
}

func (u user) Print() {
	fmt.Printf("%+v\n", u)
}

type Printer interface { // 接口类型
	Print()
}

func main() {
	var u user
	u.name = "Tom"
	u.age = 29

	var p Printer = u // 只要包含接口所需的全部方法，即表示实现了该接口
	p.Print()
}
```

> 另有空接口类型`interface{}`，用途类似OOP里的`system.Object`，可接收任意类型对象。

### 并发

整个运行时完全并发化设计。凡你看到的，几乎都在以goroutine方式运行。这是一种比普通协程或线程更加高效的并发设计，能轻松创建和运行成千上万的并发任务。

```go
package main

import (
	"fmt"
	"time"
)

func task(id int) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d: %d\n", id, i)
		time.Sleep(time.Second)
	}
}

func main() {
	go task(1)
	go task(2)

	time.Sleep(time.Second * 6)
}
```

输出：

```shell script
2: 0
1: 0
1: 1
2: 1
1: 2
2: 2
2: 3
1: 3
1: 4
2: 4
```

通道(channel)与goroutine搭配，实现用通信代替内存共享的CSP模型。

```go
package main

// 消费者
func consumer(data chan int, done chan bool) {
	for x := range data { // 接收数据，直到通道被关闭
		println("recv:", x)
	}

	done <- true // 通知main，消费结束
}

// 生产者
func producer(data chan int) {
	for i := 0; i < 4; i++ {
		data <- i // 发送数据
	}
	close(data) // 生产结束，关闭通道
}

func main() {
	done := make(chan bool) // 用于接收消费结束信号
	data := make(chan int)  // 数据管道

	go consumer(data, done) // 启动消费者
	go producer(data)       // 启动生产者

	<-done // 阻塞，直到消费者发回结束信号
}
```

输出：

```shell script
recv: 0
recv: 1
recv: 2
recv: 3
```
