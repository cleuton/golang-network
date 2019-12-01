![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Code Modularity

To modularize a source code is to divide it into different abstraction layers. We have a higher level layer and one or more lower level layers, containing functions or classes.

## Functions

A **function** may or may not take arguments and may or may not return a value. Non-returning value functions are known as **side-effect functions** because they change the program's state.
```
func delta(a float64, b float64, c float64) float64 {
	return math.Pow(float64(b), 2) - 4*a*c
}
```

The **delta** function takes 3 float arguments (float64) and returns a float value (float64). It can be invoked like this:

```
vDelta := delta(coefA, coefB, coefC)
```

Sample [**function.go**](./codigo/function.go) demonstrates the use of functions: 

```
package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func delta(a float64, b float64, c float64) float64 {
	return math.Pow(float64(b), 2) - 4*a*c
}

func readFromArgs(a []string, p int) float64 {
	f, err := strconv.ParseFloat(a[p], 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func main() {
	clArguments := os.Args[1:] // We got command line arguments
	coefA := readFromArgs(clArguments, 0)
	coefB := readFromArgs(clArguments, 1)
	coefC := readFromArgs(clArguments, 2)
	vDelta := delta(coefA, coefB, coefC)
	fmt.Println(vDelta)
}
```

**Reading command line arguments**

We use the **os** package to read command line arguments. The **readFromArgs()** function does it:
```
clArguments := os.Args[1:]
...
f, err := strconv.ParseFloat(a[p], 64)
```

The **strconv** object converts data to strings and vice versa. The **ParseFloat** method converts strings to floats, and we can specify the precision, in this case 64 bits. We get from the **os.Args[]** vector a **slice** from the second position (os.Args [1:]), to avoid taking the program name, which is always the first argument. Going forward, we will have a slice starting at zero, with each argument passed after the program name.

Another interesting thing to note is the exponentiation, provided by the **Pow()** method of the **math** object. This method requires the variables to be **float64** and that is why I am using this data type.

**Error Handling**

**Go** do not use **SEH** (Structured Exception Handling) as in C ++, Java or Python (try). We simply test if there were any errors in the function. And functions can return multiple values, for example:

```
	
func f() (int, int) {
    return 1, 2
}
...
a, b := f()
```

Some functions also return errors. Here's how I test the string to float conversion error:

```
f, err := strconv.ParseFloat(a[p], 64)
if err != nil {
    log.Fatal(err)
}
```

Method **ParseFloat()** has the signature: 

```
func ParseFloat(s string, bitSize int) (float64, error)
```

Therefore, it may return an error. I can test if the error is **nil**, (ie null). However, if there is an error, I use the **log** method **Fatal()** which displays the message and paralyzes execution.

To execute the sample code:

```
go run function.go 1 -5 6
```

However, to run this code, you will need to have **Go** installed, as [**MyCompiler**] (https://www.mycompiler.io/) does not allow us to pass command line arguments.


## Installing Go

[**Installing Go**] (https://golang.org/doc/install) is very simple. Do it according to your operating system. Here, I am using **Ubuntu**:

Download the [**package**] (https://golang.org/dl/) and extract it into /usr/local, create a tree in /usr/local/go. For example:

```
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
```
Choose the appropriate package for your installation. For example, if you are installing Go 1.2.1, for Linux x86 64-bit, the package is **go1.2.1.linux-amd64.tar.gz**.

(These commands must be executed with **sudo**)

Add /usr/local/go/bin to the **PATH** environment variable.

If you use **Windows** or **MacOS** [see instructions] (https://golang.org/doc/install).

## Go command

Open a terminal and execute **Go**: 

```
go
...
$ go
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildmode   build modes
        c           calling between Go and C
        cache       build and test caching
        environment environment variables
        filetype    file types
        go.mod      the go.mod file
        gopath      GOPATH environment variable
        gopath-get  legacy GOPATH go get
        goproxy     module proxy protocol
        importpath  import path syntax
        modules     modules, module versions, and more
        module-get  module-aware go get
        module-auth module authentication using go.sum
        module-private module configuration for non-public modules
        packages    package lists and patterns
        testflag    testing flags
        testfunc    testing functions

Use "go help <topic>" for more information about that topic.
```

The two basic options are: ```build``` and ```run```. To run the example code, [**function.go**] (./code/function.go), open the folder in a terminal and type:
```
go run function.go 1 -5 6
...
1
```

## Challenge

Complete the program by creating a function that calculates the roots of the quadratic equation. Remember that the delta value is crucial to this!

Tip: Use math.Inf to return infinity (if you know how to solve quadratic equations, you know this could happen)