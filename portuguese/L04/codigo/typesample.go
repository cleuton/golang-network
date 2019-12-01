package main

import (
	"fmt"
	"math"
	"reflect"
)

// Student This is a student
type Student struct {
	name string
}

// Course This is a course
type Course struct {
	description string
	students    []Student
}

func (c *Course) register(s Student) {
	c.students = append(c.students, s)
}

// EadCourse This is another course type
type EadCourse struct {
	course  Course
	website string
}

// NewFloat float with steroids
type NewFloat float64

func (nf NewFloat) roundBy(places float64) NewFloat {
	nplaces := math.Pow(10.00, places)
	return NewFloat(math.Round(float64(nf)*nplaces) / nplaces)
}

func main() {
	var newStudent Student
	newStudent.name = "John Doe"
	fmt.Println(newStudent, reflect.TypeOf(newStudent))

	pst := &newStudent
	fmt.Println(pst.name, reflect.TypeOf(newStudent))
	fmt.Println((*pst).name, reflect.TypeOf(newStudent))

	var other Student
	other.name = "Other Student"
	fmt.Println(pst.name, other.name)

	p := Student{"Pamela"}
	fmt.Println(p.name)

	engineering := Course{"Engineering", make([]Student, 0)}
	engineering.register(newStudent)
	engineering.register(p)
	fmt.Println(engineering)

	var n NewFloat = 5.0293019384
	fmt.Println(n.roundBy(2))
	fmt.Println(n.roundBy(3))

	ead := EadCourse{Course{"New EAD", make([]Student, 0)}, "http://eadcourse"}
	ead.course.register(newStudent)
	fmt.Println(ead)
}
