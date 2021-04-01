package main

import (
	"sync"
	"time"
)

type MyQueue struct {
	data []interface{}
	mutex sync.Mutex
	processing map[interface{}]bool
	processingmutex sync.Mutex
	exists map[interface{}]bool
	isshutdown bool
	isshuttingdown bool
}

func NewMyQueue() *MyQueue {
	var myqueue MyQueue
	myqueue.processing = make(map[interface{}]bool)
	myqueue.exists = make(map[interface{}]bool)
	return &myqueue
}

func (q *MyQueue) Add(item interface{}) {
	if q.isshuttingdown {
		return
	}

	q.mutex.Lock()
	defer q.mutex.Unlock()
	//item在处理之前只能被添加一次
	if _,ok:=q.exists[&item]; ok {
		return
	}

	q.exists[&item] = true
	q.data = append(q.data, item)

	return
}

func (q *MyQueue) Get() (item interface{}, shutdown bool) {
	if q.Len()==0 {
		return nil,q.isshutdown
	}

	q.mutex.Lock()
	q.processingmutex.Lock()
	defer  q.mutex.Unlock()
	defer q.processingmutex.Unlock()

	item = q.data[0]
	q.data = q.data[1:]
	q.processing[&item] = true
	delete(q.exists, &item)

	return item,q.isshutdown
}

func (q *MyQueue) Len() int {
	q.mutex.Lock()
	defer  q.mutex.Unlock()
	return len(q.data)
}

func (q *MyQueue) Done(item interface{}) {
	q.processingmutex.Lock()
	delete(q.processing, &item)
	q.processingmutex.Unlock()
}

func (q *MyQueue) execShutDown() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if q.Len()==0 {
				q.isshutdown = true
				return
			}
		}
	}
}

func (q *MyQueue) ShutDown() {
	q.isshuttingdown = true
	go q.execShutDown()
}

func (q *MyQueue) ShuttingDown() bool {
	return q.isshuttingdown
}


