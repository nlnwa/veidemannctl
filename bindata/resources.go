// Code generated by go-bindata.
// sources:
// res/crawllog.template
// res/pagelog.template
// res/screenshot.template
// res/version
// DO NOT EDIT!

package bindata

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _crawllogTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x93\xcd\x6e\xdb\x30\x0c\x80\xef\x7d\x0a\x82\x80\x8f\x11\x9a\x9f\xa6\xe9\x31\xc8\x32\xc0\x40\xd3\x0d\x89\xb7\x5e\xe3\x59\x74\x22\xa0\x91\x3d\x4a\xee\x96\x19\x7e\xf7\x81\x52\xe7\x24\x4b\x0f\xbd\x10\x96\x04\x92\x1f\xf8\xd1\x6d\xab\xa9\x34\x96\x00\x1f\xd3\xa7\x25\xc2\xa0\xeb\x6e\xda\xb6\x66\x63\x7d\x09\xdb\x64\x3c\x55\xe3\xa9\x83\x64\x30\xbc\x93\xf8\x30\x52\x0f\x23\xf9\x98\xa9\x99\x83\xb7\xd7\x2d\xa8\xe7\x9c\x8b\x54\x83\xfa\x64\x5c\x51\xbd\x12\x1f\xbf\xe6\x7e\x0f\x6a\x4d\x3f\x1b\x72\x9e\xf4\x37\x36\x72\x2a\x2a\xd6\xd9\xb1\xa6\x98\xb0\xa6\x92\xd8\x65\x55\xe8\xf8\x83\xcd\x6e\xef\x8b\x63\x6e\xbb\xee\x03\xfd\x87\x77\x2a\xdc\x8c\x35\x24\xc3\xa1\x86\xe4\x5e\x43\x32\xd2\x5b\x50\xcb\xdf\x54\x34\xde\x54\x56\x78\xd2\x7a\xae\x35\x93\x73\xd2\xdd\xd5\x95\x75\x14\x50\x16\x95\xf5\x64\x7d\x64\xd9\xf8\xdc\x37\x6e\x51\x69\xf9\x36\x7f\x08\xd4\x67\xf2\xc5\x3e\x33\x07\x5a\x85\x44\xcf\x86\x9c\x60\x31\x39\xf2\x01\xf7\x0a\x74\x30\xbd\x55\xd3\xdb\x30\x9a\x89\x9a\x4d\x64\x28\x1b\x5f\x71\xbe\xa3\x35\x95\x52\xa4\x24\x66\xe2\xae\x83\xb6\xf5\xe6\x40\xa0\xa4\xfe\xc6\xe7\x87\x3a\xdc\xfd\x32\x32\xb0\x25\x73\xc5\x52\x37\x8e\x83\x49\xcb\x41\x49\x20\xab\x2f\x11\xc8\xea\x60\xeb\xe6\xa4\x70\x35\x4f\x9f\xfe\x29\x34\xf6\x95\xd8\xd1\x25\xe4\xd5\x38\x03\xb0\x84\xc9\xc8\x6d\x01\x9f\xe7\xeb\x05\xa4\x1a\x01\x45\x20\x02\x9e\x1b\x0c\xc7\x42\x66\x86\x80\xe7\x02\xf1\x12\xec\x52\xe5\xc7\x39\xa2\x50\x27\x42\x1d\x24\xf7\xe1\x6a\x28\x54\xbd\xd2\x88\xd6\x4b\x0d\x40\xbd\x55\x04\x3c\xd3\x1a\xdf\x10\x50\x8c\x22\x60\x50\xba\x8a\x29\x9e\x8f\xf8\xbe\xce\x77\x61\xcf\xcc\x4a\x1c\xc5\x28\x7b\x8f\x27\xc5\xa1\x70\x74\x8c\x80\xbd\x5b\x04\x0c\x4e\xff\xeb\xc7\xb9\xdd\x11\xa8\xef\xf9\x4b\x13\x9a\x79\x3a\xd4\x2f\xb9\xef\xff\xc2\x93\xf1\x9b\xec\x4b\x36\x7f\x84\x37\x1e\xd9\xdc\xc6\x5e\xf9\x3f\xa5\xc7\x0d\x50\x5d\xf7\x37\x00\x00\xff\xff\x98\x4f\x85\xcc\xd8\x03\x00\x00")

func crawllogTemplateBytes() ([]byte, error) {
	return bindataRead(
		_crawllogTemplate,
		"crawllog.template",
	)
}

func crawllogTemplate() (*asset, error) {
	bytes, err := crawllogTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "crawllog.template", size: 984, mode: os.FileMode(436), modTime: time.Unix(1518713311, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _pagelogTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x92\xcf\x6e\xe2\x30\x10\xc6\xef\x7d\x8a\x91\xa5\x1c\x63\x35\x8e\x28\xa2\x37\x94\x65\x57\xd1\xb6\xa5\x4a\x61\x7b\xc5\x24\x53\xb0\x36\x38\xc8\x76\xba\x8b\x22\xde\x7d\x35\x6e\xc9\x1f\x16\xd4\x43\x24\x7b\x26\x9e\xf9\xbe\xf9\x4d\xd3\x14\xf8\xa6\x34\x02\xcb\x66\x2f\xf3\x65\x96\xcc\x18\x84\xc7\xe3\x4d\xd3\xac\x8d\xda\x6c\xdd\x4e\x6e\x50\x3b\x79\x3c\x36\xcd\xde\x28\xed\xde\x60\x15\xc4\x77\x3c\xbe\xb3\x10\x84\x23\x07\x41\x5c\x40\x10\x8d\xe8\x36\x11\x7c\x22\xe8\x10\x09\x1e\x89\x53\x3e\x14\x82\x0b\x61\x57\xc0\x5f\xa5\xc9\xd3\x02\xf8\x77\x53\xed\x12\x99\x6f\x11\xf8\x8b\x93\xae\xb6\x49\x55\x20\xf0\x6f\xca\xe6\xd5\x3b\x9a\xc3\xb3\x74\x5b\xe0\x4b\xa3\x80\x67\x68\xab\xda\xe4\xb8\x38\xec\x91\x6e\xba\x40\x23\xd7\x25\x02\x7f\x54\x3b\x1f\x25\x61\x06\x2d\x3a\xaf\x19\x75\xe1\xd5\xdf\x74\xb6\xe6\xcb\xc5\x43\xfa\xf4\x73\xe8\x2a\x3f\x48\xdd\xb7\x04\x21\xc9\x1e\xdf\xf2\x68\x7c\x4b\x52\xbf\xac\xfa\x3c\xfd\xd1\x0e\x4a\xe9\x77\x34\x16\x07\x23\x0a\xdb\x19\xb5\x87\x48\xc4\x76\x05\xec\x75\x9a\x25\x90\x16\x0c\xd8\xec\x2f\xe6\xb5\x53\x95\xf6\xb7\xa5\x51\x6c\xd8\x76\x5d\x95\x05\x45\x0e\x58\x96\xd5\x9f\x8b\x04\xda\xe2\xfd\xf9\xf6\xea\xfa\x31\x0e\xab\xf6\xd4\x5e\x25\x7c\x1a\xbb\xbd\x87\x40\x8c\xb8\x20\xbc\xf4\xc5\xb6\x63\xfd\x61\xc9\x53\xe6\x3e\x24\x04\x1f\x0d\x0d\x7a\xc8\x0c\x58\x86\x96\x01\x23\xac\xe4\x33\x4b\x7d\x28\x27\x7a\xfe\xa4\xe9\x5f\xe2\x79\xe6\xdf\x48\xbd\xc1\x6e\x07\x28\xe9\x70\xb7\x2f\xa5\x1b\x6c\xab\x87\x75\x42\x74\xc1\xdf\x39\xeb\x20\x8c\x26\x63\x12\x3a\xaf\x5d\xa9\xf4\x6f\x7b\x7f\xb9\xef\x67\x7a\xd8\xb6\xdd\xa6\x53\xd7\x6b\x0b\xf2\x38\x4d\x9f\x58\xbf\xdc\x2f\x59\xd6\x67\x1e\x3e\x96\xa8\x57\x69\x31\x5f\x4c\x1f\xe0\x53\x2a\xf0\xa4\xaa\xf5\x7f\x1b\xd8\x3d\xf7\x2d\xe8\xf9\xbf\x00\x00\x00\xff\xff\x7c\xa2\xe6\x10\xc4\x03\x00\x00")

func pagelogTemplateBytes() ([]byte, error) {
	return bindataRead(
		_pagelogTemplate,
		"pagelog.template",
	)
}

func pagelogTemplate() (*asset, error) {
	bytes, err := pagelogTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pagelog.template", size: 964, mode: os.FileMode(436), modTime: time.Unix(1518713311, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _screenshotTemplate = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x54\xce\xc1\x8a\x83\x40\x0c\x06\xe0\xbb\x4f\x11\x02\x7b\x34\xa0\x03\xde\x97\xc5\x83\x20\xbb\x17\xdd\xb3\xd2\x89\x65\xc0\x46\x99\x19\x4b\x61\x98\x77\x2f\xa3\x45\xda\x4b\x48\x20\xff\x97\x84\xa0\x79\x32\xc2\x80\x6d\xf3\x5b\x23\xe4\x31\x86\xb0\x5a\x23\x7e\x82\xe1\x4b\x55\x0e\x8e\x92\x17\xa5\xa2\xa2\x54\x6e\x00\x6a\x34\x50\xfd\xe0\xcb\xe6\xcd\x22\x69\xe8\xad\x49\x31\x16\x9d\xf2\x59\x08\x46\xee\x6c\x1d\x7f\x58\xb9\xaa\xe8\xa0\xce\xe6\xf0\xb0\xd1\x08\xf8\x06\x22\x60\x6f\x0d\xa6\xb4\x65\xc7\x7e\x27\xed\x28\x57\x06\xfa\x1f\xe7\x6d\x77\x3d\xdf\xd6\x79\xf4\xe7\xe3\xb4\x6f\xb1\xe8\x18\xb3\xee\xaf\xfb\x6e\xe1\x75\x1b\xe8\x67\xd9\x24\x21\xcf\x00\x00\x00\xff\xff\x79\xbe\x4d\x7f\xec\x00\x00\x00")

func screenshotTemplateBytes() ([]byte, error) {
	return bindataRead(
		_screenshotTemplate,
		"screenshot.template",
	)
}

func screenshotTemplate() (*asset, error) {
	bytes, err := screenshotTemplateBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "screenshot.template", size: 236, mode: os.FileMode(436), modTime: time.Unix(1518713311, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _version = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x32\x33\x33\x4c\x36\x4d\xb5\x34\x35\xb0\x34\x34\x34\x4c\x34\x4b\xb5\x4c\xb3\x30\x4b\x49\xb2\x30\xb7\x34\x4f\x35\x30\x36\x32\x31\x4c\x36\x34\x48\x4b\x31\x4b\x36\xd6\x4d\xc9\x2c\x2a\xa9\xe4\x02\x04\x00\x00\xff\xff\x24\x81\x9b\xee\x2f\x00\x00\x00")

func versionBytes() ([]byte, error) {
	return bindataRead(
		_version,
		"version",
	)
}

func version() (*asset, error) {
	bytes, err := versionBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "version", size: 47, mode: os.FileMode(436), modTime: time.Unix(1519832224, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"crawllog.template":   crawllogTemplate,
	"pagelog.template":    pagelogTemplate,
	"screenshot.template": screenshotTemplate,
	"version":             version,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"crawllog.template":   &bintree{crawllogTemplate, map[string]*bintree{}},
	"pagelog.template":    &bintree{pagelogTemplate, map[string]*bintree{}},
	"screenshot.template": &bintree{screenshotTemplate, map[string]*bintree{}},
	"version":             &bintree{version, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
