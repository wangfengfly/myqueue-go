package main

type Queue interface {
	Add(item interface{})
	Get() (item interface{}, shutdown bool)
	Len() int
	Done(item interface{})
	ShutDown()
	ShuttingDown() bool
}
