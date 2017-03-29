package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/rai_config.yml": {
		local:   "../rai_config.yml",
		size:    726,
		modtime: 1490817928,
		compressed: `
H4sIAAAJbogA/3yQzc6iSBiF914FezMfPwqiCckAAoLiL4KyMSWUUCAUFAWoVz/5nO60vejenvfkPU8e
UFWzAcOUoIAzhgA0AC1Nv5OK4A7FkMyY74QbMEyMC4DKd6siOIMR/XqfviJcDBgmuiNY0guKZ4yiqMZh
JCiKODUSaZia81uOYajnu+5wxFyT6M0ianlt6mIrUdd5/3w5dUL9245MbEtC5yHoLRlcd9tjUQKg1g45
3670sfk108CIQPoxJVSblKQjvFebIZZSPpzkkTvupVVlZitpMsyETvNcqzpM+Yxqm7HQL60st4+eTutj
m7tPYcmKfvFa8YJ4jtpt2SP/GqL+5JGCyBs2CeQFO99YRqBnthjXEVikduOavhF4orW+9jxb4Md13Jhj
ZQALgO6/W/xOkrb88PgufRGA/vmp83+RDW5JBGdMhRtagIZC8u+fqj+eXkBEUQcvoEKXHD4/rGw9Xtv7
7svNuOm0grfs7OFXlQg05p0xbcyS5eOFU6ohYXXHvgdwJfWBGLPckhjnJEC1q9LXTa0FPDzqZi3Kd62D
TuWFo224rlasv1OUD44356UDdxQDinD5F5qDPHf3zmmBs0QYJQ9vLqzTsNU5LW6lkSE/BGtyP2Wo2D5Z
n6MEhNmSN9hzl57YwN4DuZavXIIxby7afcoSaj81y93yMLoZyuC/AAAA//9tonaP1gIAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},
}
