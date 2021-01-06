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

type Tango struct {
	raw    string
	clean  string
	style  Style
	marked []TextCell
	input  *WordLine
}

func NewMarkup(text string, style Style) (markup *Tango, err error) {
	if !strings.HasPrefix(text, "<markup") {
		text = "<markup>" + text + "</markup>"
	}
	markup = &Tango{
		raw:   text,
		style: style,
	}
	err = markup.init()
	if err != nil {
		markup = nil
	}
	return
}

func (m *Tango) Raw() string {
	return m.raw
}

func (m *Tango) Clean() string {
	return m.clean
}

func (m *Tango) TextBuffer() TextBuffer {
	tb := NewEmptyTextBuffer(m.style)
	tb.SetInput(m.input)
	return tb
}

func (m *Tango) init() error {
	m.clean = ""
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
			if t.Name.Local == "markup" {
				cstyle = mstyle
			} else {
				cstyle = pstyle
			}
		case xml.CharData:
			content := xml.CharData(t) // CharData []byte
			for idx := 0; idx < len(content); idx++ {
				v, _ := utf8.DecodeRune(content[idx:])
				m.clean += string(v)
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
				m.input.words[wid].AppendRune(v, cstyle)
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

func (m *Tango) parseStyleAttrs(attrs []xml.Attr) (style Style) {
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
