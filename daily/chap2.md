# 第二章 类型

## 变量

在数学概念中，变量(variable)表示没有固定值且可变的数。从计算机系统实现角度看，变量是一段或多段用来存储数据的内存。

作为静态类型语言，Go变量总是有固定的数据类型，类型决定了变量内存的长度和存储格式。我们只能修改变的值，无法修改类型。

> 通过类型转换或指针操作，我们可用不同方式修改变量值。但这并不意味着改变了变量类型。

因为内存分配发生在运行期，所以在编码阶段我们用一个易于阅读的名字来表示这段内存。实际上，编译后的机器码从不使用变量名，而是直接通过内存地址访问目标数据。保存在符号表中的变量名等信息可被删除，或用于输出更详细的错误信息。

### 定义

关键字`var`用于定义变量，和C不同，类型被放在变量名后。另外，运行时内存分配操作回确保变量自动初始化为二进制零值(zero value)，避免出现不可预测行为。如显示提供初始化值，可省略变量类型，由编译器推断。

```go
var x int			// 自动初始化为0
var y = false // 自动推断为bool类型
```

可以一次定义多个变量，包括用不同初始值定义不同类型。

```go
var x, y int 	// 相同类型的多个变量
var a, s = 100, "abc" // 不同类型初始值
```

依照惯例，建议以组方式整理多行变量定义。

```go
var (
	x, y int
  a, s =100, "abc"
)
```

### 简短模式

除`var`关键字外，还可以使用更加简短的变量定义和初始化语法。

```go
func main() {
  x := 100
  a, s := 1, "abc"
}
```

只是注意，简短模式(short variable delaration)有些限制：

- 定义变量，同时显式初始化。
- 不能提供数据类型。
- 只能用在函数内部。

对于粗心的新手，这可能会造成意外错误。比如原本打算修改全局变量，结果变成重新定义同名局部变量。

```go
var x = 100

func main() {
  println(&x, x)					// 全局变量
  
  x := "abc"							// 重新定义和初始化同名局部变量
  println(&x, x）
}
```

输出：

```shell
0xae020						100			// 对比内存地址，为两个不同的变量
0xc8200041f38			abc
```

> 简短模式在定义函数多返回值，以及`if/for/switch`等语句中定义局部变量非常方便。

简短模式并不是重新定义变量，也可能是部分退化的赋值操作。

```go
func main() {
  x := 100
  println(&x)
  
  x, y := 200, "abc" 			// 注意，x 退化为赋值操作，仅有y是变量定义
  
  println(&x, y)
  println(y)
}
```

输出：

```shell
0xc000046748 100
0xc000046748 200				// 对比内存地址，可以确认x属于同一变量
abc
```

退化操作的前提条件是：最少一个变量被定义，且必须是同一作用域。

```go
func main() {
  x := 100
  println(&x, x)
  
  x := 200							// 错误，no new varables on left side of :=
  println(&x, x)
}
```

```go
func main() {
  x := 100
  println(&x, x)
  
  {
    x, y := 200, 300		// 不同作用域，全部是新变量定义
    println(&x, x, y)
  }
}
```

输出：

```shell
0xc000046748 100
0xc000046740 200 300
```

在处理函数错误返回值时，退化赋值允许我们重复使用`err`变量，这是相当有益的。

```go
package main

import (
	"log"
	"os"
)

func main() {

	f, err := os.Open("/dev/random")
	...
	buf := make([]byte,1024)
	n, err := f.Read(buf)			// err 退化赋值，n新定义
	...
}
```

### 多变量赋值

在进行多变量赋值操作时，首先计算出所有右值，然后在依次完成赋值操作。

```go
func main() {

	x, y := 1, 2
	x, y = y + 3, x + 2			// 先计算出右值y + 3, x + 2，然后再对x, y变量赋值
	
	println(x, y)
}
```

输出：

```shell
$ go build && ./test                  
5 3
$ go tool objdump -s "main\.main" test
TEXT main.main(SB) test.go
  MOVQ GS:0x30, CX                    // 先使用AX, CX寄存器完成表达式x+2、y+3操作    
  CMPQ 0x10(CX), SP                       
  JBE 0x10515b5                           
  SUBQ $0x10, SP                          
  MOVQ BP, 0x8(SP)                    // 将计算结果依次写入x,y    
  LEAQ 0x8(SP), BP                        
  CALL runtime.printlock(SB)    			// 依次println(x),println(y)          
  MOVQ $0x5, 0(SP)                        
  CALL runtime.printint(SB)               
  CALL runtime.printsp(SB)                
  MOVQ $0x3, 0(SP)                        
  CALL runtime.printint(SB)               
  CALL runtime.printnl(SB)                
  CALL runtime.printunlock(SB)            
  MOVQ 0x8(SP), BP                        
  ADDQ $0x10, SP                          
  RET                                     
  CALL runtime.morestack_noctxt(SB)       
  JMP main.main(SB) 
```

赋值操作，必须确保左右值类型相同。

### 未使用错误

编译器将未使用局部变量当做错误。不要觉得麻烦，这有助于培养良好的编码习惯。

```go
var x int			// 全局变量没问题
func main() {
  y := 10
}
```

输出：

```shell
$ go build
./test.go: y declared and not used
```

### 命名

对变量、常量、函数、自定义类型进行命名，通常优先选用有实际含义、易于阅读和理解的字母或单次组合。

#### 命名建议

-  以字母或下划线开始，由多个字母、数字和下划线组合而成。
- 区分大小写。
- 使用驼峰(camel case)拼写格式。
- 局部变量优先使用短名。
- 不要使用保留关键字。
- 不建议使用与预定义常量、类型、内置函数相同的名字。
- 专有名词通常会全部大写，例如escapeHTML。

> 尽管Go支持用汉字等Unicode字符命名，但从编程习惯上来说，这并不是好选择。

```go
func main() {

	var c int							// c代替count
	for i := 0; i < 10; i++ {			// i代替index
		c++
	}

	println(c)
}
```

符号名字首字母大小写决定了其作用域。首字母大写的为导出成员，可被包外引用，而小写则仅能在包内使用。相关细节，可参考后续章节。

### 空标识符

和Python类似，Go也有个名为“_”的特殊成员(blank identifier)。通常作为忽略占位符使用，可作表达式左值，无法读取内容。

```go
import "strconv"

func main() {

	x, _ := strconv.Atoi("12")		// 忽略Atoi的err返回值
	println(x)
}
```

空标识符可用来临时规避编译器对未使用变量和导入包的错误检查。但请注意，它是预置成员，不能重新定义。

## 常量

常量表示编译期确定的字符、字符串、数字或布尔值。

```go
const x, y int = 123, 0x22
const s = "hello, world!"
const c = '我'

const (
    i, f = 1, 0.123
    b    = false
)
```

可以在代码块中定义常量，不曾使用的常量不会因为编译错误。

```go
func main() {
    const x = 123
    println(x)
    
    const y = 1.23

    {
        const x = "abc"
        println(x)
    }   
}
```

如果显示指定类型，必须确保常量左右值类型一致，需要时可做显式转换。

```go
const (
    x, y int  = 99, -999
    b    byte = byte(x)
    n         = uint8(y) // 错误， 超出uint长度
)
```

常量值也可以是某些编译器能计算出结果的表达式，如unsalfe.Sizeof、len、cap等。

```go
const (
    ptrSize = unsafe.Sizeof(uintptr(0))
    strSize = len("hello, world!")
)
```

在常量组中如不指定类型和初始化值，则与上一行非空常量右值（表达式文本）相同。

```go
import "fmt"

func main() {
	const (
		x uint16 = 120
		y        // 与上一行x类型、右值相同
		s = "abc"
		z 		// 与s类型、右值相同

	)

	fmt.Printf("%T, %v\n", y, y) // 输出类型和值
	fmt.Printf("%T, %v\n", z, z)
}
```

### 枚举
Go并没有明确意义上的enum定义，不过可借助iota标识符实现一组自增常量值来实现枚举类型。

```go
const (
    x = iota    // 0
    y           // 1
    z           // 2
)
```

```go
const (
    _  = iota // 0
    KB = 1 << (10 * iota)
    MB
    GB
)
```
自增作用范围为常量组。可在多常量定义中使用多个iota。

```go
const (
    _, _ = iota, iota * 10
    a, b
    c, d
)
```

如果中断了iota自增，必须显式恢复。且后续自增值按行序递增，

```go
const (
    a = iota		// 0
    b				// 1
    c = 100			// 100
    d				// 100
    e = iota		// 4
    f				// 5
)
```

自增默认数据类型为int，可显式指定类型。

```go
const (
    a         = iota // int
    b float32 = iota // float32
    c         = iota // float32
)
```

在实际编码中，建议用自定义类型实现用途明确的枚举类型。但这并不能将取值范围限定在预定义的枚举值内。

```go
type color byte

const (
	black color = iota
	red
	blue
)

func test(c color) {
	println(c)
}

func main() {
	test(red)
	test(100)		// 100 并未超出color/byte类型取值范围
	
	x := 2
	test(x)			// 错误： x的类型为int
}
```

### 展开

常量除"只读"外，和变量究竟有什么不同？

```go
var x = 0x100

const y = 0x200

func main() {
	println(&x, x)
	println(&y, y) // 错误：cannot take the address of y
}
```

不同于变量在运行期分配存储内存，常量通常会被编译器在预处理阶段直接展开，作为指令数据使用。

```go
const y = 0x200

func main() {
    println(y)
}
```

数字常量不会分配存储空间，无须像变量那样通过内存寻址来取值，因此无须获取地址。


## 基本类型














