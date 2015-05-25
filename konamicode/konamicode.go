package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var (
	KONAMI    = [...]int{38, 38, 40, 40, 37, 39, 37, 39, 66, 65}
	jQuery    = jquery.NewJQuery
	document  = jQuery(js.Global.Get("document"))
	result    = jQuery("#result")
	keyupChan = make(chan int)
	posOk     = func(a, pos int) (currentOk bool, isLast bool) {
		if pos > len(KONAMI)-1 || pos < 0 {
			return false, false
		}
		return KONAMI[pos] == a, pos == len(KONAMI)-1
	}
	keyHandler = func(this *js.Object) {
		keyCode := this.Get("keyCode").Int()
		go func() {
			keyupChan <- keyCode
		}()
	}
)

func main() {

	document.On("keyup", keyHandler)
	go func() {
		var keyPos int
		for {
			select {
			case key := <-keyupChan:

				keyOk, isLast := posOk(key, keyPos)
				if keyOk && isLast {
					result.SetHtml("KONAMI").Show().FadeOut(2000)
					keyPos = 0
				} else if keyOk {
					keyPos += 1
				} else {
					keyPos = 0
				}
			}
		}
	}()
}
