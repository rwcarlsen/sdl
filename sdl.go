package sdl

/*
#cgo pkg-config: sdl2

#include "SDL2/SDL.h
*/
import "C"
import (
	"errors"
	"log"
	"runtime"
	"unsafe"
)

func init() {
	log.SetFlags(0)
	success := bool(C.SDL_Init(C.SDL_INIT_EVERYTHING) != 0)
	if !success {
		err := C.GoString(C.SDL_GetError())
		log.Fatalf("SDL_Init Error: %v", err)
	}
}

const (
	WinposCentered  = C.SDL_WINDOWPOS_CENTERED
	WinposUndefined = C.SDL_WINDOWPOS_UNDEFINED
)

type Window struct {
	win *C.SDL_Window
}

func NewWindow(title string, xpos, ypos, w, h int) (*Window, error) {
	cs := C.CString(title)
	defer C.free(unsafe.Pointer(cs))

	w := C.SDL_CreateWindow(cs, C.int(xpos), C.int(ypos), C.int(w), C.int(h), 0)
	if w == nil {
		return nil, sdlerr()
	}

	win := &Window{w}
	runtime.SetFinalizer(w, freewin)
	return win, nil
}

func (w *Window) NewRenderer() (*Renderer, error) {
	ren := C.SDL_CreateRenderer(w.win, -1, C.SDL_RENDERER_ACCELERATED|C.SDL_RENDERER_PRESENTVSYNC)
	if ren == nil {
		return nil, sdlerr()
	}

	r := &Renderer{ren}
	runtime.SetFinalizer(r, freeren)
	return r, nil
}

func (w *Window) Show() { C.SDL_ShowWindow(w.win) }

func (w *Window) Hide() { C.SDL_HideWindow(w.win) }

func (w *Window) Maximize() { C.SDL_MaximizeWindow(w.win) }

func (w *Window) Minimize() { C.SDL_MinimizeWindow(w.win) }

func (w *Window) Fullscreen() { C.SDL_SetWindowFullscreen(w.win, C.SDL_WINDOW_FULLSCREEN) }

func (w *Window) Windowed() { C.SDL_SetWindowFullscreen(w.win, 0) }

func (w *Window) SetGrab(grab bool) { C.SDL_SetWindowGrab(w.win, grab) }

func sdlerr() error { return errors.New(C.GoString(C.SDL_GetError())) }

type Renderer struct {
	ren *C.SDL_Renderer
}

func freewin(w *Window) { C.SDL_DestroyWindow(w.win) }

func freeren(r *Renderer) { C.SDL_DestroyRenderer(r.ren) }
