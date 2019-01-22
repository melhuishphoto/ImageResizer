//go:generate goversioninfo
//go:generate rice embed-go
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/shibukawa/configdir"
	"github.com/zserge/lorca"

	"github.com/melhuishphoto/ImageResizer/resize"
)

func main() {
	url := Serve()
	dir := configdir.New("ImageResizer", "Resizer").QueryCacheFolder().Path
	ui, e := lorca.New(url, dir, 550, 250)
	if e != nil {
		panic(e)
	}
	ui.Bind("choose", Choose(ui))
	ui.Bind("resize", Resize(ui))
	ui.Eval(`window.moveTo((screen.width-window.outerWidth)/2,(screen.height-window.outerHeight)/2)`)
	<-ui.Done()
}

func Choose(ui lorca.UI) func() string {
	return resize.ChooseDir
}

func Resize(ui lorca.UI) func(string, int, int, int) {
	return func(dir string, imgSize, thumbSize, quality int) {
		opts := &resize.Options{
			ImageSize: imgSize,
			ThumbSize: thumbSize,
			Quality:   quality,
		}
		e := resize.Resize(dir, opts, func(progress, total int) {
			ui.Eval(fmt.Sprintf(`setProgress(%v,%v);`, progress, total))
		})
		if e != nil {
			log.Println(e)
		}
	}
}

func Serve() string {
	l, e := net.Listen("tcp", "localhost:")
	if e != nil {
		panic(e)
	}
	l.Close()

	go func(l net.Listener) {
		web := rice.MustFindBox("./assets")
		http.ListenAndServe(l.Addr().String(), http.FileServer(web.HTTPBox()))
	}(l)

	//time.Sleep(time.Millisecond * 50)

	return fmt.Sprint("http://", l.Addr().String())
}
