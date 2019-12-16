![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

## Concurrent programming with Goroutines

![](../../images/goroutines1.png)

**Before starting:** I consider [**multiprogramming**](https://www.geeksforgeeks.org/difference-between-multitasking-multithreading-and-multiprocessing/) and [**multithreading**](https://en.wikipedia.org/wiki/Thread_(computing)) a ​​bad practice in designing and building applications. I believe that, as far as possible, the developer should leave this to the infrastructure that will serve the application, and not within the source code, mixed with functional code. This practice, in addition to adding [**accidental complexity**](https://www.nutshell.com/blog/accidental-complexity-software-design/) to code, also increases the cost and complexity of testing and control. quality, increasing the risks of software design. If possible, use alternatives such as [**FaaS**](http://faas.guru) and leave [**scalability**](https://en.wikipedia.org/wiki/Scalability) for infrastructure.

**Multiprogramming vs. Multiprocessing**: is a misunderstood issue among programmers and often misused. Let's look at both concepts and their differences:

- **Multiprogramming**: It is the ability to execute more than one code **concurrently**, for example by dividing the CPU into time slices and allowing each code to execute a little bit, or when a code performs a code operation. I/O, control is passed to another;
- **Multiprocessing**: It is the ability to execute more than one code **simultaneously**! On a system with multiple CPUs, we can execute multiple threads of code at a time.

Both techniques have advantages, disadvantages and risks. The biggest advantage would be to make better use of computational resources through the increased efficiency that concurrent or concurrent execution brings. But they make your code more complex, at least in most programming languages. And they carry a huge risk of **deadlocks** and **starvation**. Computer scientist E.W. **Dijkstra** has demonstrated this brilliantly with the [**Problem of Dining Philosophers**](https://en.wikipedia.org/wiki/Dining_philosophers_problem).

![](../../images/dinning-philosophers.jpeg)

- **Deadlock**: This occurs when a code that is running needs a resource that the other is using and vice versa.

## Goroutines

Let's start with our well-known [**Fibonacci sequence**](https://en.wikipedia.org/wiki/Fibonacci_number). The [**sample code**](../../code/goroutines1/fibo.go) shows how to iteratively calculate it:

```
func FibonacciLoop(n int) int {
    f := make([]int, n+1, n+2)
    if n < 2 {
        f = f[0:2]
    }
    f[0] = 0
    f[1] = 1
    for i := 2; i <= n; i++ {
        f[i] = f[i-1] + f[i-2]
    }
    return f[n]
}
```

This first version invokes the function in an iterative and blocking manner, that is, calling the function blocks the main thread (in the case of Go a **main Goroutine**):

```
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Type a term number or other character to finish: ")
		nints, _ := reader.ReadString('\n')
		nint, errInt := strconv.ParseInt(strings.TrimSpace(nints), 10, 64)
		if errInt != nil {
			break
		}
		fmt.Printf("You typed: %d and the term is %d\n", nint,FibonacciLoop(int(nint)))
	}
	
}
```

This code could be a RESTful component calling a lambda function, right? This call would be synchronous and blocking, meaning the code would have to wait for the call to **FibonacciLoop()** function to finish before continuing.

With the concept of **Goroutines** we can implement **multiprogramming** in our code in a simple and practical way. Let's look at this in the [**following example**](../../code/goroutines1/fibo2/fibo2.go):

```
func FibonacciLoop(n int) {
    f := make([]int, n+1, n+2)
    if n < 2 {
        f = f[0:2]
    }
    f[0] = 0
    f[1] = 1
    for i := 2; i <= n; i++ {
        f[i] = f[i-1] + f[i-2]
	}
	fmt.Printf("The term %d is %d\n",n,f[n])
}
...
func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Type a term number or other character to finish: ")
		nints, _ := reader.ReadString('\n')
		nint, errInt := strconv.ParseInt(strings.TrimSpace(nints), 10, 64)
		if errInt != nil {
			break
		}
		go FibonacciLoop(int(nint))
	}
}
```

As you can see, the **FibonacciLoop()** function now calculates and displays the result. A **Goroutine** is started with the **go** command before its invocation:

```
go FibonacciLoop(int(nint))
```

This creates a lightweight thread processing that lets you run the code in parallel. The call to **Goroutine** returns immediately (so we cannot have return values) and the code continues its execution. See this interesting results:

```
Type a term number or other character to finish: 5
Type a term number or other character to finish: The term 5 is 5
```

The **main()** code continued its excution and asked for one more value, and the result of the previous calculation came later!

## Communication through channels

This implementation of the **Goroutine** is ugly. It's making **I/O** within the function body! I would prefer it to return a value... Ok, so let's look at the concept of **channel**, which is a way of creating a channel between **Goroutines**. Using **channels** one **Goroutine** can write and the other will read.

We declare a **channel** specifying its data type, and initialize it with the **make** command:

```
var aChannel chan string
var otherChannel chan int
...
aChannel = make(chan string)
...
another := make(chan int)
```

We can declare and initialize in the same command!

To send or read data from a channel, we use the arrow operator "<-":

```
aChannel <- fmt.Sprintf("The term %d is %d\n",n,f[n])
...
answer := <- aChannel
```

In the first example, we send a formatted string to the channel, and in the second, we read from the channel to a variable.

Channel reading and writing is synchronous operations! When a Goroutine writes to a channel, its processing is blocked until another Goroutine reads from the channel. Respectively, when a Goroutine reads data from a channel, its processing is equally blocked until something is written on the channel.

Let's look at the [**third example**](../../code/goroutines1/fibo3/fibo3.go):

```
func FibonacciLoop(n int, aChannel chan string) {
    f := make([]int, n+1, n+2)
    if n < 2 {
        f = f[0:2]
    }
    f[0] = 0
    f[1] = 1
    for i := 2; i <= n; i++ {
        f[i] = f[i-1] + f[i-2]
	}
	aChannel <- fmt.Sprintf("The term %d is %d\n",n,f[n])
}
...
func main() {
	reader := bufio.NewReader(os.Stdin)
	chAnswer := make(chan string)
	for {
		fmt.Print("Type a term number or other character to finish: ")
		nints, _ := reader.ReadString('\n')
		nint, errInt := strconv.ParseInt(strings.TrimSpace(nints), 10, 64)
		if errInt != nil {
			break
		}
		go FibonacciLoop(int(nint),chAnswer)
		answer := <- chAnswer
		fmt.Println(answer)
	}
}
```

It has now become a synchronous call because the channel read code (```answer: = <- chAnswer```) will be blocked until the **FibonacciLoop()** function writes something to it. There are several solutions for creating non-blocking channels. We can use a **select** with a **default** option, which will be invoked without blocking the code. See this in the [**final example**](../../code/goroutines1/fibo4/fibo4.go):

```
func main() {
	reader := bufio.NewReader(os.Stdin)
	chAnswer := make(chan string)
	for {
		fmt.Print("Type a term number or other character to finish: ")
		nints, _ := reader.ReadString('\n')
		nint, errInt := strconv.ParseInt(strings.TrimSpace(nints), 10, 64)
		if errInt != nil {
			break
		}
		go FibonacciLoop(int(nint),chAnswer)
		select {
		case answer := <-chAnswer:
			fmt.Println("Got an answer: ",answer)
		default:
			fmt.Println("waiting...")
		}
	}
}
```

The execution demonstrates this behavior:

```
Type a term number or other character to finish: 6
waiting...
Type a term number or other character to finish: 8
Got an answer:  The term 6 is 8
```

It did not block execution. 


