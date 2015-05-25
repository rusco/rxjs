package main

//demo only, use html5 "draggable" attribute instead
import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

const (
	MOUSEMOVE = "mousemove"
	MOUSEUP   = "mouseup"
	MOUSEDOWN = "mousedown"
	CLIENTX   = "clientX"
	CLIENTY   = "clientY"
)

type Pos struct {
	top  int
	left int
}

var (
	jQuery     = jquery.NewJQuery
	dragTarget = jQuery("#dragTarget")
	document   = jQuery(js.Global.Get("document"))

	mouseupChan   = make(chan bool)
	mousemoveChan = make(chan Pos)
	mousedownChan = make(chan Pos)

	mousedownFn = func(this *js.Object) {
		this.Call("preventDefault")

		offset := dragTarget.Offset()
		left := this.Get(CLIENTX).Int() - offset.Left
		top := this.Get(CLIENTY).Int() - offset.Top

		go func(l, t int) {
			mousedownChan <- Pos{left: l, top: t}
		}(left, top)
	}
	mouseupFn = func(this *js.Object) {
		go func() {
			mouseupChan <- true
		}()
	}
	mousemoveFn = func(this *js.Object) {
		left := this.Get(CLIENTX).Int()
		top := this.Get(CLIENTY).Int()

		go func(l, t int) {
			mousemoveChan <- Pos{left: l, top: t}
		}(left, top)
	}
)

func main() {

	dragTarget.On(MOUSEDOWN, mousedownFn)
	dragTarget.On(MOUSEUP, mouseupFn)
	document.On(MOUSEMOVE, mousemoveFn)

	go func() {
		var imgOffsetTop, imgOffsetLeft int
		var dragStarted bool

		for {
			select {
			case md := <-mousedownChan:
				imgOffsetTop, imgOffsetLeft, dragStarted = md.top, md.left, true
			case mm := <-mousemoveChan:
				if dragStarted {
					dragTarget.SetCss(js.M{"top": mm.top - imgOffsetTop, "left": mm.left - imgOffsetLeft})
				}
			case <-mouseupChan:
				dragStarted = false
			}
		}
	}()
}
