/*
 * XML DOM
 * Copyright 2024 John Douglas Pritchard, Syntelos
 */
package xdom

import (
	"errors"
	"fmt"
	"os"
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
	Read(string, []byte) (Node, error)
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
 * Node type includes parser states, e.g. element as "open",
 * "solitary" and "close", and text as "text" and "data".
 */
type Kind uint8
const (
	KindUndefined   Kind = 0
	KindDocument    Kind = 0b00000001
	KindAttribute   Kind = 0b00000010
	KindDeclaration Kind = 0b00000100
	KindInstruction Kind = 0b00001000
	KindOpen        Kind = 0b00010000
	KindSolitary    Kind = 0b00010001
	KindClose       Kind = 0b00010010
	KindText        Kind = 0b00100000
	KindData        Kind = 0b00100001
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
func (this Document) ReadFile (src *os.File) (n Node, er error){
	var fi os.FileInfo
	fi, er = src.Stat()
	if nil != er {
		return this, fmt.Errorf("Read error file not found: %w",er)
	} else {
		var sz int64 = fi.Size()
		var content []byte = make([]byte,sz)
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
func (this Document) Read (url string, content []byte) (n Node, er error){
	this.source = url
	this.content = content
	{
		var x, z int = 0, len(content)
		for x < z {
			var first, last int = x, this.content.read(x)
			if first < last {
				var begin, end int = first, (last+1)
				var text Text = this.content[begin:end]

				if KindUndefined != text.KindOf() {

					this.children = append(this.children,text)
				}
				x = end
			} else {
				x += 1
			}
		}
	}
	return this, nil
}
func (this Element) KindOf() (Kind){
	if nil != this.content {
		return this.content.KindOf()
	} else {
		return KindUndefined
	}
}
func (this Element) Content() (Text){
	return this.content
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
		var str []byte = make([]byte,depth)
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
func (this Element) Read(url string, content []byte) (n Node, er error) {
	this.content = content

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
			w = this.content.class(x,z,ws)
			if 0 < w {
				x = (w+1)

				y = this.content.class(x,z,at)
				if 0 < y {
					var at_be, at_en int = x, (y+1)
					var atx Text = this.content[at_be:at_en]

					var at Attribute
					n, er = at.Read(url,atx)
					if nil != er {
						return this, er
					} else {
						this.attributes = append(this.attributes,n.(Attribute))
						x = at_en
					}
				} else {
					break
				}
			} else {
				break
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
func (this Attribute) Read(url string, content []byte) (n Node, er error) {
	this.content = content

	var x, z int = 0, len(this.content)
	if x < z {
		var y int = (z-1)
		if '"' == content[x] || '"' == content[y] {

			if '"' == content[x] && '"' == content[y] {
				this.value = string(content)
			} else {
				return this, errors.New("Attribute quote missing.")
			}
		} else {
			y = this.content.class(x,z,id)
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
				y = this.class(x,z,ws)
				if (y+1) == z {

					return KindUndefined
				}
			}
		}
	}
	return KindText
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
func (this Text) Read(url string, content []byte) (n Node, er error) {
	this = content

	return this,nil
}
/*
 * Text span operator.
 */
func (this Text) identifier(ofs int) (string) {
	var x, z int = ofs, len(this)
	for ; x < z; x++ {

		switch this[x] {
		case ' ', '>', '[', ']', '&', '\r', '\n':

			return string(this[ofs:x])
		}
	}
	return string(this)
}

type cc func (byte)(bool)

func at (ch byte) (bool) {
        switch {
        case 'a' <= ch && 'z' >= ch :
                return true
        case 'A' <= ch && 'Z' >= ch :
                return true
        case '0' <= ch && '9' >= ch :
                return true
        case '_' == ch || '-' == ch || '+' == ch || '.' == ch || '=' == ch || ':' == ch :
                return true
        case '/' == ch || '\'' == ch || '"' == ch:
                return true
	case '?' == ch || '%' == ch || '!' == ch || '#' == ch || '$' == ch :
		return true
	case '(' == ch || ')' == ch || '[' == ch || ']' == ch || '*' == ch :
		return true
        default:
                return false
        }
}
func id (ch byte) (bool) {
        switch {
        case 'a' <= ch && 'z' >= ch :
                return true
        case 'A' <= ch && 'Z' >= ch :
                return true
        case '0' <= ch && '9' >= ch :
                return true
        case '_' == ch || '-' == ch || '+' == ch || '.' == ch || ':' == ch :
                return true
        default:
                return false
        }
}
func xc (ch byte) (bool) {
        switch {
	case '<' == ch || '>' == ch:
		return true
	case '?' == ch || '!' == ch:
		return true
        default:
                return false
        }
}
func ws (ch byte) (bool) {
        switch {
        case '\r' == ch || '\n' == ch || '\t' == ch || ' ' == ch:
                return true
        default:
                return false
        }
}
func (this Text) class ( ofs, len int, clop cc) (spx int) {
	/*
	 * Clamp to relationship
	 */
        spx = -1

        for ; ofs < len; ofs++ {

                if clop(this[ofs]) {

                        spx = ofs
                } else {
                        return spx
                }
        }
        return spx
}
func (this Text) scan(ofs, len int, ch byte) (spx int) {
	/*
	 * Clamp to first
	 */
	spx = ofs

	for ; ofs < len; ofs++ {

		if ch == this[ofs] {
			/*
			 * Found object next
			 */
			return ofs
		}
	}
	return spx
}
func (this Text) read(x int) (int) {
	var z = len(this)
	if x < z {
		/*
		 * Clamp to last
		 */
		var y int = (z-1)
		if x < y {
			if '<' == this[x] {
				/*
				 * Span code
				 */
				return this.scan(x,z,'>')

			} else {
				/*
				 * Span text
				 */
				return this.scan(x,z,'<')-1
			}
			return y
		}
	}
	return -1
}
