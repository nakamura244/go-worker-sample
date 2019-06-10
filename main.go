package main

import (
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"
)

func taskA()  {
	fmt.Println("task A start")
	// do something
	time.Sleep(3 * time.Second)
	fmt.Println("task A end")
}

func taskB()  {
	fmt.Println("task B start")
	// do something
	time.Sleep(3 * time.Second)
	fmt.Println("task B end")
}

func task(task func(), stopCh, doneCh chan struct{}) {
	defer func() {
		close(doneCh)
	}()
	f := reflect.ValueOf(task)
	for {
		select {
		case <-stopCh:
			fmt.Printf("%v:catch the stop request.\n", runtime.FuncForPC(f.Pointer()).Name())
			return
		default:
			task()
		}
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	stopCh := make(chan struct{})
	taskADoneCh := make(chan struct{})
	taskBDoneCh := make(chan struct{})
	go task(taskA, stopCh, taskADoneCh)
	go task(taskB, stopCh, taskBDoneCh)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		fmt.Printf("Catch the signal at main: %v\n", s)
		close(stopCh)
	}()
	fmt.Println("Waiting for signal at main")
	<-stopCh
	fmt.Println("Exit from main")
	<-taskADoneCh
	<-taskBDoneCh
	fmt.Println("Exit All")
}
