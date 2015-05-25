package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

const (
	MOUSEMOVE   = "mousemove"
	MOUSEUP     = "mouseup"
	MOUSEDOWN   = "mousedown"
	MOUSEEVENTS = MOUSEMOVE + " " + MOUSEUP + " " + MOUSEDOWN
	CLIENTX     = "clientX"
	CLIENTY     = "clientY"
)

type Offset struct {
	x int
	y int
}

var (
	jQuery   = jquery.NewJQuery
	canvas   = js.Global.Get("document").Call("getElementById", "tutorial")
	ctx      = canvas.Call("getContext", "2d")
	tutorial = jQuery("#tutorial")

	mouseupChan   = make(chan bool)
	mousemoveChan = make(chan Offset)
	mousedownChan = make(chan Offset)

	mouseHandler = func(this *js.Object) {
		evtType := this.Get("type").String()

		if evtType == MOUSEUP {
			go func() {
				mouseupChan <- true
			}()
			return
		}

		offset := jQuery(this.Get("target")).Offset()
		x := this.Get(CLIENTX).Int() - offset.Left
		y := this.Get(CLIENTY).Int() - offset.Top

		go func() {
			if evtType == MOUSEMOVE {
				mousemoveChan <- Offset{x: x, y: y}
			} else if evtType == MOUSEDOWN {
				mousedownChan <- Offset{x: x, y: y}
			}
		}()
	}
)

func main() {

	tutorial.On(MOUSEEVENTS, mouseHandler)

	go func() {
		var prevX, prevY int
		var dragStarted bool

		for {
			select {
			case md := <-mousedownChan:
				prevX, prevY, dragStarted = md.x, md.y, true
				ctx.Call("beginPath")
			case mm := <-mousemoveChan:
				if dragStarted {
					ctx.Call("moveTo", prevX, prevY)
					ctx.Call("lineTo", mm.x, mm.y)
					ctx.Call("stroke")
					prevX, prevY = mm.x, mm.y
				}
			case <-mouseupChan:
				dragStarted = false
			}
		}
	}()
}
