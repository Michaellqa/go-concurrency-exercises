//////////////////////////////////////////////////////////////////////
//
// DO NOT EDIT THIS PART
// Your task is to edit `main.go`
//

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

// RunMockServer pretends to be a video processing service. It
// simulates user interacting with the Server.
func RunMockServer() {
	free := User{ID: 0, IsPremium: false}
	premium := User{ID: 1, IsPremium: true}

	wg.Add(5)

	go createMockRequest(1, shortProcess, &free) // 4 sec -> true
	time.Sleep(1 * time.Second)

	go createMockRequest(2, longProcess, &premium) // 12 sec -> true
	time.Sleep(2 * time.Second)

	go createMockRequest(3, shortProcess, &free) // time left: [7...4...2...0] 6.5 sec -> false
	time.Sleep(1 * time.Second)

	go createMockRequest(4, longProcess, &free)     // 6.5 sec -> false
	go createMockRequest(5, shortProcess, &premium) // 8 sec -> true

	wg.Wait()
}

func createMockRequest(pid int, fn func(), u *User) {
	fmt.Println("UserID:", u.ID, "\tProcess", pid, "started.")
	res := HandleRequest(fn, u)

	if res {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "done.")
	} else {
		fmt.Println("UserID:", u.ID, "\tProcess", pid, "killed. (No quota left)")
	}

	wg.Done()
}

func shortProcess() {
	time.Sleep(4 * time.Second)
}

func longProcess() {
	time.Sleep(11 * time.Second)
}
