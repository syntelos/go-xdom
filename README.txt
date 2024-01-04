XML DOM

  /*
   * Principal user interface.
   */
  type Node interface {

	  KindOf() (Kind)
	  Content() (Text)
	  String() (string)
	  Print()
	  Depth() (uint8)
	  Append(Node) (Node)
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

Usage

   var doc Document
   var fil *os.File
   var er error
   fil, er = os.Open(filename)
   if nil != er {
	   t.Fatalf("Opening '%s': %v",filename,er)
   } else {
	   defer fil.Close()

	   doc, er = doc.ReadFile(fil)
	   if nil != er {
		   t.Fatalf("Reading '%s': %v",filename,er)
	   } else {
		   doc.Print()
	   }
   }


References

  [XML] https://www.w3.org/TR/REC-xml/REC-xml-20081126.xml

