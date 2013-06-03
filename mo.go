package gettext

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/textproto"
	"strings"
)

const (
	leMagic = 0x950412de
	beMagic = 0xde120495
)

var (
	ErrInvalidMagic = errors.New("magic has failed")
	ErrTruncated    = errors.New("truncated")
)

func ParseMO(r io.ReadSeeker) (*Catalog, error) {
	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return nil, err
	}

	var bo binary.ByteOrder
	if magic == leMagic {
		bo = binary.LittleEndian
	} else if magic == beMagic {
		bo = binary.BigEndian
	} else {
		return nil, ErrInvalidMagic
	}

	var version, stringCount, origOffset, transOffset uint32
	if err := binary.Read(r, bo, &version); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bo, &stringCount); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bo, &origOffset); err != nil {
		return nil, err
	}
	if err := binary.Read(r, bo, &transOffset); err != nil {
		return nil, err
	}

	// TODO: Check version

	stringOffsets := make([]struct{ origStart, origLen, transStart, transLen int32 }, stringCount)

	if o, err := r.Seek(int64(origOffset), 0); err != nil {
		return nil, err
	} else if o != int64(origOffset) {
		return nil, ErrTruncated
	}

	for i := 0; i < int(stringCount); i++ {
		if err := binary.Read(r, bo, &stringOffsets[i].origLen); err != nil {
			return nil, err
		}
		if err := binary.Read(r, bo, &stringOffsets[i].origStart); err != nil {
			return nil, err
		}
	}

	if o, err := r.Seek(int64(transOffset), 0); err != nil {
		return nil, err
	} else if o != int64(transOffset) {
		return nil, ErrTruncated
	}

	for i := 0; i < int(stringCount); i++ {
		if err := binary.Read(r, bo, &stringOffsets[i].transLen); err != nil {
			return nil, err
		}
		if err := binary.Read(r, bo, &stringOffsets[i].transStart); err != nil {
			return nil, err
		}
	}

	catalog := Catalog{
		Strings:       make(map[string]*Translation, stringCount),
		PluralFormula: GermanicPluralFormula,
	}

	for _, so := range stringOffsets {
		if o, err := r.Seek(int64(so.origStart), 0); err != nil {
			return nil, err
		} else if o != int64(so.origStart) {
			return nil, ErrTruncated
		}
		origBytes := make([]byte, so.origLen)
		if _, err := r.Read(origBytes); err != nil {
			return nil, err
		}

		if o, err := r.Seek(int64(so.transStart), 0); err != nil {
			return nil, err
		} else if o != int64(so.transStart) {
			return nil, ErrTruncated
		}
		transBytes := make([]byte, so.transLen)
		if _, err := r.Read(transBytes); err != nil {
			return nil, err
		}

		if len(origBytes) == 0 {
			// Translation meta header
			header, err := textproto.NewReader(
				bufio.NewReader(bytes.NewReader(transBytes))).ReadMIMEHeader()
			if err != nil {
				catalog.Header = header
			} else {
				catalog.Header = textproto.MIMEHeader(make(map[string][]string))
			}
		} else {
			origParts := strings.Split(string(origBytes), "\x00")
			transParts := strings.Split(string(transBytes), "\x00")
			if len(transParts) > 0 {
				t := &Translation{
					Translation: transParts,
				}
				catalog.Strings[origParts[0]] = t
				if len(origParts) > 1 {
					t.Plural = origParts[1]
				}
			}
		}
	}

	return &catalog, nil
}
