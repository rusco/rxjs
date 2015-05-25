package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"strconv"
	"time"
)

const (
	MOUSEMOVE               = "mousemove"
	TEXT                    = "time flies like an arrow"
	DELAY     time.Duration = 120
)

var jQuery = jquery.NewJQuery

type Offset struct {
	offsetX int
	offsetY int
}

func main() {

	container := jQuery("#textContainer")
	document := jQuery(js.Global.Get("document"))
	offsetChan := make(chan Offset)
	visible := false

	for idx, txt := range TEXT {
		container.Append("<span id='span" + strconv.Itoa(idx) + "' style='position:absolute;'>" + string(txt) + "</span>")
	}

	document.On(MOUSEMOVE, func(this *js.Object) {

		if !visible {
			jQuery("#textContainer").SetCss("color", "black")
			visible = true
		}

		offset := container.Offset()
		offX := this.Get("pageX").Int() - offset.Left
		offY := this.Get("pageY").Int() - offset.Top

		go func(ox int, oy int) {
			offsetChan <- Offset{offsetX: ox, offsetY: oy}
		}(offX, offY)
	})

	go func() {
		for {
			oc := <-offsetChan
			for idx := range TEXT {
				go func(i, x, y int) {

					time.Sleep(time.Duration(i) * DELAY * time.Millisecond)
					jQuery("#span" + strconv.Itoa(i)).SetCss(js.M{"top": y, "left": x + i*10 + 15})

				}(idx, oc.offsetX, oc.offsetY)
			}
		}
	}()
}
