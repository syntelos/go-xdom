XML DOM for GOPL

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
  type NodeList interface {

	  CountChildren() (uint32)
	  GetChild(uint32) (Node)
  }
  type AttributeList interface {

	  CountAttributes() (uint32)
	  GetAttribute(uint32) (Node)
  }
  /*
   * Read document.
   */
  func (Document) ReadFile (*os.File) (Document, error)
  func (Document) Read (string, []byte) (Document, error)


References

  [XML] https://www.w3.org/TR/REC-xml/REC-xml-20081126.xml

