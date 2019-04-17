# Notes on Golang
### Mal Minhas, Jan 2019

Getting started
---------------
1. Set GOPATH for workspace directory.  Must be absolute:
$ export GOPATH=D:\CODE\go
2. Go inside the workspace:
$ cd D:\CODE\go
3. Create helloWorld.go:
```
package main

import (
 "fmt"
)

func main(){
  fmt.Println("Hello World!")
}
```
4. Either compile and execute:
```
$ go build helloWorld.go
$ helloWorld.exe
Hello World!
```
5. Or simply do it in one go:
```
$ go run helloWorld.go
```

Variables
-----------
Variables in Go are declared explicitly. 
Go is a statically typed language:

```
var a int
var a = 1
message := "hello world" // shorthand for variable decl
var b, c int = 2, 3
```

Data Structures
---------------
```
var a bool = true
var b int = 1
var c string = 'hello world'
var d float32 = 1.222
var x complex128 = cmplx.Sqrt(-5 + 12i)
```

Arrays
------
Arrays are fixed in size:
```
var a [5]int
```

Slices store a sequence of elements and can be expanded at any time:
```
var b[]int
```

You can init a slice to a particular length then append:
```
numbers := make([]int,5,10)
numbers = append(numbers, 1, 2, 3, 4)
```

Subslicing:
```
slice3 := number2[1:4]
fmt.Println(slice3) // -> [2 3 4]
```

Creating an array of strings and adding to them:

```
blocks := make([]string,0)	
blocks = append(blocks,block)
```

Maps
----
```
var m map[string]int
// adding key/value
m['clearity'] = 2
m['simplicity'] = 3
// printing the values
fmt.Println(m['clearity']) // -> 2
fmt.Println(m['simplicity']) // -> 3
```

Typecasting
-----------
Converting from one type to another:
```
a := 1.1
b := int(a)
fmt.Println(b)
```

If ... else
-----------
```
if num := 9; num < 0 {
 fmt.Println(num, "is negative")
} else if num < 10 {
 fmt.Println(num, "has 1 digit")
} else {
 fmt.Println(num, "has multiple digits")
}
```

Switch
------
```
i := 2
switch i {
case 1:
 fmt.Println("one")
case 2:
 fmt.Println("two")
default:
 fmt.Println("none")
}
```

Loops
-----
Go has a single keyword for the loop. 
A single for loop command help achieve different kinds of loops.  
This is equivalent to a while loop:
```
i := 0
sum := 0
for i < 10 {
 sum += 1
  i++
}
fmt.Println(sum)
```

This is a normal for loop:
```
sum := 0
for i := 0; i < 10; i++ {
  sum += i
}
fmt.Println(sum)
```

This is an infinite loop:
```
for {
}
```

Pointers
--------
tbd

Interfaces
----------
tbd

Functions
---------
The main function defined in the main package is 
the entry point for a go program to execute. 
More functions can be defined and used:

```
func add(a int, b int) int {
  c := a + b
  return c
}
func main() {
  fmt.Println(add(2, 1))
}

Note you can predefine the return value thus:
func add(a int, b int) (c int) {
  c = a + b
  return
}
```

or thus:
```
func add(a int, b int) (int, string) {
  c := a + b
  return c, "successfully added"
}
```

Note if you get this kind of error:
```
processBlock(pblocks, &statementBlock) used as value
```
it means you haven't set a return value for the function.

Structs, Methods and Interfaces
-------------------------------
* A struct is a typed, collection of different fields. 
A struct is used to group data together. 
* Methods are a special type of function with a receiver. 
A receiver can be both a value or a pointer. 
* Go interfaces are a collection of methods. 
Interfaces help group together the properties of a type. 

Packages
--------
We write all code in Go in a package. 
There are two types of packages. An executable package and an utility package. 
A executable package is your main application since you will be running it. 
An utility package is not self executable.  Instead it enhances functionality 
of an executable package by providing utility functions and other important assets.
A utility package = a subdirectory in src.
The main package is the entry point for the program execution. 
There are lots of built-in packages in Go. 
The most commonly-used one is the fmt package.
Here is an example package:
```
package person
func Description(name string) string {
  return "The person name is: " + name
}
func secretName(name string) string {
  return "Do not share"
}
```

And here is how to use it:
```
package main
import(
  "custom_package/person"
  "fmt"
)
func main(){ 
  p := person.Description("Milap")
  fmt.Println(p)
}
// => The person name is: Milap
```

Like export syntax in JavaScript, Go exports a variable if a variable name 
starts with Uppercase. All other variables not starting with an uppercase 
letter is private to the package.
For an executable package, a file with main function is the entry for execution.
Go first searches for package directory inside GOROOT/src directory and if 
it doesn’t find the package, then it looks in GOPATH/src.
VSCode compiles the package when you save it if you have Go plugin installed.

Installing a package locally before using:
go get github.com/yosssi/gohtml
go get github.com/basgys/goxml2json
go get github.com/docopt/docopt-go

Error Handling
--------------
```
package main

import (
  "fmt"
  "net/http"
)

func main(){
  resp, err := http.Get("http://example.com/")
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(resp)
}
```

Defer
-----
A defer statement makes the program execute the line 
at the end of the execution of the program

Concurrency
-----------
Go routines are functions which can run in parallel or 
concurrently with another function. 
Creating a Go routine is very simple. 
Simply by adding a keyword Go in front of a function, 
we can make it execute in parallel. 
Go routines are very lightweight, 
so we can create thousands of them.

```
package main
import (
  "fmt"
  "time"
)
func main() {
  go c()
  fmt.Println("I am main")
  time.Sleep(time.Second * 2)
}
func c() {
  time.Sleep(time.Second * 2)
  fmt.Println("I am concurrent")
}
//=> I am main
//=> I am concurrent
```
