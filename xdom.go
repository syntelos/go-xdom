/*
 * XML DOM for GOPL
 * Copyright 2024 John Douglas Pritchard, Syntelos
 */
package xdom

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	span "github.com/syntelos/go-span"
	"strings"
)
/*
 * Principal user interface.
 */
type Node interface {

	KindOf() (Kind)
	Content() (Text)
	String() (string)
	Print()
	Depth() (uint8)
	Read(string, Text) (Node, error)
}
/*
 */
type NodeList interface {

	CountChildren() (uint32)
	GetChild(uint32) (Node)
}
/*
 */
type AttributeList interface {

	CountAttributes() (uint32)
	GetAttribute(uint32) (Node)
}
/*
 * Node type includes parser states, e.g. "code" and "text",
 * element as "open", "solitary" and "close".
 */
type Kind uint8
const (
	KindUndefined   Kind = 0
	KindCode        Kind = 0b10000000
	KindDocument    Kind = 0b10000001
	KindAttribute   Kind = 0b10000010
	KindDeclaration Kind = 0b10000100
	KindInstruction Kind = 0b10000101
	KindOpen        Kind = 0b10001000
	KindSolitary    Kind = 0b10001001
	KindClose       Kind = 0b10001010
	KindText        Kind = 0b00010000
	KindData        Kind = 0b00010001
)

type Document struct {
	source string
	content Text
	children []Node
}

type Element struct {
	parent Node
	content Text
	name string
	attributes []Attribute
	children []Node
}

type Attribute struct {
	content Text
	name string
	value string
}

type Text []byte

type TextList []Text

func (this Kind) IsCode() (bool) {
	return (KindCode == (this & KindCode))
}
func (this Kind) IsText() (bool) {
	return (KindText == (this & KindText))
}
func (this Kind) IsHead() (bool) {
	return (KindDeclaration == (this & KindDeclaration))
}
func (this Kind) IsBody() (bool) {
	return (KindOpen == this & KindOpen)
}
func (this Kind) IsOpen() (bool) {
	return (KindOpen == this)
}
func (this Kind) String() (string) {
	switch this {

	case KindDeclaration:
		return "<DECL>"
	case KindInstruction:
		return "<INSTR>"
	case KindDocument:
		return "<DOC>"
	case KindOpen:
		return "<OPEN>"
	case KindSolitary:
		return "<SOL>"
	case KindClose:
		return "<CLOSE>"
	case KindText:
		return "<TEXT>"
	case KindData:
		return "<DATA>"

	default:
		return "<UNKN>"
	}
}
func (this Document) KindOf() (Kind){
	return KindDocument
}
func (this Document) Content() (Text){
	return this.content
}
func (this Document) Source() (string){
	return this.source
}
func (this Document) String() (string){
	return this.content.String()
}
func (this Document) Print() {

	for index, node := range this.children {
		var kind Kind = node.KindOf()
		switch kind {
		case KindDeclaration, KindInstruction, KindDocument, KindOpen, KindSolitary, KindClose, KindText, KindData:
			fmt.Printf("%03o\t%s\t%s\n",index,kind,node)

			node.Print()
		}
	}
}
func (this Document) Depth() (uint8) {

	return 0
}
func (this Document) CountChildren() (index uint32) {
	return uint32(len(this.children))
}
func (this Document) GetChild(index uint32) (Node) {
	if index < this.CountChildren() {

		return this.children[index]
	} else {
		return nil
	}
}
/*
 * XML document parser.
 */
func (this Document) ReadFile (src *os.File) (n Node, er error){
	var fi os.FileInfo
	fi, er = src.Stat()
	if nil != er {
		return this, fmt.Errorf("Read error file not found: %w",er)
	} else {
		var sz int64 = fi.Size()
		var content Text = make([]byte,sz)
		var ct int
		ct, er = src.Read(content)

		if nil != er {
			return this, fmt.Errorf("ReadFile error '%s': %w",fi.Name(),er)
		} else if int64(ct) != sz {
			return this, fmt.Errorf("ReadFile error '%s': expected (%d) found (%d).",fi.Name(),sz,ct)
		} else {
			var url string = "file:"+src.Name()

			return this.Read(url,content)
		}
	}
}
/*
 * XML document parser.
 */
func (this Document) Read (url string, content Text) (n Node, er error){
	this.source = url
	this.content = content
	{
		var source TextList
		var kind Kind
		var text Text
		var body []byte

		n, er = source.Read(url,content)

		if nil == er {
			source = n.(TextList)

			for _, text = range source {
				kind = text.KindOf()
				if kind.IsCode() {

					if kind.IsBody() {
						/*
						 * Document body
						 */
						body = span.Cat(body,text)
					} else {
						/*
						 * Document head
						 */
						var el Element
						n, er = el.Read(url,text)
						if nil != er {
							return this, er
						} else {
							el = n.(Element)

							this.children = append(this.children,el)
						}
					}
				} else if kind.IsText() {

					if 0 == len(body) {
						this.children = append(this.children,text)
					} else {
						body = span.Cat(body,text)
					}
				}
			}
		}
		/*
		 * Document body
		 */
		var el Element
		n, er = el.Read(url,body)
		if nil != er {
			return this, er
		} else {
			el = n.(Element)

			this.children = append(this.children,el)
		}
	}
	return this, nil
}
func (this Element) KindOf() (Kind){
	if 0 != len(this.content) {
		return this.content.KindOf()
	} else {
		return KindUndefined
	}
}
func (this Element) Parent() (Node){
	return this.parent
}
func (this Element) Content() (Text){
	return this.content
}
func (this Element) Name() (string){
	return this.name
}
func (this Element) String() (string){
	var str strings.Builder
	{
		str.WriteString(this.name)
		str.WriteByte(' ')
		for ix, at := range this.attributes {
			if 0 != ix {
				str.WriteByte(' ')
			}
			str.WriteString(at.name)
		}
	}
	return str.String()
}
func (this Element) Print() {
	var indent string
	{
		var depth uint8 = this.Depth()
		var str Text = make([]byte,depth)
		var ix uint8
		for ix = 0; ix < depth; ix++ {
			str[ix] = '\t'
		}
		indent = string(str)
	}

	for index, node := range this.children {
		var kind Kind = node.KindOf()
		switch kind {
		case KindDeclaration, KindInstruction, KindOpen, KindSolitary:
			fmt.Printf("%s%03o\t%s\t%s\n",indent,index,kind,node)

			node.Print()
		}
	}
}
func (this Element) Depth() (uint8) {
	var p Node = this.parent
	var c uint8 = 1
	for nil != p {
		if KindOpen == p.KindOf() {
			c += 1
			{
				var e Element = p.(Element)
			
				p = e.parent
			}
		} else {
			c += 1
			break
		}
	}
	return c
}
func (this Element) Read(url string, content Text) (n Node, er error) {
	this.content = content
	/*
	 * Element attributes
	 */
	var kind Kind = this.KindOf()
	var w, x, y, z int = 0, 0, 0, len(content)

	switch kind {
	case KindDeclaration:
		x = 2
	case KindInstruction:
		x = 2
	case KindOpen:
		x = 1
	case KindSolitary:
		x = 1
	case KindClose:
		x = 2
	}
	this.name = this.content.identifier(x)

	switch kind {
	case KindDeclaration, KindInstruction, KindOpen, KindSolitary:
		x += len(this.name)

		for x < z {
			w = span.Class(this.content,x,z,span.WS)
			if 0 < w {
				x = (w+1)

				if '"' == this.content[x] {
					y = span.Forward(this.content,x,z,'"','"')
					if 0 < y {
						var at_be, at_en int = x, (y+1)
						var atx Text = this.content[at_be:at_en]

						var at Attribute
						n, er = at.Read(url,atx)
						if nil != er {
							return this, er
						} else {
							at = n.(Attribute)

							this.attributes = append(this.attributes,at)
							x = at_en
						}
					} else {
						break
					}
				} else if '\'' == this.content[x] {
					y = span.Forward(this.content,x,z,'\'','\'')
					if 0 < y {
						var at_be, at_en int = x, (y+1)
						var atx Text = this.content[at_be:at_en]

						var at Attribute
						n, er = at.Read(url,atx)
						if nil != er {
							return this, er
						} else {
							at = n.(Attribute)

							this.attributes = append(this.attributes,at)
							x = at_en
						}
					} else {
						break
					}
				} else {
					y = span.Class(this.content,x,z,span.XI)
					if 0 < y {
						if '=' == this.content[y] {
							y += 1
							if '"' == this.content[y] {
								y = span.Forward(this.content,y,z,'"','"')
								if 0 < y {
									var at_be, at_en int = x, (y+1)
									var atx Text = this.content[at_be:at_en]

									var at Attribute
									n, er = at.Read(url,atx)
									if nil != er {
										return this, er
									} else {
										at = n.(Attribute)

										this.attributes = append(this.attributes,at)
										x = at_en
									}
								} else {
									break
								}
							} else if '\'' == this.content[y] {
								y = span.Forward(this.content,y,z,'\'','\'')
								if 0 < y {
									var at_be, at_en int = x, (y+1)
									var atx Text = this.content[at_be:at_en]

									var at Attribute
									n, er = at.Read(url,atx)
									if nil != er {
										return this, er
									} else {
										at = n.(Attribute)

										this.attributes = append(this.attributes,at)
										x = at_en
									}
								} else {
									break
								}
							} else {
								y = span.Class(this.content,y,z,span.XA)
								if 0 < y {
									var at_be, at_en int = x, (y+1)
									var atx Text = this.content[at_be:at_en]

									var at Attribute
									n, er = at.Read(url,atx)
									if nil != er {
										return this, er
									} else {
										at = n.(Attribute)

										this.attributes = append(this.attributes,at)
										x = at_en
									}
								} else {
									break
								}
							}
						} else {
							var at_be, at_en int = x, (y+1)
							var atx Text = this.content[at_be:at_en]

							var at Attribute
							n, er = at.Read(url,atx)
							if nil != er {
								return this, er
							} else {
								at = n.(Attribute)

								this.attributes = append(this.attributes,at)
								x = at_en
							}
						}
					} else {
						break
					}
				}
			} else {
				break
			}
		}
	}
	/*
	 * Element content [TODO] (review)
	 */
	if KindOpen == kind {

		w = span.First(content,x,z,'>')
		if x < w && w < z {

			y = span.Last(content,(z-1),z,'<')
			if x < y && y < z {

				var text []byte = content[w:y]
				var list TextList
				var n Node
				n, er = list.Read(url,text)
				if nil != er {
					return this, fmt.Errorf("Parsing '%s': %w", text, er)
				} else {
					list = n.(TextList)

					var text Text
					var stack int = 0
					var body []byte

					for _, text = range list {

						switch text.KindOf() {
						case KindOpen:
							body = span.Cat(body,text)
							stack += 1
						case KindSolitary:
							if 0 == stack {
								var el Element
								n, er = el.Read(url,text)
								if nil != er {
									return this, er
								} else {
									el = n.(Element)

									this.children = append(this.children,el)
								}
							} else {
								body = span.Cat(body,text)
							}
						case KindClose:
							body = span.Cat(body,text)
							stack -= 1
							if 0 == stack {
								var el Element
								n, er = el.Read(url,body)
								if nil != er {
									return this, er
								} else {
									el = n.(Element)

									this.children = append(this.children,el)

									body = make([]byte,0)
								}
							}
						case KindText, KindData:
							body = span.Cat(body,text)
						}
					}
				}
			}
		}
	}
	return this,nil
}
func (this Element) CountChildren() (index uint32) {
	return uint32(len(this.children))
}
func (this Element) GetChild(index uint32) (ch Node) {
	if index < this.CountChildren() {

		return this.children[index]
	} else {
		return ch
	}
}
func (this Element) CountAttributes() (index uint32) {
	return uint32(len(this.attributes))
}
func (this Element) GetAttribute(index uint32) (at Attribute) {
	if index < this.CountAttributes() {

		return this.attributes[index]
	} else {
		return at
	}
}
func (this Attribute) KindOf() (Kind) {
	return KindAttribute
}
func (this Attribute) Content() (Text) {
	return this.content
}
func (this Attribute) Name() (string){
	return this.name
}
func (this Attribute) Value() (string){
	return this.value
}
func (this Attribute) String() (string) {
	if "" != this.name {
		return this.name
	} else if "" != this.value {
		return this.value
	} else {
		return ""
	}
}
func (this Attribute) Print() {
}
func (this Attribute) Depth() (uint8) {
	return 0
}
func (this Attribute) Read(url string, content Text) (n Node, er error) {
	this.content = content

	var x, z int = 0, len(this.content)
	if x < z {
		var y int = (z-1)
		if '"' == content[x] {

			if '"' == content[x] && '"' == content[y] {
				this.value = string(content)
			} else {
				return this, fmt.Errorf("Attribute quote missing in '%s'.",content)
			}
		} else {
			y = span.Class(this.content,x,z,span.XI)
			if 0 < y {
				y += 1
				if y < z {
					if '=' == this.content[y] {
						this.name = string(this.content[x:y])
						this.value = string(this.content[y+1])
					} else {
						return this, fmt.Errorf("Attribute syntax of content '%s'.",content)
					}
				} else {
					this.name = string(content)
				}
			} else {
				this.value = string(content)
			}
		}
		return this, nil
	} else {
		return this, errors.New("Attribute empty.")
	}
}
func (this Text) KindOf() (Kind){
	var x int = 0
	var z int = len(this)
	if x < z {
		var y int = (z-1)
		if x < y {
			if '<' == this[x] && '>' == this[y] {
				x += 1
				y -= 1
				if x < y {
					switch this[x] {

					case '?':
						return KindInstruction
					case '!':
						x += 1
						if x < y {
							if '[' == this[x] {
								return KindData
							} else {
								return KindDeclaration
							}
						}
					case '/':
						return KindClose
					default:
						if '/' == this[y] {
							return KindSolitary
						} else {
							return KindOpen
						}
					}
				}
			} else {
				var w int = span.Class(this,x,z,span.WS)
				if w == y {

					return KindUndefined
				} else {
					return KindText
				}
			}
		}
	}
	return KindUndefined
}
func (this Text) Content() (Text){
	return this
}
func (this Text) String() (string) {
	var x, z = 0, len(this)
	if x < z {
		var y int = z
		/*
		 * Clamp to line
		 */
		for ; x < y; x++ {
			if '\n' == this[x] {
				y = x
				break
			} else if 20 == x {
				y = 20
				break
			}
		}
		/*
		 * Clamp to twenty
		 */
		if 20 < y {
			return string(this[0:20])
		} else if y < z {
			return string(this[0:y])
		} else {
			return string(this)
		}
	}
	return ""
}
func (this Text) Print() {
}
func (this Text) Depth() (uint8) {
	return 0
}
func (this Text) Read(url string, content Text) (n Node, er error) {
	this = content

	return this,nil
}
/*
 * Text span operator.
 */
func (this Text) identifier(ofs int) (string) {
	var x, z int = ofs, len(this)
	var y int = span.Class(this,x,z,span.XI)
	if x <= y && y <= z {
		return string(this[x:y+1])
	} else {
		return string(this)
	}
}

func (this TextList) KindOf() (Kind) {

	return this.Content().KindOf()
}
func (this TextList) Content() (Text) {
	var buf bytes.Buffer
	for _, txt := range this {
		buf.Write(txt)
	}
	return buf.Bytes()
}
func (this TextList) String() (string) {
	var content string = string(this.Content())
	var x, z = 0, len(content)
	if x < z {
		var y int = z
		/*
		 * Clamp to line
		 */
		for ; x < y; x++ {
			if '\n' == content[x] {
				y = x
				break
			} else if 20 == x {
				y = 20
				break
			}
		}
		/*
		 * Clamp to twenty
		 */
		if 20 < y {
			return string(content[0:20])
		} else if y < z {
			return string(content[0:y])
		} else {
			return string(content)
		}
	}
	return content
}
func (this TextList) Print() {

	for index, node := range this {

		fmt.Printf("%03o\t%s\n",index,node)
	}
}
func (this TextList) Depth() (uint8) {
	return 0
}
/*
 * XML stream parser.
 */
func (this TextList) ReadFile (src *os.File) (n Node, er error){
	var fi os.FileInfo
	fi, er = src.Stat()
	if nil != er {
		return this, fmt.Errorf("Read error file not found: %w",er)
	} else {
		var sz int64 = fi.Size()
		var content Text = make([]byte,sz)
		var ct int
		ct, er = src.Read(content)

		if nil != er {
			return this, fmt.Errorf("ReadFile error '%s': %w",fi.Name(),er)
		} else if int64(ct) != sz {
			return this, fmt.Errorf("ReadFile error '%s': expected (%d) found (%d).",fi.Name(),sz,ct)
		} else {
			var url string = "file:"+src.Name()

			return this.Read(url,content)
		}
	}
}
/*
 * XML stream parser.
 */
func (this TextList) Read (url string, content Text) (n Node, er error){
	var x, z int = 0, len(content)
	for x < z {
		var first, last int = x, span.Forward(content,x,z,'<','>')
		if first < last {
			var begin, end int = first, (last+1)
			var text Text = content[begin:end]
			if KindUndefined != text.KindOf() {

				this = append(this,text)
			}
			x = end
		} else {
			x += 1
		}
	}
	return this, nil
}
