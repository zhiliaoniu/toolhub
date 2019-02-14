package main

import (
	"errors"
	"fmt"
)

type Animal struct {
	name string
}

func (a Animal) GetName() (n string, err error) {
	if a.name == "" {
		return "", errors.New("Animal name is empty")
	}
	n = a.name
	return
}

type Rabbit struct {
	Animal
	name string
}

func NewRabbit(n1, n2 string) *Rabbit {
	rabbit := new(Rabbit)
	//animal := new(Animal)
	rabbit.Animal = Animal{}
	rabbit.Animal.name = n1
	rabbit.name = n2
	return rabbit
}

//func (r Rabbit) GetName() (n string, err error) {
//	if r.name == "" {
//		return "", errors.New("Rabbit name is empty")
//	}
//	n = r.name
//	return
//}

func main() {
	a := &Animal{"animal1"}
	animal1, err := a.GetName()
	if err != nil {
		fmt.Println("a.GetName animal1=", animal1, "err:", err)
	} else {
		fmt.Println(animal1)
	}

	//var b *Animal = new(Rabbit)
	//b := new(Rabbit) //{Animal: {"cc"}, name: "bb"}
	//b := NewRabbit("animal_name", "rabbit_name")
	//b.Animal.name = "aa"
	b := &Rabbit{Animal: *new(Animal), name: "rabbit_name"}
	rabbit1, err := b.GetName()
	if err != nil {
		fmt.Println("b.GetName rabbit1=", rabbit1, "err:", err)
	} else {
		fmt.Println(rabbit1)
	}

	type S struct {
		name string
	}
	arr := make([]*S, 0)
	s := &S{"hah"}
	arr = append(arr, s)
	s2 := &S{"hah2"}
	arr = append(arr, s2)
	for k, v := range arr {
		fmt.Printf("k:%d \tv:%s\n", k, v.name)
	}
	arr = arr[:1]
	for k, v := range arr {
		fmt.Printf("k:%d \tv:%s\n", k, v.name)
	}
	m2 := [5]int{1, 2, 3, 4, 5}
	arr2 := m2[1:4]
	for k, v := range arr2 {
		fmt.Printf("k:%d \tv:%d\n", k, v)
	}
	arr2 = m2[:2]
	for k, v := range arr2 {
		fmt.Printf("k:%d \tv:%d\n", k, v)
	}
	fmt.Println("\n=======================\n")

}
