// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package gosym implements access to the Go symbol
// and line number tables embedded in Go binaries generated
// by the gc compilers.

// Package gosym implements access to the Go symbol
// and line number tables embedded in Go binaries generated
// by the gc compilers.
package gosym

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "strconv"
    "strings"
    "sync"
)

// DecodingError represents an error during the decoding of
// the symbol table.
type DecodingError struct {
	off int
	msg string
	val interface{}
}


// A Func collects information about a single function.
type Func struct {
	Entry uint64

	End       uint64
	Params    []*Sym
	Locals    []*Sym
	FrameSize int
	LineTable *LineTable
	Obj       *Obj
}


// A LineTable is a data structure mapping program counters to line numbers.
//
// In Go 1.1 and earlier, each function (represented by a Func) had its own
// LineTable, and the line number corresponded to a numbering of all source
// lines in the program, across all files. That absolute line number would then
// have to be converted separately to a file name and line number within the
// file.
//
// In Go 1.2, the format of the data changed so that there is a single LineTable
// for the entire program, shared by all Funcs, and there are no absolute line
// numbers, just line numbers within specific files.
//
// For the most part, LineTable's methods should be treated as an internal
// detail of the package; callers should use the methods on Table instead.

// A LineTable is a data structure mapping program counters to line numbers.
//
// In Go 1.1 and earlier, each function (represented by a Func) had its own
// LineTable, and the line number corresponded to a numbering of all source
// lines in the program, across all files. That absolute line number would then
// have to be converted separately to a file name and line number within the
// file.
//
// In Go 1.2, the format of the data changed so that there is a single LineTable
// for the entire program, shared by all Funcs, and there are no absolute line
// numbers, just line numbers within specific files.
//
// For the most part, LineTable's methods should be treated as an internal
// detail of the package; callers should use the methods on Table instead.
type LineTable struct {
	Data []byte
	PC   uint64
	Line int

	// Go 1.2 state
	mu       sync.Mutex
	go12     int // is this in Go 1.2 format? -1 no, 0 unknown, 1 yes
	binary   binary.ByteOrder
	quantum  uint32
	ptrsize  uint32
	functab  []byte
	nfunctab uint32
	filetab  []byte
	nfiletab uint32
	fileMap  map[string]uint32
}


// An Obj represents a collection of functions in a symbol table.
//
// The exact method of division of a binary into separate Objs is an internal
// detail of the symbol table format.
//
// In early versions of Go each source file became a different Obj.
//
// In Go 1 and Go 1.1, each package produced one Obj for all Go sources and one
// Obj per C source file.
//
// In Go 1.2, there is a single Obj for the entire program.

// An Obj represents a collection of functions in a symbol table.
//
// The exact method of division of a binary into separate Objs is an internal
// detail of the symbol table format.
//
// In early versions of Go each source file became a different Obj.
//
// In Go 1 and Go 1.1, each package produced one Obj for all Go sources and one
// Obj per C source file.
//
// In Go 1.2, there is a single Obj for the entire program.
type Obj struct {
	// Funcs is a list of functions in the Obj.
	Funcs []Func

	// In Go 1.1 and earlier, Paths is a list of symbols corresponding
	// to the source file names that produced the Obj.
	// In Go 1.2, Paths is nil.
	// Use the keys of Table.Files to obtain a list of source files.
	Paths []Sym // meta
}


// A Sym represents a single symbol table entry.
type Sym struct {
	Value  uint64
	Type   byte
	Name   string
	GoType uint64
	// If this symbol is a function symbol, the corresponding Func
	Func *Func
}


// Table represents a Go symbol table.  It stores all of the
// symbols decoded from the program and provides methods to translate
// between symbols, names, and addresses.

// Table represents a Go symbol table. It stores all of the
// symbols decoded from the program and provides methods to translate
// between symbols, names, and addresses.
type Table struct {
	Syms  []Sym
	Funcs []Func
	Files map[string]*Obj // nil for Go 1.2 and later binaries
	Objs  []Obj           // nil for Go 1.2 and later binaries

	go12line *LineTable // Go 1.2 line number table
}


// UnknownFileError represents a failure to find the specific file in
// the symbol table.
type UnknownFileError string


// UnknownLineError represents a failure to map a line to a program
// counter, either because the line is beyond the bounds of the file
// or because there is no code on the given line.
type UnknownLineError struct {
	File string
	Line int
}


// NewLineTable returns a new PC/line table
// corresponding to the encoded data.
// Text must be the start address of the
// corresponding text segment.
func NewLineTable(data []byte, text uint64) *LineTable

// NewTable decodes the Go symbol table in data,
// returning an in-memory representation.
func NewTable(symtab []byte, pcln *LineTable) (*Table, error)

func (*DecodingError) Error() string

// LineToPC returns the program counter for the given line number,
// considering only program counters before maxpc.
// Callers should use Table's LineToPC method instead.
func (*LineTable) LineToPC(line int, maxpc uint64) uint64

// PCToLine returns the line number for the given program counter.
// Callers should use Table's PCToLine method instead.
func (*LineTable) PCToLine(pc uint64) int

// BaseName returns the symbol name without the package or receiver name.
func (*Sym) BaseName() string

// PackageName returns the package part of the symbol name,
// or the empty string if there is none.
func (*Sym) PackageName() string

// ReceiverName returns the receiver type name of this symbol,
// or the empty string if there is none.
func (*Sym) ReceiverName() string

// Static reports whether this symbol is static (not visible outside its file).
func (*Sym) Static() bool

// LineToPC looks up the first program counter on the given line in
// the named file.  It returns UnknownPathError or UnknownLineError if
// there is an error looking up this line.

// LineToPC looks up the first program counter on the given line in
// the named file. It returns UnknownPathError or UnknownLineError if
// there is an error looking up this line.
func (*Table) LineToPC(file string, line int) (pc uint64, fn *Func, err error)

// LookupFunc returns the text, data, or bss symbol with the given name,
// or nil if no such symbol is found.
func (*Table) LookupFunc(name string) *Func

// LookupSym returns the text, data, or bss symbol with the given name,
// or nil if no such symbol is found.
func (*Table) LookupSym(name string) *Sym

// PCToFunc returns the function containing the program counter pc,
// or nil if there is no such function.
func (*Table) PCToFunc(pc uint64) *Func

// PCToLine looks up line number information for a program counter.
// If there is no information, it returns fn == nil.
func (*Table) PCToLine(pc uint64) (file string, line int, fn *Func)

// SymByAddr returns the text, data, or bss symbol starting at the given
// address.
func (*Table) SymByAddr(addr uint64) *Sym

func (*UnknownLineError) Error() string

func (UnknownFileError) Error() string

