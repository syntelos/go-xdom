/*
 * XML DOM for GOPL
 * Copyright 2024 John Douglas Pritchard, Syntelos
 */
package xdom

import (
	"fmt"
	"os"
	"testing"
)

func TestDocument(t *testing.T){
	var fn string = "tst/text.svg"
	var fil *os.File
	var er error
	fil, er = os.Open(fn)
	if nil != er {
		t.Fatalf("Opening '%s': %v",fn,er)
	} else {
		defer fil.Close()

		var doc Document
		var n Node

		n, er = doc.ReadFile(fil)
		if nil != er {
			t.Fatalf("Reading '%s': %v",fn,er)
		} else {
			doc = n.(Document)
			n.Print()
		}
		fmt.Println()
	}
}
