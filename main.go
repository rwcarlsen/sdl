package main

import (
	"image/color"
	"log"
	"time"

	"github.com/rwcarlsen/sdl2/sdl"
)

func main() {
	w, h := 640, 480

	win, err := sdl.NewWindow("hello world", sdl.WinposCentered, sdl.WinposCentered, w, h)
	if err != nil {
		log.Fatal(err)
	}

	surf, err := sdl.NewSurface(w, h)
	if err != nil {
		log.Fatal(err)
	}

	err = surf.FillRect(nil, color.RGBA{252, 0, 0, 0})
	if err != nil {
		log.Fatal(err)
	}

	tex, err := win.NewTexture(surf)
	if err != nil {
		log.Fatal(err)
	}

	win.Clear()
	if err := win.Copy(tex, nil, nil); err != nil {
		log.Fatal(err)
	}
	win.Present()

	<-time.After(5 * time.Second)
}
