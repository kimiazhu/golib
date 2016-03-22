// Description: safego
// Author: ZHU HAIHUA
// Since: 2016-03-22 15:57
package safego

import (
	"fmt"
	log "github.com/kimiazhu/log4go"
	"testing"
	"time"
)

func TestPanic(t *testing.T) {
	Go(func() {
		a := []int{0}
		fmt.Println(a[1])
	})
	time.Sleep(time.Second)
}

func TestNoPanic(t *testing.T) {
	Go(func() {
		a := []int{0}
		fmt.Println(a[0])
	})
	time.Sleep(time.Second)
}

func TestWithLog4go(t *testing.T) {
	Go(func() {
		a := []int{0}
		fmt.Println(a[1])
	}, func(err interface{}) {
		// Critical will automatically print the stack
		log.Critical("panic catched: %v", err)
	})
	time.Sleep(time.Second)
}
