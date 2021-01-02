package hashmember

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/hashicorp/go-multierror"
	cuckoo "github.com/seiflotfy/cuckoofilter"
)

type hashmember struct {
	*cuckoo.Filter
}

type Hashmember interface {
	Lookup(data []byte) bool
	Insert(data []byte) bool

	Encode() []byte
}

// ErrDecode wrapps errors returned by the encoder
var ErrDecode = errors.New("hashmember decode error")

// ErrUnhandledVersion returned when the version of the hashmember input is not
// parsable by the decoder likely because it was created by a newer version of hashmember
var ErrUnhandledVersion = errors.New("unhandled hashmember version")

// Decode reads an encoded Hashmember from a reader and returns a Hashmember or
// an ErrDecode wrapping the more specific error
//
// ErrDecode will wrap an error wrapping ErrUnhandledVersion on attempting to parse
// a newer version of hashmember than the current version is capable.
func Decode(r io.Reader) (Hashmember, error) {
	br := bufio.NewReader(r)

	line, _, err := br.ReadLine()
	if err != nil {
		return nil, multierror.Append(ErrDecode, err)
	}

	ver, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, multierror.Append(ErrDecode, err)
	}

	if ver != 1 {
		return nil, multierror.Append(ErrDecode, fmt.Errorf("%w: lte v1 expected", ErrUnhandledVersion))
	}

	data, err := ioutil.ReadAll(br)
	if err != nil {
		return nil, multierror.Append(ErrDecode, fmt.Errorf("read error: %w", err))
	}

	decData, err := filterDecode(data)
	if err != nil {
		return nil, multierror.Append(ErrDecode, fmt.Errorf("filter decode error: %w", err))
	}

	cf, err := cuckoo.Decode(decData)
	if err != nil {
		return nil, multierror.Append(ErrDecode, err)
	}

	return &hashmember{cf}, nil
}

// ErrEncode wrapps errors returned by the encoder
var ErrEncode = errors.New("hashmember encode error")

// Encode encodes and writes a given Hashmember to a given writer
//
// Returns a nil on success or an ErrEncode wrapping the specific IO error
// on failure
func Encode(w io.Writer, h Hashmember) error {
	_, err := w.Write([]byte("1\n"))
	if err != nil {
		return multierror.Append(ErrEncode, err)
	}

	encBytes, err := filterEncode(h.Encode())
	if err != nil {
		return multierror.Append(ErrEncode, err)
	}

	_, err = w.Write(encBytes)
	if err != nil {
		return multierror.Append(ErrEncode, err)
	}

	return nil
}

// New returns a new initilized Hashmember
func New() Hashmember {
	return hashmember{
		cuckoo.NewFilter(1000000),
	}
}

func filterEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(input)
	if err != nil {
		return nil, err
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return ioutil.ReadAll(&buf)
}

func filterDecode(data []byte) ([]byte, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
