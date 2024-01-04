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
  }
  /*
   * Read file into document.
   */
  func (Document) Read(*os.File) (Document, error)


Usage

   var doc Document
   var fil *os.File
   var er error
   fil, er = os.Open(filename)
   if nil != er {
	   t.Fatalf("Opening '%s': %v",filename,er)
   } else {
	   defer fil.Close()

	   doc, er = doc.Read(fil)
	   if nil != er {
		   t.Fatalf("Reading '%s': %v",filename,er)
	   } else {
		   doc.Print()
	   }
   }


References

  [XML] https://www.w3.org/TR/REC-xml/REC-xml-20081126.xml

