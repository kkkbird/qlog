package qlog

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

const (
	longTimeStamp   = "2006/01/02 15:04:05.000000Z07:00"
	shortTimeStamp  = "06/01/02 15:04:05.000"
	keyPrettyCaller = "prettycaller" // omitfunc, truncated
)

var (
	gRegisteredFormatters = make(map[string]reflect.Type)

	gPrettyCallFuncMap = map[string]func(*runtime.Frame) (function string, file string){
		"omitfunc":  prettyCallerOmitFunc,
		"truncated": prettyCallerTruncated,
	}
)

func registeFormatter(name string, typ reflect.Type) {
	gRegisteredFormatters[name] = typ
}

// copy from filepath.Base
func truncatedPath(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && os.IsPathSeparator(path[len(path)-1]) {
		path = path[0 : len(path)-1]
	}
	// Throw away volume name
	path = path[len(filepath.VolumeName(path)):]
	// Find the last element
	i := len(path) - 1
	for i >= 0 && !os.IsPathSeparator(path[i]) {
		i--
	}
	if i >= 0 {
		// find first dir
		j := i - 1
		for j >= 0 && !os.IsPathSeparator(path[j]) {
			j--
		}
		if j >= 0 {
			path = path[j+1:]
		}
	}
	// If empty now, it had only slashes.
	if path == "" {
		return string(filepath.Separator)
	}
	return path
}

func prettyCallerTruncated(caller *runtime.Frame) (function string, file string) {
	//funcVal := ""
	fileVal := fmt.Sprintf("%s:%d", truncatedPath(caller.File), caller.Line)
	return "", fileVal
}

func prettyCallerOmitFunc(caller *runtime.Frame) (function string, file string) {
	fileVal := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	return "", fileVal
}

func newFormatter(name string, key string) (logrus.Formatter, error) {
	var err error
	var typ reflect.Type
	var ok bool

	if typ, ok = gRegisteredFormatters[name]; !ok {
		return nil, fmt.Errorf("[qlog] formatter name(%s) not registered", name)
	}

	f := reflect.New(typ)

	if err = v.UnmarshalKey(key, f.Interface()); err != nil {
		return nil, err
	}

	// check if we need truncate caller
	prettyCaller := v.GetString(key + "." + keyPrettyCaller)

	if len(prettyCaller) > 0 {
		if prettyFunc, ok := gPrettyCallFuncMap[prettyCaller]; ok {
			prettyFuncField := f.Elem().FieldByName("CallerPrettyfier")
			if prettyFuncField.IsValid() {
				prettyFuncField.Set(reflect.ValueOf(prettyFunc))
			} else {
				return nil, fmt.Errorf("[qlog] formatter name(%s) doesn't support truncate caller", name)
			}
		} else {
			return nil, fmt.Errorf("[qlog] formatter name(%s) init with unsupported pretty func:%s", name, prettyCaller)
		}
	}

	return f.Interface().(logrus.Formatter), nil
}
