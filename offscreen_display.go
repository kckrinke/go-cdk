// Copyright 2016 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cdk

import (
	"fmt"
	"sync"
	"unicode/utf8"

	"golang.org/x/text/transform"

	"github.com/kckrinke/go-cdk/utils"
)

const OffscreenDisplayTtyPath = "<offscreen>"

func MakeOffscreenDisplay(charset string) (OffscreenDisplay, error) {
	s := NewOffscreenDisplay(charset)
	if s == nil {
		return nil, fmt.Errorf("failed to get simulation display")
	}
	if e := s.Init(); e != nil {
		return nil, fmt.Errorf("failed to initialize display: %v", e)
	}
	return s, nil
}

// NewOffscreenDisplay returns a OffscreenDisplay.  Note that
// OffscreenDisplay is also a Display.
func NewOffscreenDisplay(charset string) OffscreenDisplay {
	if utils.IsEmpty(charset) {
		charset = GetCharset()
	}
	s := &COffscreenDisplay{charset: charset}
	return s
}

// OffscreenDisplay represents a display simulation.  This is intended to
// be a superset of normal Screens, but also adds some important interfaces
// for testing.
type OffscreenDisplay interface {
	// InjectKeyBytes injects a stream of bytes corresponding to
	// the native encoding (see charset).  It turns true if the entire
	// set of bytes were processed and delivered as KeyEvents, false
	// if any bytes were not fully understood.  Any bytes that are not
	// fully converted are discarded.
	InjectKeyBytes(buf []byte) bool

	// InjectKey injects a key event.  The rune is a UTF-8 rune, post
	// any translation.
	InjectKey(key Key, r rune, mod ModMask)

	// InjectMouse injects a mouse event.
	InjectMouse(x, y int, buttons ButtonMask, mod ModMask)

	// SetSize resizes the underlying physical display.  It also causes
	// a resize event to be injected during the next Show() or Sync().
	// A new physical contents array will be allocated (with data from
	// the old copied), so any prior value obtained with GetContents
	// won't be used anymore
	SetSize(width, height int)

	// GetContents returns display contents as an array of
	// cells, along with the physical width & height.   Note that the
	// physical contents will be used until the next time SetSize()
	// is called.
	GetContents() (cells []OffscreenCell, width int, height int)

	// GetCursor returns the cursor details.
	GetCursor() (x int, y int, visible bool)

	Display
}

// OffscreenCell represents a simulated display cell.  The purpose of this
// is to track on display content.
type OffscreenCell struct {
	// Bytes is the actual character bytes.  Normally this is
	// rune data, but it could be be data in another encoding system.
	Bytes []byte

	// Style is the style used to display the data.
	Style Style

	// Runes is the list of runes, unadulterated, in UTF-8.
	Runes []rune
}

type COffscreenDisplay struct {
	physw int
	physh int
	fini  bool
	style Style
	evch  chan Event
	quit  chan struct{}

	front     []OffscreenCell
	back      *CellBuffer
	clear     bool
	cursorx   int
	cursory   int
	cursorvis bool
	mouse     bool
	paste     bool
	charset   string
	encoder   transform.Transformer
	decoder   transform.Transformer
	fillchar  rune
	fillstyle Style
	fallback  map[rune]string

	sync.Mutex
}

func (s *COffscreenDisplay) Init() error {
	return s.InitWithTty("")
}

func (s *COffscreenDisplay) InitWithTty(_ string) error {
	s.evch = make(chan Event, 10)
	s.quit = make(chan struct{})
	s.fillchar = 'X'
	s.fillstyle = StyleDefault
	s.mouse = false
	s.physw = 80
	s.physh = 25
	s.cursorx = -1
	s.cursory = -1
	s.style = StyleDefault
	s.back = NewCellBuffer()

	if enc := GetEncoding(s.charset); enc != nil {
		s.encoder = enc.NewEncoder()
		s.decoder = enc.NewDecoder()
	} else {
		return ErrNoCharset
	}

	s.front = make([]OffscreenCell, s.physw*s.physh)
	s.back.Resize(80, 25)

	// default fallbacks
	s.fallback = make(map[rune]string)
	for k, v := range RuneFallbacks {
		s.fallback[k] = v
	}
	return nil
}

func (s *COffscreenDisplay) Close() {
	s.fini = true
	s.back.Resize(0, 0)
	if s.quit != nil {
		close(s.quit)
	}
	s.physw = 0
	s.physh = 0
	s.front = nil
}

func (s *COffscreenDisplay) SetStyle(style Style) {
	s.style = style
}

func (s *COffscreenDisplay) Clear() {
	s.Fill(' ', s.style)
}

func (s *COffscreenDisplay) Fill(r rune, style Style) {
	s.back.Fill(r, style)
}

func (s *COffscreenDisplay) SetCell(x, y int, style Style, ch ...rune) {
	if len(ch) > 0 {
		s.SetContent(x, y, ch[0], ch[1:], style)
	} else {
		s.SetContent(x, y, ' ', nil, style)
	}
}

func (s *COffscreenDisplay) SetContent(x, y int, mainc rune, combc []rune, st Style) {
	s.back.SetContent(x, y, mainc, combc, st)
}

func (s *COffscreenDisplay) GetContent(x, y int) (rune, []rune, Style, int) {
	var mainc rune
	var combc []rune
	var style Style
	var width int
	mainc, combc, style, width = s.back.GetContent(x, y)
	return mainc, combc, style, width
}

func (s *COffscreenDisplay) drawCell(x, y int) int {

	mainc, combc, style, width := s.back.GetContent(x, y)
	if !s.back.Dirty(x, y) {
		return width
	}
	if x >= s.physw || y >= s.physh || x < 0 || y < 0 {
		return width
	}
	simc := &s.front[(y*s.physw)+x]

	if style == StyleDefault {
		style = s.style
	}
	simc.Style = style
	simc.Runes = append([]rune{mainc}, combc...)

	// now emit runes - taking care to not overrun width with a
	// wide character, and to ensure that we emit exactly one regular
	// character followed up by any residual combing characters

	simc.Bytes = nil

	if x > s.physw-width {
		simc.Runes = []rune{' '}
		simc.Bytes = []byte{' '}
		return width
	}

	lbuf := make([]byte, 12)
	ubuf := make([]byte, 12)
	nout := 0

	for _, r := range simc.Runes {

		l := utf8.EncodeRune(ubuf, r)

		nout, _, _ = s.encoder.Transform(lbuf, ubuf[:l], true)

		if nout == 0 || lbuf[0] == '\x1a' {

			// skip combining

			if subst, ok := s.fallback[r]; ok {
				simc.Bytes = append(simc.Bytes,
					[]byte(subst)...)

			} else if r >= ' ' && r <= '~' {
				simc.Bytes = append(simc.Bytes, byte(r))

			} else if simc.Bytes == nil {
				simc.Bytes = append(simc.Bytes, '?')
			}
		} else {
			simc.Bytes = append(simc.Bytes, lbuf[:nout]...)
		}
	}
	s.back.SetDirty(x, y, false)
	return width
}

func (s *COffscreenDisplay) ShowCursor(x, y int) {
	s.cursorx, s.cursory = x, y
	s.showCursor()
}

func (s *COffscreenDisplay) HideCursor() {
	s.ShowCursor(-1, -1)
}

func (s *COffscreenDisplay) showCursor() {

	x, y := s.cursorx, s.cursory
	if x < 0 || y < 0 || x >= s.physw || y >= s.physh {
		s.cursorvis = false
	} else {
		s.cursorvis = true
	}
}

func (s *COffscreenDisplay) hideCursor() {
	// does not update cursor position
	s.cursorvis = false
}

func (s *COffscreenDisplay) Show() {
	s.resize()
	s.draw()
}

func (s *COffscreenDisplay) clearScreen() {
	// We emulate a hardware clear by filling with a specific pattern
	for i := range s.front {
		s.front[i].Style = s.fillstyle
		s.front[i].Runes = []rune{s.fillchar}
		s.front[i].Bytes = []byte{byte(s.fillchar)}
	}
	s.clear = false
}

func (s *COffscreenDisplay) draw() {
	s.hideCursor()
	if s.clear {
		s.clearScreen()
	}

	w, h := s.back.Size()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			width := s.drawCell(x, y)
			x += width - 1
		}
	}
	s.showCursor()
}

func (s *COffscreenDisplay) EnableMouse() {
	s.mouse = true
}

func (s *COffscreenDisplay) DisableMouse() {
	s.mouse = false
}

func (s *COffscreenDisplay) EnablePaste() {
	s.paste = true
}

func (s *COffscreenDisplay) DisablePaste() {
	s.paste = false
}

func (s *COffscreenDisplay) Size() (w, h int) {
	w, h = s.back.Size()
	return
}

func (s *COffscreenDisplay) resize() {
	w, h := s.physw, s.physh
	ow, oh := s.back.Size()
	if w != ow || h != oh {
		s.back.Resize(w, h)
		ev := NewEventResize(w, h)
		s.PostEvent(ev)
	}
}

func (s *COffscreenDisplay) Colors() int {
	return 256
}

func (s *COffscreenDisplay) PollEvent() Event {
	select {
	case <-s.quit:
		return nil
	case ev := <-s.evch:
		return ev
	}
}

func (s *COffscreenDisplay) PostEventWait(ev Event) {
	s.evch <- ev
}

func (s *COffscreenDisplay) PostEvent(ev Event) error {
	select {
	case s.evch <- ev:
		return nil
	default:
		return ErrEventQFull
	}
}

func (s *COffscreenDisplay) InjectMouse(x, y int, buttons ButtonMask, mod ModMask) {
	ev := NewEventMouse(x, y, buttons, mod)
	s.PostEvent(ev)
}

func (s *COffscreenDisplay) InjectKey(key Key, r rune, mod ModMask) {
	ev := NewEventKey(key, r, mod)
	s.PostEvent(ev)
}

func (s *COffscreenDisplay) InjectKeyBytes(b []byte) bool {
	failed := false

outer:
	for len(b) > 0 {
		if b[0] >= ' ' && b[0] <= 0x7F {
			// printable ASCII easy to deal with -- no encodings
			ev := NewEventKey(KeyRune, rune(b[0]), ModNone)
			s.PostEvent(ev)
			b = b[1:]
			continue
		}

		if b[0] < 0x80 {
			mod := ModNone
			// No encodings start with low numbered values
			if Key(b[0]) >= KeyCtrlA && Key(b[0]) <= KeyCtrlZ {
				mod = ModCtrl
			}
			ev := NewEventKey(Key(b[0]), 0, mod)
			s.PostEvent(ev)
			b = b[1:]
			continue
		}

		utfb := make([]byte, len(b)*4) // worst case
		for l := 1; l < len(b); l++ {
			s.decoder.Reset()
			nout, nin, _ := s.decoder.Transform(utfb, b[:l], true)

			if nout != 0 {
				r, _ := utf8.DecodeRune(utfb[:nout])
				if r != utf8.RuneError {
					ev := NewEventKey(KeyRune, r, ModNone)
					s.PostEvent(ev)
				}
				b = b[nin:]
				continue outer
			}
		}
		failed = true
		b = b[1:]
		continue
	}

	return !failed
}

func (s *COffscreenDisplay) Sync() {
	s.clear = true
	s.resize()
	s.back.Invalidate()
	s.draw()
}

func (s *COffscreenDisplay) CharacterSet() string {
	return s.charset
}

func (s *COffscreenDisplay) SetSize(w, h int) {
	newc := make([]OffscreenCell, w*h)
	for row := 0; row < h && row < s.physh; row++ {
		for col := 0; col < w && col < s.physw; col++ {
			newc[(row*w)+col] = s.front[(row*s.physw)+col]
		}
	}
	s.cursorx, s.cursory = -1, -1
	s.physw, s.physh = w, h
	s.front = newc
	s.back.Resize(w, h)
}

func (s *COffscreenDisplay) GetContents() ([]OffscreenCell, int, int) {
	cells, w, h := s.front, s.physw, s.physh
	return cells, w, h
}

func (s *COffscreenDisplay) GetCursor() (int, int, bool) {
	x, y, vis := s.cursorx, s.cursory, s.cursorvis
	return x, y, vis
}

func (s *COffscreenDisplay) RegisterRuneFallback(r rune, subst string) {
	s.fallback[r] = subst
}

func (s *COffscreenDisplay) UnregisterRuneFallback(r rune) {
	delete(s.fallback, r)
}

func (s *COffscreenDisplay) CanDisplay(r rune, checkFallbacks bool) bool {
	if enc := s.encoder; enc != nil {
		nb := make([]byte, 6)
		ob := make([]byte, 6)
		num := utf8.EncodeRune(ob, r)

		enc.Reset()
		dst, _, err := enc.Transform(nb, ob[:num], true)
		if dst != 0 && err == nil && nb[0] != '\x1A' {
			return true
		}
	}
	if !checkFallbacks {
		return false
	}
	if _, ok := s.fallback[r]; ok {
		return true
	}
	return false
}

func (s *COffscreenDisplay) HasMouse() bool {
	return false
}

func (s *COffscreenDisplay) Resize(int, int, int, int) {}

func (s *COffscreenDisplay) HasKey(Key) bool {
	return true
}

func (s *COffscreenDisplay) Beep() error {
	return nil
}

func (t *COffscreenDisplay) Export() *CellBuffer {
	t.Lock()
	defer t.Unlock()
	cb := NewCellBuffer()
	w, h := t.back.Size()
	cb.Resize(w, h)
	for idx, cell := range t.back.cells {
		cb.cells[idx] = cell
	}
	return cb
}

func (t *COffscreenDisplay) Import(cb *CellBuffer) {
	t.Lock()
	defer t.Unlock()
	w, h := cb.Size()
	t.back.Resize(w, h)
	for idx, cell := range cb.cells {
		t.back.cells[idx] = cell
	}
}
