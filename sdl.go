package sdl

/*
#cgo pkg-config: sdl2

#include "SDL2/SDL.h"

SDL_Rect* make_rect(int x, int y, int w, int h) {
	SDL_Rect* r = malloc(sizeof(SDL_Rect));
	(*r).x = x;
	(*r).y = y;
	(*r).w = w;
	(*r).h = h;
	return r;
};

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
	ren *C.SDL_Renderer
}

func NewWindow(title string, xpos, ypos, w, h int) (*Window, error) {
	cs := C.CString(title)
	defer C.free(unsafe.Pointer(cs))

	ww := C.SDL_CreateWindow(cs, C.int(xpos), C.int(ypos), C.int(w), C.int(h), C.Uint32(0))
	if ww == nil {
		return nil, sdlerr()
	}

	r := C.SDL_CreateRenderer(ww, -1, C.SDL_RENDERER_ACCELERATED|C.SDL_RENDERER_PRESENTVSYNC)
	if r == nil {
		return nil, sdlerr()
	}

	win := &Window{ww, r}
	runtime.SetFinalizer(win, freewin)
	return win, nil
}

func (w *Window) Copy(tex *Texture, src, dst *Rect) error {
	var csrc *C.SDL_Rect = nil
	if src != nil {
		csrc = sdlRect(src)
		defer C.free(unsafe.Pointer(csrc))
	}

	var cdst *C.SDL_Rect = nil
	if dst != nil {
		cdst = sdlRect(dst)
		defer C.free(unsafe.Pointer(cdst))
	}

	status := C.SDL_RenderCopy(w.ren, tex.tex, csrc, cdst)
	if status != 0 {
		return sdlerr()
	}
	return nil
}

func (w *Window) Present() { C.SDL_RenderPresent(w.ren) }

func (w *Window) Clear() { C.SDL_RenderClear(w.ren) }

func (w *Window) NewTexture(s *Surface) (*Texture, error) {
	t := C.SDL_CreateTextureFromSurface(w.ren, s.surf)
	if t == nil {
		return nil, sdlerr()
	}

	tex := &Texture{t}
	runtime.SetFinalizer(tex, freetex)
	return tex, nil
}

func (w *Window) Show() { C.SDL_ShowWindow(w.win) }

func (w *Window) Hide() { C.SDL_HideWindow(w.win) }

func (w *Window) Maximize() { C.SDL_MaximizeWindow(w.win) }

func (w *Window) Minimize() { C.SDL_MinimizeWindow(w.win) }

func (w *Window) Fullscreen() { C.SDL_SetWindowFullscreen(w.win, C.SDL_WINDOW_FULLSCREEN) }

func (w *Window) Windowed() { C.SDL_SetWindowFullscreen(w.win, 0) }

func (w *Window) SetGrab(grab bool) {
	if grab {
		C.SDL_SetWindowGrab(w.win, 1)
	} else {
		C.SDL_SetWindowGrab(w.win, 0)
	}
}

func freewin(w *Window) {
	C.SDL_DestroyRenderer(w.ren)
	C.SDL_DestroyWindow(w.win)
}

type Surface struct {
	surf *C.SDL_Surface
}

func NewSurface(w, h int) (*Surface, error) {
	var curr C.SDL_DisplayMode
	if C.SDL_GetCurrentDisplayMode(0, &curr) != 0 {
		return nil, sdlerr()
	}
	format := C.SDL_AllocFormat(curr.format)

	s := C.SDL_CreateRGBSurface(0, C.int(w), C.int(h),
		C.int(format.BitsPerPixel),
		format.Rmask, format.Gmask,
		format.Bmask, format.Amask)
	if s == nil {
		return nil, sdlerr()
	}

	surf := &Surface{s}
	runtime.SetFinalizer(surf, freesurf)
	return surf, nil
}

func freesurf(s *Surface) { C.SDL_FreeSurface(s.surf) }

type Texture struct {
	tex *C.SDL_Texture
}

func freetex(t *Texture) { C.SDL_DestroyTexture(t.tex) }

type Rect struct {
	X, Y int
	W, H int
}

func sdlRect(r *Rect) *C.SDL_Rect {
	return C.make_rect(C.int(r.X), C.int(r.Y), C.int(r.W), C.int(r.H))
}

func sdlerr() error { return errors.New(C.GoString(C.SDL_GetError())) }
