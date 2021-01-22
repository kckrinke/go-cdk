// Copyright 2020 The CDK Authors
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

// Text-based version of Pango

/*
<span
  style=[normal,italic]
  weight=[dim,normal,bold]
  foreground=[color]
  background=[color]
  underline=[bool]
  strikethrough=[bool]
>
 CONTENT
</span>
<b></b>
<i></i>
<s></s>
<u></u>

*/

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	TabStops = 4
)

type Tango interface {
	Raw() string
	TextBuffer() TextBuffer
}

type CTango struct {
	raw    string
	style  Style
	marked []TextCell
	input  WordLine
}

func NewMarkup(text string, style Style) (markup Tango, err error) {
	if !strings.HasPrefix(text, "<markup") {
		text = "<markup>" + text + "</markup>"
	}
	m := &CTango{
		raw:   text,
		style: style,
	}
	if err = m.init(); err == nil {
		markup = m
	} else {
		markup = nil
	}
	return
}

func (m *CTango) Raw() string {
	return m.raw
}

func (m *CTango) TextBuffer() TextBuffer {
	tb := NewEmptyTextBuffer(m.style)
	tb.SetInput(m.input)
	return tb
}

func (m *CTango) init() error {
	m.marked = []TextCell{}
	m.input = NewEmptyWordLine()
	r := strings.NewReader(m.raw)
	parser := xml.NewDecoder(r)

	wid := 0

	mstyle := m.style
	cstyle := m.style
	pstyle := m.style

	isWord := false
	var err error
	var token xml.Token
	for {
		token, err = parser.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			pstyle = cstyle
			switch t.Name.Local {
			case "markup":
				pstyle = mstyle
				mstyle = m.parseStyleAttrs(t.Attr)
			case "span":
				pstyle = cstyle
				cstyle = m.parseStyleAttrs(t.Attr)
			case "b":
				cstyle = cstyle.Bold(true)
			case "i":
				cstyle = cstyle.Italic(true)
			case "s":
				cstyle = cstyle.StrikeThrough(true)
			case "u":
				cstyle = cstyle.Underline(true)
			case "d":
				cstyle = cstyle.Dim(true)
			}
		case xml.EndElement:
			switch t.Name.Local {
			case "markup":
				cstyle = pstyle
			case "span":
				cstyle = pstyle
			case "b":
				cstyle = cstyle.Bold(false)
			case "i":
				cstyle = cstyle.Italic(false)
			case "s":
				cstyle = cstyle.StrikeThrough(false)
			case "u":
				cstyle = cstyle.Underline(false)
			case "d":
				cstyle = cstyle.Dim(false)
			}
		case xml.CharData:
			content := xml.CharData(t) // CharData []byte
			for idx := 0; idx < len(content); idx++ {
				v, _ := utf8.DecodeRune(content[idx:])
				m.marked = append(m.marked, NewRuneCell(v, cstyle))
				if unicode.IsSpace(v) {
					if isWord {
						isWord = false
						m.input.AppendWordCell(NewEmptyWordCell())
						wid = m.input.Len() - 1
					}
				} else {
					if !isWord {
						isWord = true
						m.input.AppendWordCell(NewEmptyWordCell())
						wid = m.input.Len() - 1
					}
				}
				if wid >= m.input.Len() {
					m.input.AppendWordCell(NewEmptyWordCell())
				}
				m.input.AppendWordRune(wid, v, cstyle)
			} // for idx len(content)
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		default:
			fmt.Println("Unknown")
		}
	}
	return nil
}

func (m *CTango) parseStyleAttrs(attrs []xml.Attr) (style Style) {
	style = m.style
	for _, attr := range attrs {
		switch attr.Name.Local {
		case "style":
			switch attr.Value {
			case "normal":
				style = style.Italic(false)
			case "italic":
				style = style.Italic(true)
			}
		case "weight":
			switch attr.Value {
			case "dim":
				style = style.Dim(true)
			case "normal":
				style = style.Dim(false).Bold(false)
			case "bold":
				style = style.Bold(true)
			}
		case "foreground":
			style = style.Foreground(GetColor(attr.Value))
		case "background":
			style = style.Background(GetColor(attr.Value))
		case "underline":
			style = style.Underline(attr.Value == "true" || attr.Value == "1")
		case "strikethrough":
			style = style.StrikeThrough(attr.Value == "true" || attr.Value == "1")
		}
	}
	return
}
