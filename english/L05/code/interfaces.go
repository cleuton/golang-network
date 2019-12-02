package main

import (
	"fmt"
	"reflect"
)

// Vehicle this represents a vehicle behavior
type Vehicle interface {
	TurnOn() bool
	TurnOff() bool
	Move(direction int, speed int) bool
	Stop()
}

// Gas this represents a gasoline fueled object
type Gas interface {
	FillUp()
}

// Car this struct represents a car
type Car struct {
	model string
}

// TurnOn this method turns a car on
func (c Car) TurnOn() bool {
	fmt.Println("Turning the car on")
	return true
}

// TurnOff this method turns a car on
func (c Car) TurnOff() bool {
	fmt.Println("Turning the car off")
	return true
}

// Move this method makes a car move
func (c Car) Move(direction int, speed int) bool {
	fmt.Println("Moving the car")
	return true
}

// Stop this method stops a car
func (c Car) Stop() {
	fmt.Println("Stopping the car")
}

// FillUp this method Fills the tank of a Car
func (c Car) FillUp() {
	fmt.Println("Filling the car's tank")
}

// Truck this struct represents a Truck
type Truck struct {
	model string
}

// TurnOn this method turns a Truck on
func (t Truck) TurnOn() bool {
	fmt.Println("Turning the Truck on")
	return true
}

// TurnOff this method turns a Truck on
func (t Truck) TurnOff() bool {
	fmt.Println("Turning the Truck off")
	return true
}

// Move this method makes a Truck move
func (t Truck) Move(direction int, speed int) bool {
	fmt.Println("Moving the Truck")
	return true
}

// Stop this method stops a Truck
func (t Truck) Stop() {
	fmt.Println("Stopping the Truck")
}

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

func main() {
	newCar := Car{"Mustang"}
	fmt.Println(newCar, reflect.TypeOf(newCar))
	newCar.TurnOn()
	newCar.Move(1, 1)
	newCar.Stop()
	newCar.TurnOff()

	var newVehicle Vehicle
	fmt.Println(newVehicle, reflect.TypeOf(newVehicle))
	newVehicle = Car{"Dodge"}
	fmt.Println(newVehicle, reflect.TypeOf(newVehicle))
	newVehicle.TurnOn()
	newVehicle.Move(1, 1)
	newVehicle.Stop()
	newVehicle.TurnOff()

	newTruck := Truck{"GMC"}
	newVehicle = newTruck
	newVehicle.TurnOn()
	newVehicle.Move(1, 1)
	newVehicle.Stop()
	newVehicle.TurnOff()

	var aVehicle interface{}
	fmt.Println(aVehicle, reflect.TypeOf(aVehicle))
	aVehicle = newTruck
	fmt.Println(aVehicle, reflect.TypeOf(aVehicle))

	value, answer := newVehicle.(Truck)
	fmt.Println(value, answer)

	newVehicle = newCar
	v2, a2 := newVehicle.(Car)
	fmt.Println(v2, a2)
	v3, a3 := newVehicle.(Gas)
	fmt.Println(v3, a3)

	// Using pointer

	// Error: Laptop does not implement Computer (TurnOff method has pointer receiver)
	// var c Computer = Laptop{"Asus", false}
	var c Computer = &Laptop{"Asus", false}
	fmt.Println(c)
	c.TurnOn()

}
