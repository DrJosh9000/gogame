package game

// Sigh.

import (
	"sdl"
	"sort"
)

type windowManager struct {
	all     objectSlice
	current object
}

func (w *windowManager) add(b clicker) {
	w.all = append(w.all, b)
}

func (w *windowManager) life() {
	inbox := make(chan message, 10)
	kmp("input.event", inbox)
	for msg := range inbox {
		switch m := msg.v.(type) {
		case *sdl.MouseButtonDownEvent:
			sort.Sort(w.all)
			for i := len(w.all) - 1; i >= 0; i-- {
				b, ok := w.all[i].(clicker)
				if !ok {
					continue
				}
				if b.invisible() {
					continue
				}
				if b.bounds().Contains(m.X, m.Y) {
					b.setDown(true)
					w.current = b
				}
			}
		case *sdl.MouseButtonUpEvent:
			if w.current == nil {
				continue
			}
			b, ok := w.current.(clicker)
			if !ok {
				w.current = nil
				continue
			}
			if b.bounds().Contains(m.X, m.Y) {
				b.click()
				b.setDown(false)
			}
			w.current = nil
		case *sdl.MouseMotionEvent:
			if w.current == nil {
				continue
			}
			b, ok := w.current.(clicker)
			if !ok {
				w.current = nil
				continue
			}
			b.setDown(b.bounds().Contains(m.X, m.Y))
		}
	}
}
