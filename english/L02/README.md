![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Datatypes, assignment, blocks

Go has a rich set of data types, including:

**Integers**

| Type | Description |
|------------|--------|
| uint8 | 8 unsigned bits |
| uint16 | 16 unsigned bits |
| uint32 | 32 unsigned bits |
| uint64 | 64 unsigned bits |
| int8   | 8 signed bits |
| int16  | 16 signed bits |
| int32  | 32 signed bits |
| int64  | 64 signed bits |

Let's look at an example of declared and inferred type variable declaration:

```
package main

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x uint8 = 10
	// An error will happen if you uncomment the following line:
	// y = 11
	y := 11
	fmt.Println(reflect.TypeOf(x), x, reflect.TypeOf(y), y)
...
uint8 10 int 11
```

In this code we see the declaration of two variables: "x", 8-bit unsigned integer, and "y", integer (**int**). The type **int** is size dependent on the architecture of the computer where you are compiling the code.

More interesting is that we see two ways we can declare variables: With explicit and implicit type declaration. When it is necessary to specify the type, we use the command **var**:

```
var <name 1>, <name 2> <data type> = <initial value 1>, <initial value 2>
```

In the other statement, the data type is inferred by the assigned value. In this case, we need to add a colon (":") character before the assignment sign (=) indicating that we are declaring and initializing a variable. Note that initialization is not required when using **var**. Here is the implicit statement:

```
y := 11
```

If you try to declare the variable as you do in **Python**, ie simply by assigning the value, the compiler will show you an error indicating that the variable was not declared.

Another interesting thing was how I showed the data type of each variable: the "reflect" package, which, among other things, has the "TypeOf()" method. To use it, I had to declare it in the list of **import**.


**Floating point numbers**

| Type | Description |
|------------|--------|
| float32 | 32 bits |
| float64 | 64 bits |

See this sample code [**datatypes.go**](./code/datatypes.go):

```
var z float32 = 22 / 7.0
pi := 22 / 7.0
fmt.Println(reflect.TypeOf(z), z, reflect.TypeOf(pi), pi)
...
float32 3.142857 float64 3.142857142857143
```

I calculated **PI** dividing 22 by 7. Note that I used decimal point on digit 7. This was to turn it into a **float** number, otherwise it would give an integer result.

Note that the first variable, **z**, was declared as **float32**, and the second, **pi**, was declared as implicit type, and the compiler assumed **float64**. Notice how the accuracy of the calculated value was higher. The larger the floating type, the greater the accuracy.

**Text**

**Go** has the type **string** to declare strings:

```
var name string = "John"
surname := " Doe"
fmt.Println(reflect.TypeOf(name), name, reflect.TypeOf(surname), surname)
fmt.Println(name + surname)

message :=
	`This is a multiline string, 
because there are 2 lines of text. `
fmt.Println(message)

other := "this is also a\nmultiline message."
fmt.Println(other)
... 
string John string  Doe
John Doe
This is a multiline string, 
	because there are 2 lines of text. 
this is also a
multiline message.
```

We can declare **string** literals with double quotation marks or backticks ("`"). When delimited with backticks, we can expand the text over multiple lines. When using double quotes, we can include "\n" for newline or "\t" for tab. And we can concatenate strings as we do in **Java**.

**Boolean type**

**Go** has the **bool** type to declare logical variables. And we have the boolean operators: 

- And: && 
- Or: || 
- Not: !
- Exclusive Or: Not exists. 

```
var (
    p bool = true
    q bool = false
)
r := p && !q
fmt.Println(r)
fmt.Println("p XOR q", p != q)
fmt.Println("p XOR r", p != r)
...
true
p XOR q true
p XOR r false
```

I demonstrated a **var** block with multiple variables being declared. I also simulated **XOR** with a "trick" (p! = Q).

**Constants**

constants can be declared with statement **const**: 

```
const tipec = "*"
fmt.Println(tipec)
...
*
```

**Arrays and slices**

**Go** has a special way of dealing with multivalued variables: the **array** data type. We can also deal with array **slices**. The first array element has index zero, and the last has **length - 1**.

We declare arrays by stating the size, enclosed in braces, and the type of values. And we can initialize them:
```
var grades = [5]float32{5.5, 7.5, 8.0, 5.3, 9.2}
```

We can declare **slices** indicating the start position and end positions. End position is exclusive:

```
firstOnes := grades[0:2]
```

We take from the first element up to the third, but not including it (exclusive).

And we can declare a slice omitting one of the informations:

```
lastOnes := grades[3:]
```

In this example, we take from fourth position up to the end.

And we can iterate over arrays using the **for** command with the **range** option. At each iteration, we will receive two new variables: The index and the value of that position:

```
for ix, val := range grades {
	fmt.Println("Index", ix, "value", val)
}
```

Here's the full source code: 

```
var grades = [5]float32{5.5, 7.5, 8.0, 5.3, 9.2}
fmt.Println(grades)
for ix, val := range grades {
	fmt.Println("Index", ix, "value", val)
}
grades[1] = 8.0
fmt.Println(grades[1])

firstOnes := grades[0:2]
lastOnes := grades[3:]
fmt.Println("firstOnes", firstOnes)
fmt.Println("lastOnes", lastOnes)

for i := 0; i < 5; i++ {
	fmt.Println("grades", i, grades[i])
}
...
[5.5 7.5 8 5.3 9.2]
Index 0 value 5.5
Index 1 value 7.5
Index 2 value 8
Index 3 value 5.3
Index 4 value 9.2
8
firstOnes [5.5 8]
lastOnes [5.3 9.2]
grades 0 5.5
grades 1 8
grades 2 8
grades 3 5.3
grades 4 9.2
```

You can also use the traditional form of **for** and also the **break** command:

```
for i := 0; i < 5; i++ {
	fmt.Println("grades", i, grades[i])
}
```

**If**

To finish this lesson, we have the **if** command:

```
if condition {

}
else {

}
```

The comparison operators in **Go** are:  ```==, !=, <, <=, >, and >=```

## Challenge

Given a vector containing "n" integers, indicate the largest repeated quantity of numbers. Example: 

```
{1,1,0,0,0,3,1,1,4,4,4,0,0,7,7,7,7,7,1}
```
In this example, the answer would be: 5, because the number 7 repeats 5 times.

