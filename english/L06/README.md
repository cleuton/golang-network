![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Base I/O

The most basic I/O operations are read and write on the console. We already know about console writing because I've used the **fmt.Println()** method from the beginning. But it does not allow formatting strings or numbers. Let's see other **fmt** methods, such as: **Printf()** and **Sprintf()**.

Reading console data basically returns **strings**. We can specify other data types, as the [**bufio package**] (https://golang.org/pkg/bufio/) specifies, but we basically use strings. Let's see how this works. In the [example program](./code/baseio.go) we read a name from the console:

```
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Your name: ")
	text, _ := reader.ReadString('\n')
    ...
```

It expects you to type a name and press <ENTER>. The resulting string will contain everything up to the <ENTER> character. The **ReadString()** method returns the read string and an error code, which only happens if the string does not contain the delimiter. In this case, we will read until we find an <ENTER>. We could use any character, such as a "$":
```
text, _ := reader.ReadString('$')
```

It will read everything until it finds a '$' character. Note that I used single quotes because it is a **rune** or a 32-bit integer representing a UTF-8 character.

And we can show what was typed. If we do with **fmt.Println()** it will look like this:

```
fmt.Printf("Hello, %v how are you today?\n", text)
...
Hello, Cleuton
 how are you today?
```

Here I used the **Printf()** method that allows us to format data. You can refer to the [**valid formats on this page**](https://golang.org/pkg/fmt/), but I will summarize here:

| Format | Description |
| --- | --- |
| %v | Any value |
| %t | true or false |
| %d | base 10 integer |
| %f | Floating point number |

Why this line break after the name? The reason is the character '\n' you enter is part of the string. We can use the [**strings library**](https://golang.org/pkg/strings/) and remove the leading and trailing spaces from the typed text:

```
newName := fmt.Sprintf("%v", strings.TrimSpace(text))
fmt.Printf("Now, %v, without the extra line\n", newName)
...
Now, Cleuton, without the extra line
```

What if we want to read numbers? For example, integers? It's simple: 

```
fmt.Print("Now, type an integer: ")
nints, _ := reader.ReadString('\n')
nint, errInt := strconv.ParseInt(strings.TrimSpace(nints), 10, 64)
if errInt != nil {
    log.Fatal("Invalid integer")
}
fmt.Printf("You typed: %d\n", nint)
```

We read a string, eliminate the **whitespaces**, and transform it into an integer 64-bit base 10 with the **ParseInt()** method. But we have to remember to test if there was an error as the user can type an invalid number. Remember **log.Fatal()**? It terminates execution with the error message.

How about floating point numbers?

```
// Floating point
fmt.Print("Now, type a float: ")
nfs, _ := reader.ReadString('\n')
nf, errf := strconv.ParseFloat(strings.TrimSpace(nfs), 64)
if errf != nil {
    log.Fatal("Invalid float number")
}
fmt.Printf("You typed: %f\n", nf)
```

**i18n**


**i18n** Go has no direct internationalization support, but there are packages that can help, such as [**golang.org/x/text**] (https://godoc.org/golang.org/x/text). You need to install it with:

```
go get -v golang.org/x/text
```



## Text files

I will show the basic I/O using text files. First, let's create a text file:

```
// Creating a text file
stringArr := []byte("Good morning.\nEat an apple.\n")
// Permission: -rw-r--r--
err := ioutil.WriteFile("/tmp/arq.txt", stringArr, 0644)
check(err)
...
cat /tmp/arq.txt
Good morning.
Eat an apple.
```

As you can see, the **ioutil.WriteFile()** method writes an array of bytes. As we use Unicode, it can record accents smoothly. The permission is the same as we would use on the **chmod** command from Linux. I created a [**check function**](https://gobyexample.com/writing-files) that uses the panic command in case of problems. Personally, I prefer **log.Fatal()**.

Now let's read that file: 

```
// Reading a text file
data, err := ioutil.ReadFile("/tmp/arq.txt")
check(err)
fmt.Printf("\nType: %T\n", data)
textContent := string(data)
fmt.Println(textContent)
...
Type: []uint8
Good morning.
Eat an apple.
```

When we read a file the data type of the variable is byte array (unsigned 8 bits). We can convert it to **string** and use normally.

## Challenge

Create a program that reads from the console the coefficient values ​​of a quadratic equation, calculates the roots, and displays it on the screen, writing a file with each calculation.