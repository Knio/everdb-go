package everdb

import "errors"
import "os"

const BLOCK_SIZE = 4096

type FileBlockDevice struct {
	file *os.File
	len  uint32
}

func NewFileBlockDevice(fname string, readonly bool, overwrite bool) (*FileBlockDevice, error) {
	if readonly && overwrite {
		return nil, errors.New("Both overwrite and readonly specified")
	}

	f := new(FileBlockDevice)
	var err error

	if readonly {
		f.file, err = os.OpenFile(fname, os.O_RDONLY, 0)
	} else if overwrite {
		f.file, err = os.OpenFile(fname, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	} else {
		f.file, err = os.OpenFile(fname, os.O_CREATE|os.O_RDWR, 0666)
	}

	if nil != err {
		return nil, err
	}

	fi, err := f.file.Stat()
	if nil != err {
		return nil, errors.New("Failed to stat file")
	}

	if fi.Size()%BLOCK_SIZE != 0 {
		return nil, errors.New("File is corrupt (not a multiple of BLOCK_SIZE")
	}

	f.len = (uint32)(fi.Size() / BLOCK_SIZE)

	return f, nil
}

func (f *FileBlockDevice) Len() uint32 {
	return f.len
}

func (f *FileBlockDevice) Resize(len uint32) error {
	sz := (int64)(len) * BLOCK_SIZE

	err := f.file.Truncate(sz)
	if nil != err {
		return err
	}

	fi, err := f.file.Stat()
	if nil != err {
		return errors.New("Failed to stat file")
	}
	if fi.Size() != sz {
		return errors.New("Failed to reize file")
	}
	f.len = len

	return nil
}

func (f *FileBlockDevice) Get(block int) ([]byte, error) {
	b := make([]byte, BLOCK_SIZE)
	n, err := f.file.ReadAt(b, (int64)(block)*BLOCK_SIZE)
	if nil != err {
		return nil, err
	}
	if n != BLOCK_SIZE {
		return nil, errors.New("Failed to read BLOCK_SIZE bytes")
	}
	return b, nil
}
