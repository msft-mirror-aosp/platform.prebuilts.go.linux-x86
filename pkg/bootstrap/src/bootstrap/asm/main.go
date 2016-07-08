// Do not edit. Bootstrap copy of /usr/local/google/buildbot/src/android/build-tools/out/obj/go/src/cmd/asm/main.go

//line /usr/local/google/buildbot/src/android/build-tools/out/obj/go/src/cmd/asm/main.go:1
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"bootstrap/asm/internal/arch"
	"bootstrap/asm/internal/asm"
	"bootstrap/asm/internal/flags"
	"bootstrap/asm/internal/lex"

	"bootstrap/internal/bio"
	"bootstrap/internal/obj"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("asm: ")

	GOARCH := obj.Getgoarch()

	architecture := arch.Set(GOARCH)
	if architecture == nil {
		log.Fatalf("unrecognized architecture %s", GOARCH)
	}

	flags.Parse()

	ctxt := obj.Linknew(architecture.LinkArch)
	if *flags.PrintOut {
		ctxt.Debugasm = 1
	}
	ctxt.LineHist.TrimPathPrefix = *flags.TrimPath
	ctxt.Flag_dynlink = *flags.Dynlink
	ctxt.Flag_shared = *flags.Shared || *flags.Dynlink
	ctxt.Bso = bufio.NewWriter(os.Stdout)
	defer ctxt.Bso.Flush()

	// Create object file, write header.
	out, err := os.Create(*flags.OutputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer bio.MustClose(out)
	buf := bufio.NewWriter(bio.MustWriter(out))

	fmt.Fprintf(buf, "go object %s %s %s\n", obj.Getgoos(), obj.Getgoarch(), obj.Getgoversion())
	fmt.Fprintf(buf, "!\n")

	lexer := lex.NewLexer(flag.Arg(0), ctxt)
	parser := asm.NewParser(ctxt, architecture, lexer)
	diag := false
	ctxt.DiagFunc = func(format string, args ...interface{}) {
		diag = true
		log.Printf(format, args...)
	}
	pList := obj.Linknewplist(ctxt)
	var ok bool
	pList.Firstpc, ok = parser.Parse()
	if ok {
		// reports errors to parser.Errorf
		obj.Writeobjdirect(ctxt, buf)
	}
	if !ok || diag {
		log.Printf("assembly of %s failed", flag.Arg(0))
		os.Remove(*flags.OutputFile)
		os.Exit(1)
	}
	buf.Flush()
}
