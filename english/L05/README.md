![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Interfaces

An interface is a data type, just like a **struct**, which serves to declare a **set of methods** that an object must implement.

In addition to the methods declared inside a **struct** (or in type, as we have seen) we can add methods to an interface. Thus we can separate objects according to their behavior.

Since we do not have the concept of **inheritance** we need a way to generalize the behavior of similar objects. For example, what behaviors do all **vehicles** have in common?

- TurnOn
- TurnOff
- Move
- Stop

We can create an interface that abstracts these behaviors. See my [**sample**](./code/interfaces.go): 

```
// Vehicle this represents a vehicle behavior
type Vehicle interface {
	TurnOn() bool
	TurnOff() bool
	Move(direction int, speed int) bool
	Stop()
}
```

An interface is a type, so we can create variables from it:

```
var newVehicle Vehicle
fmt.Println(newVehicle, reflect.TypeOf(newVehicle))
...
<nil> <nil>
```

Because the interface is just a set of methods, it has no static value, and if not initialized with a compatible object, the result is **nil** (null).

An interface can take as values ​​any objects that implement it. In our example: **Car** and **Truck** implement the **Vehicle** interface. How do I know this? Well, look at the [**source code**] (./code/intefaces.go) and note that both **structs** declare ALL **Vehicle** methods, so they implement this interface. 

```
var newVehicle Vehicle
newVehicle = Car{"Dodge"}
newTruck := Truck{"GMC"}
newVehicle = newTruck
```

This is only possible if the structs implement ALL interface methods.

We can use any interface methods on any object, and the method invoked will be the one declared in that **struct**, ie **implicit polymorphism**: 

```
newVehicle = Car{"Dodge"}
newVehicle.TurnOn()
...
Turning the car on
...
newTruck := Truck{"GMC"}
newVehicle = newTruck
...
Turning the Truck on
```

In this example, the **TurnOn()** method invoked is always the one belonging the concrete type of the object, since the interface has no implementation of the methods.

## Empty interface

All the custom types implement the **empty interface** ie: **interface{}**: 

```
var aVehicle interface{}
fmt.Println(aVehicle, reflect.TypeOf(aVehicle))
...
<nil> <nil>
```

A variable or an argument of type **interface{}** can take any object, since everyone implements it:

```
aVehicle = newTruck
fmt.Println(aVehicle, reflect.TypeOf(aVehicle))
...
{GMC} main.Truck
```

An object can implement more than one interface, for example, every **Car** object also implements the **Gas​​** interface:

```
// Gas this represents a gasoline fueled object
type Gas interface {
	FillUp()
}
...
// FillUp this method Fills the tank of a Car
func (c Car) FillUp() {
	fmt.Println("Filling the car's tank")
}
```

When we have a variable **interface** variable, then we can check which interfaces the underlining object implements. For this we use the syntax: ```i.(Type)```:

```
newVehicle = newCar
v2, a2 := newVehicle.(Car)
fmt.Println(v2, a2)
v3, a3 := newVehicle.(Gas)
fmt.Println(v3, a3)
...
{Mustang} true
{Mustang} true
```

The **i.(Type)** syntax tests whether the object pointed to by an **interface** variable implements another interface, and returns the value and a **bool** indicating the response.

## Ponter vs value receiver

We have already seen that a method can have either a **pointer** or a **value** receiver. For example, this method is given a value, so it cannot modify the object:

```
// FillUp this method Fills the tank of a Car
func (c Car) FillUp() {
	fmt.Println("Filling the car's tank")
}
```

However, in lesson [**L04**](../L04) we can have methods whose receiver is a pointer:

```
func (c *Course) register(s Student) {
	c.students = append(c.students, s)
}
```

This is because the **register()** method changes properties of the **Course** object. When we have **struct** variables the receiver type doesn't matter. But when we use interfaces, we may have problems. See the code below:

```
// Computer this represents a Computer behavior
type Computer interface {
	TurnOn() bool
	TurnOff() bool
}

// Laptop this struct represents a Laptop
type Laptop struct {
	model string
	on    bool
}

// TurnOn this method turns a laptop on
func (l *Laptop) TurnOn() bool {
	fmt.Println("Turning the laptop on")
	l.on = true
	return true
}

// TurnOff this method turns a laptop off
func (l *Laptop) TurnOff() bool {
	fmt.Println("Turning the laptop off")
	l.on = false
	return true
}
...
var c Computer = Laptop{"Asus", false}
```

The code would result in an error. In the last line we would see an error of type: "Laptop does not implement Computer (TurnOff method has pointer receiver)". When we use methods of a **struct** object, the receiver kind does not matter, but when we want to use interface variables, we have to know that. If a method has a **pointer** receiver, we have to use pointer notation:

```
var c Computer = &Laptop{"Asus", false}
fmt.Println(c)
c.TurnOn()
```

## Challenge

Take the challenge from the previous chapter and declare **Course** methods using an interface.




