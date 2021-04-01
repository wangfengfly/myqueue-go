package main

import (
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

func TestAdd(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		myqueue.Add(1)
		wg.Done()
	}()
	go func() {
		myqueue.Add(2)
		wg.Done()
	}()

	wg.Wait()
	if myqueue.Len()!=2 {
		t.Errorf("myqueue Len is %d, but %d is expected.", myqueue.Len(), 2)
	}
}


func TestGet(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	var item1 int
	var item2 int
	go func() {
		v,_ := myqueue.Get()
		item1 = v.(int)
		wg.Done()
	}()

	go func() {
		v,_ := myqueue.Get()
		item2 = v.(int)
		wg.Done()
	}()

	wg.Wait()

	if item1==item2 {
		t.Errorf("read duplicated items from myqueue")
	}
}
