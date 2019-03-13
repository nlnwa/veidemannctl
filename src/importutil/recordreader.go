package importutil

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type State struct {
	fileName string
	recNum   int
	err      error
}

func (r *State) GetFilename() string { return r.fileName }
func (r *State) GetRecordNum() int   { return r.recNum }
func (r *State) GetError() error     { return r.err }

type recordReader struct {
	recordDecoder RecordDecoder
	curFile       *os.File
	dir           *os.File
	curFileName   string
	curRecNum     int
	filePattern   string
}

type RecordDecoder interface {
	Init(r io.Reader, suffix string)
	Read(v interface{}) (err error)
}

func NewRecordReader(fileOrDir string, decoder RecordDecoder, filePattern string) (l *recordReader, err error) {
	l = &recordReader{
		recordDecoder: decoder,
		filePattern:   filePattern,
	}

	if fileOrDir == "-" {
		l.recordDecoder.Init(os.Stdin, "")
	} else {
		if hasMeta(fileOrDir) {
			l.filePattern = filepath.Base(fileOrDir)
			fileOrDir = filepath.Dir(fileOrDir)
		}
		var f *os.File
		f, err = os.Open(fileOrDir)
		if err != nil {
			log.Fatalf("Could not open file '%s': %v", fileOrDir, err)
			return
		}
		if fi, _ := f.Stat(); fi.IsDir() {
			l.dir = f
			l.initRecordReader()
		} else {
			l.curFileName = f.Name()
			log.Infof("Reading file %s", l.curFileName)
			l.recordDecoder.Init(f, filepath.Ext(f.Name()))
		}
	}
	return
}

// hasMeta reports whether path contains any of the magic characters
// recognized by filepath.Match.
func hasMeta(path string) bool {
	magicChars := `*?[`
	if runtime.GOOS != "windows" {
		magicChars = `*?[\`
	}
	return strings.ContainsAny(path, magicChars)
}

func (l *recordReader) initRecordReader() (err error) {
	if l.curFile != nil {
		l.curFile.Close()
	}

	if l.dir == nil {
		return io.EOF
	}

	var f []os.FileInfo
	for {
		f, err = l.dir.Readdir(1)
		if err != nil {
			return
		}
		fi := f[0]

		if !fi.IsDir() {
			var match bool
			if match, err = filepath.Match(l.filePattern, fi.Name()); match && err == nil {
				l.curFile, err = os.Open(filepath.Join(l.dir.Name(), fi.Name()))
				if err != nil {
					log.Fatalf("Could not open file '%s': %v", fi.Name(), err)
					return
				}
				l.curRecNum = 0
				l.curFileName = l.curFile.Name()

				log.Infof("Reading file %s", l.curFileName)
				l.recordDecoder.Init(l.curFile, filepath.Ext(l.curFile.Name()))
				return
			}
		}
	}
	return
}

func (l *recordReader) Next(target interface{}) (s *State, err error) {
	err = l.recordDecoder.Read(target)
	if err == io.EOF {
		if err = l.initRecordReader(); err != nil {
			return
		}
		return l.Next(target)
	}
	l.curRecNum++

	s = &State{}
	s.fileName = l.curFileName
	s.recNum = l.curRecNum
	s.err = err
	return
}

type LineAsStringDecoder struct {
	r *bufio.Reader
}

func (l *LineAsStringDecoder) Init(r io.Reader, suffix string) {
	l.r = bufio.NewReader(r)
}

func (l *LineAsStringDecoder) Read(v interface{}) (err error) {
	s, err := l.r.ReadString('\n')
	if err != nil {
		return
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("invalid target: %v", reflect.TypeOf(v))
	}

	s = strings.TrimSpace(s)
	rv.Elem().Set(reflect.ValueOf(&s).Elem())

	return
}

type JsonYamlDecoder struct {
	dec jyDecoder
}

func (j *JsonYamlDecoder) Init(r io.Reader, suffix string) {
	if suffix == ".yaml" || suffix == ".yml" {
		dec := yaml.NewDecoder(r)
		dec.SetStrict(true)
		j.dec = dec
		return
	} else {
		dec := json.NewDecoder(r)
		dec.DisallowUnknownFields()
		j.dec = dec
		return
	}
}

func (j *JsonYamlDecoder) Read(v interface{}) (err error) {
	return j.dec.Decode(v)
}

type jyDecoder interface {
	Decode(v interface{}) (err error)
}
