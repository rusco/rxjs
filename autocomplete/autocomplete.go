package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"time"
)

var (
	jQuery = jquery.NewJQuery
	input  = jQuery("#textInput")
	ul     = jQuery("#results")

	keyupChan  = make(chan string)
	keyHandler = func(this *js.Object) {

		keyCode := this.Get("target").Get("value").String()
		if len(keyCode) < 3 {
			return
		}
		go func() {
			keyupChan <- keyCode
		}()
	}
	search = func(term string) {

		ajaxOpt := js.M{
			"url":      "http://en.wikipedia.org/w/api.php",
			"dataType": "jsonp",
			"data": js.M{
				"action": "opensearch",
				"format": "json",
				"search": js.Global.Call("encodeURI", term),
			},
			"error": func(error *js.Object) {
				ul.Empty()
				jQuery("<li>" + error.Get("errorThrown").String() + "</li>").AppendTo(ul)
			},
			"success": func(value *js.Object) {

				results := value.Index(1).Interface().([]interface{})
				ul.Empty()
				for _, r := range results {
					jQuery("<li>" + r.(string) + "</li>").AppendTo(ul)
				}
			},
		}
		jquery.Ajax(ajaxOpt)
	}
)

func main() {

	input.On("keyup", keyHandler)

	go func() {
		searchFor, modified, interval := "", false, time.Tick(800*time.Millisecond)

		for {
			select {
			case <-interval:
				if modified {
					search(searchFor)
					modified = false
				}
			case searchFor = <-keyupChan:
				modified = true
			}
		}
	}()
}
