/*
 * XML DOM
 * Copyright 2024 John Douglas Pritchard, Syntelos
 */
package xdom

import (
	"os"
	"testing"
)

const tst_text_svg string = "tst/text.svg"

func TestDocument(t *testing.T){
	var doc Document
	var fil *os.File
	var er error
	fil, er = os.Open(tst_text_svg)
	if nil != er {
		t.Fatalf("Opening '%s': %v",tst_text_svg,er)
	} else {
		defer fil.Close()

		doc, er = doc.ReadFile(fil)
		if nil != er {
			t.Fatalf("Reading '%s': %v",tst_text_svg,er)
		} else {
			doc.Print()
		}
	}
}
