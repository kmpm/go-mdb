package data

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/yageek/go-mdb/pages"
	"github.com/yageek/go-mdb/version"
)

//Offset
const (
	HeaderFreeSpaceOffset int = 0x02

	HeaderTableDefinitionPointerOffset = 0x04

	Jet3HeaderNumRowsOffset = 0x08
	Jet4HeaderNumRowsOffset = 0x0c
)

type PageKind byte

const (
	PageKindData  PageKind = 0x01
	PageKindTdef  PageKind = 0x02
	PageKindIndex PageKind = 0x03
)

func (pk PageKind) String() string {
	switch pk {
	case PageKindData:
		return "PageKindData"
	case PageKindTdef:
		return "PageKindTdef"
	case PageKindIndex:
		return "PageKindIndex"
	default:
		return fmt.Sprintf("PageKind %0X", byte(pk))
	}
}

type Jet3DatapageHeader struct {
	Kind    PageKind
	_       byte
	Space   uint16
	Pointer uint32
	Count   uint16
}

func (j *Jet3DatapageHeader) PageKind() PageKind  { return j.Kind }
func (j *Jet3DatapageHeader) FreeSpace() uint16   { return j.Space }
func (j *Jet3DatapageHeader) PagePointer() uint32 { return j.Pointer }
func (j *Jet3DatapageHeader) RowsCount() uint16   { return j.Count }
func (j *Jet3DatapageHeader) Size() int           { return 10 }

// Jet4DatapageHeader is TDEF Header for Jet4
type Jet4DatapageHeader struct {
	// page_type
	Kind PageKind
	_    byte
	// tdef_id
	Space uint16
	// next_pg
	Pointer uint32
	_       [4]byte
	Count   uint16
}

func (j *Jet4DatapageHeader) PageKind() PageKind  { return j.Kind }
func (j *Jet4DatapageHeader) FreeSpace() uint16   { return j.Space }
func (j *Jet4DatapageHeader) PagePointer() uint32 { return j.Pointer }
func (j *Jet4DatapageHeader) RowsCount() uint16   { return j.Count }
func (j *Jet4DatapageHeader) Size() int           { return 12 }

type DatapageHeader interface {
	PageKind() PageKind
	FreeSpace() uint16
	PagePointer() uint32
	RowsCount() uint16
	Size() int
}

// NewDataPageHeader returns a new datapage header from page
func NewDataPageHeader(page []byte, v version.JetVersion) (DatapageHeader, error) {

	buff := bytes.NewReader(page)
	var header DatapageHeader

	if v == version.Jet3 {
		header = new(Jet3DatapageHeader)
	} else {
		header = new(Jet4DatapageHeader)
	}
	err := binary.Read(buff, binary.LittleEndian, header)

	if err != nil {
		return nil, err
	}

	if byte(header.PageKind()) != v.MagicNumber() {
		return nil, pages.ErrInvalidPageCode
	}
	return header, nil
}
