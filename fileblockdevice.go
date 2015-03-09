package everdb

import "errors"
import "os"

const (
	BlockSize = 4096
)

type FileBlockDevice struct {
	file *os.File
	len  uint32
}

func NewFileBlockDevice(fname string, readonly bool, overwrite bool) (*FileBlockDevice, error) {
	if readonly && overwrite {
		return nil, errors.New("Both overwrite and readonly specified")
	}

	f := new(FileBlockDevice)

	flags := 0
	switch {
	case readonly:
		flags = os.O_RDONLY
	case overwrite:
		flags = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	default:
		flags = os.O_CREATE | os.O_RDWR
	}

	var err error
	f.file, err = os.OpenFile(fname, flags, 0666)
	if nil != err {
		return nil, err
	}

	fi, err := f.file.Stat()
	if nil != err {
		return nil, errors.New("Failed to stat file")
	} else if fi.Size()%BlockSize != 0 {
		return nil, errors.New("File is corrupt (not a multiple of BlockSize")
	}

	f.len = (uint32)(fi.Size() / BlockSize)

	return f, nil
}

func (f *FileBlockDevice) Len() uint32 {
	return f.len
}

func (f *FileBlockDevice) Resize(len uint32) error {
	sz := (int64)(len) * BlockSize

	err := f.file.Truncate(sz)
	if nil != err {
		return err
	}

	fi, err := f.file.Stat()
	if nil != err {
		return errors.New("Failed to stat file")
	} else if fi.Size() != sz {
		return errors.New("Failed to resize file")
	}
	f.len = len

	return nil
}

func (f *FileBlockDevice) Get(block int) ([]byte, error) {
	b := make([]byte, BlockSize)

	n, err := f.file.ReadAt(b, (int64)(block)*BlockSize)
	if nil != err {
		return nil, err
	} else if n != BlockSize {
		return nil, errors.New("Failed to read BlockSize bytes")
	}

	return b, nil
}
