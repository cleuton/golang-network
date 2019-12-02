![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Modularization: Structs

In today's programming languages, there is the concept of **Object Orientation**, allowing us to declare **classes** and instantiate **objects** from them. **Go** does not have this concept.

**Go has no concept of classes**

**Go** lets us create complex **Types**. These types are represented by the **interface** **Type** (we will see later).

A **Type** is a **struct** (a complex type) like this:

```
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	name string
}

func main() {
	var newStudent Student
	newStudent.name = "John Doe"
	fmt.Println(newStudent, reflect.TypeOf(newStudent))
}
...
{John Doe} main.Student
```

Open file [**typesample.go**] (./code/typesample.go).

We created a **struct** called **Student** containing only one member: The student's name. Pay attention to the use of capital letters in the name **Student**, it allows the type to be exported to other packages. Anything you declare capitalized in the package scope will be exported to other packages, but you must use the full declaration (with **var**).

After we run, we see the content of the **Student** object and its type name, which is the package name and the struct name.

**Pointers and memory**

**Go** has pointers. We can get the address of a variable and store it in a pointer with the "&" operator:
```
pst := &newStudent
fmt.Println(pst.name, reflect.TypeOf(newStudent))
...
John Doe main.Student
```

The **pst** variable now points to the object (memory location) of the variable **newStudent**. In other languages, such as C ++, we would need to *dereference* explicitly using the "*" operator:

```
fmt.Println((*pst).name, reflect.TypeOf(newStudent))
```

When we declare a variable using a **struct** we create a new memory area for it:

```
var other Student
other.name = "Other Student"
fmt.Println(pst.name, other.name)
...
John Doe Other Student
```

As we can see, each variable has its own memory area. It is actually a pointer to a created memory area.

And we can initialize a struct using a **literal struct** that contains the values ​​in the same order as struct types: 

```
p := Student{"Pamela"}
fmt.Println(p.name)
```

**Methods**

We can create **functions** within **structs**, better known as **methods** in other programming languages. For example, if we wanted to create a *class* **Course** that contains a description and student list, we would do this in Java:

```
class Course {
        String description;
        List<Student> students;
}
```

In **Go** an array can be of type **struct**, but not static arrays! Your type is declared along with its size! Only **slices** are dynamic. So how do we declare **struct**?

```
// Course This is a course
type Course struct {
	description string
	students    []Student
}
```

*(First of all **struct** put a comment with the name and some bla-bla-bla, otherwise the compiler will be disturbing you)*

I created a variable **students** which is a **slice** (a vector with no size declaration) and will create a method for adding new students:

```
func (c *Course) register(s Student) {
	c.students = append(c.students, s)
}
```

A **method** is a function that receives a **receiver** (before its name). A **receiver** is a function initiator and is what turns it into a method. The **receiver** can be the value or a pointer to an instance of the **struct**. If we just want to use the properties of the **struct**, we inform **receiver** as a common variable, but if we want to modify any property of the **struct**, we inform it as **pointer**.

In this case, the **register()** method gets a **pointer** to a **Course** object (note the asterisk), so it can modify any property of the **Course** object. And he does it by appending the new student to the **slice of Students**.

Now let's look at how we use **Course** type and how we register new students:

```
engineering := Course{"Engineering", make([]Student, 0)}
engineering.register(newStudent)
engineering.register(p)
fmt.Println(engineering)
...
{Engineering [{John Doe} {Pamela}]}
```

**Inheritance**

**Go** does not have the concept of inheritance, in which one class inherits the properties and methods of another. It has the concept of composition, in which one class can include another and access its properties. For example, let's create a distance learning course modality:

```
// EadCourse This is another course type
type EadCourse struct {
	course  Course
	website string
}
...
ead := EadCourse{Course{"New EAD", make([]Student, 0)}, "http://eadcourse"}
ead.course.register(newStudent)
fmt.Println(ead)
...
"{{New EAD [{John Doe}]} http://eadcourse}"
```

**Methods in other types**

We can declare methods in other types that we create. For example, I will create a method for rounding a **float64** to a certain number of decimal places:

```
// NewFloat float with steroids
type NewFloat float64

func (nf NewFloat) roundBy(places float64) NewFloat {
	nplaces := math.Pow(10.00, places)
	return NewFloat(math.Round(float64(nf)*nplaces) / nplaces)
}
...
var n NewFloat = 5.0293019384
fmt.Println(n.roundBy(2))
fmt.Println(n.roundBy(3))
...
5.03
5.029
```

As you can see, I created a method for a **NewFloat** class derived from **float64**.

## Challenge

Create a method for removing students from a course.

Tip: Use **slice** to remove students (append (slice [: i], slice [i + 1:] ...)) and use a **for** to find out the student's index (i) be removed.