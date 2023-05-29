package main

import (
	"fmt"
	"log"
	"time"
)

func produceMessages(c chan<- string) {
	for i := 0; i < 10; i++ {
		fmt.Println("Produce a new message")
		c <- fmt.Sprint("Message", i)
	}
}

func consumeMessages(c <-chan string) {
	for {
		message := <-c
		fmt.Println(message)
	}
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	var sum V
	for _, v := range m {
		sum += v
	}
	return sum
}

func SumGeneric[T int8 | int16 | int32 | int64 | int](ints []T) T {
	var sum T
	for _, i := range ints {
		sum += i
	}
	return sum
}

type List[T any] struct {
	next *List[T]
	val  T
}

func (li *List[T]) Read() T {
	return li.val
}

func (li *List[T]) Add(n T) {
	li.next = new(List[T])
	li.next.val = n
}

func main() {
	fmt.Println(SumGeneric([]int16{123, 212}))

	listOfInts := new(List[int])
	listOfStrings := new(List[string])
	listOfInts.Add(3)
	listOfStrings.Add("asas")
	fmt.Println(listOfInts.Read())
	fmt.Println(listOfInts.next.Read())
	fmt.Println(listOfStrings.Read())
	fmt.Println(listOfStrings.next.Read())

	return
	strChannel := make(chan string, 3)
	go produceMessages(strChannel)
	go consumeMessages(strChannel)
	time.Sleep(time.Second)

	log.Println("Go fuer KI")

	einTupel := Tupel{3, 4}
	einTupel.Add(1, 4)

	fmt.Println("a:", einTupel.a, " b:", einTupel.b)
	// var i interface{} = einTupel

	myMap := make(map[int]string)

	defer fmt.Println("Ende")
	for i := 0; i < 3; i++ {
		fmt.Print(i)
	}

	var msg = "Hello"
	fmt.Println(msg)
	var i interface{} = einTupel
	switch v := i.(type) {
	case fmt.Stringer:
		fmt.Println("hat String() Funktion", v.String())
	default:
		fmt.Println("leider nein")
	}

	var tuu Tupel
	fmt.Println(tuu)

	// var pointr *int
	// pointrr := *pointr / 7
	// fmt.Println(pointrr)
	var com complex128
	fmt.Println(com)

	// var a uint = 5
	// var b int = 7
	// var c = a + b
}

type Tupel struct {
	a, b int
}

func (t *Tupel) Add(c, d int) {
	t.a += c
	t.b += d
}

func (t Tupel) String() string {
	return fmt.Sprint("Tupel(a:", t.a, ", b:", t.b, ")")
}

type Stringer interface {
	String() string
}
