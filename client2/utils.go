package main

import (
	"github.com/rivo/tview"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func Center(width, height int, p tview.Primitive) tview.Primitive {
	return tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(p, height, 1, true).
			AddItem(tview.NewBox(), 0, 1, false), width, 1, true).
		AddItem(tview.NewBox(), 0, 1, false)
}

func reverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func checkForAts(username, line string) bool {
	return strings.Contains(line, "@"+strings.Title(username))
}

func doEvery(d time.Duration, f func() error) {
	for range time.Tick(d) {
		f()
	}
}

func getTimeNow() string {
	t := time.Now()
	now := t.Format("15:04:05")
	return now
}

func startSpinner(button *tview.Button, action func()) {
	done := make(chan bool)

	go func() {
		action()
		done <- true
		close(done)
	}()

	go func() {
		//spinners := []string{"|", "\\", "-", "/"}
		var i int
		for {
			select {
			case _ = <-done:
				//app.QueueUpdateDraw(func() {
				//	button.SetLabel("Confirm")
				//})
				return
			case <-time.After(200 * time.Millisecond):
				//spin := i % len(spinners)
				//app.QueueUpdateDraw(func() {
				//	button.SetLabel(spinners[spin] + " Loading")
				//})
				i++
			}
		}
	}()
}
