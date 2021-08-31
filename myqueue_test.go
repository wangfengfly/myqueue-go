package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var myqueue *MyQueue

func setup() {
	myqueue = NewMyQueue()
}

func teardown() {
	myqueue = nil
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

type Student struct {
	Name string
}

func worker(wg *sync.WaitGroup, item *Student) {
	defer wg.Done()

	ret := myqueue.Add(item)
	fmt.Println(fmt.Sprintf("add to myqueue result: %v", ret))

}

func TestAdd(t *testing.T) {
	student := Student{"Tom"}

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go worker(&wg, &student)
	}

	wg.Wait()

	if myqueue.Len() != 1 {
		t.Errorf("the myqueue length shoud be 1.")
	}
}

func TestGet(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	var item1 *Student
	var item2 *Student
	go func() {
		v, _ := myqueue.Get()
		item1, _ = v.(*Student)
		wg.Done()
	}()

	go func() {
		v, _ := myqueue.Get()
		item2, _ = v.(*Student)
		wg.Done()
	}()

	wg.Wait()

	if item1 == item2 {
		t.Errorf("read duplicated items from myqueue")
	}
}
