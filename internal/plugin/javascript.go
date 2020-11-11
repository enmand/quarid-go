package plugin

import (
	"fmt"

	"github.com/dop251/goja"
	"github.com/enmand/quarid-go/internal/adapter"
	"github.com/enmand/quarid-go/internal/engines"
	"github.com/enmand/quarid-go/internal/engines/js"
	"github.com/kr/pretty"
	"github.com/internal/errors"
)

type javascript struct {
	path   string
	source string
	jsvm   engines.Engine
}

func newJavaScript(path string, source string) (plugin, error) {
	vm, err := js.New()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create JavaScript engine")
	}
	if err := vm.Load(); err != nil {
		return nil, errors.Wrap(err, "unable to load JavaScript engine")
	}

	return &javascript{
		path:   path,
		source: source,
		jsvm:   vm,
	}, nil
}

func (js *javascript) Filters() ([]adapter.Filter, error) {
	s, err := js.jsvm.Compile(js.path, js.source)
	if err != nil {
		return nil, errors.Wrap(err, "unable to compile script")
	}
	fmt.Printf("%# v", pretty.Formatter(s))
	os := s.(goja.Value)
	fmt.Printf("%# v", pretty.Formatter(os))
	obj := os.ToObject(js.jsvm.Runtime().(*goja.Runtime))
	fmt.Printf("%# v", pretty.Formatter(obj))

	fs := []*goja.Object{}

	for _, k := range obj.Keys() {
		v := obj.Get(k)
		if goja.IsUndefined(v) || goja.IsNull(v) {
			return nil, errors.Wrap(err, "unable to parse filter keys")
		}
		f := v.ToObject(js.jsvm.Runtime().(*goja.Runtime)).Get("filters")
		if goja.IsUndefined(f) || goja.IsNull(f) {
			return nil, errors.Wrap(err, "unable to find filters")
		}

		fs = append(fs, f.ToObject(js.jsvm.Runtime().(*goja.Runtime)))
	}

	return nil, nil
}
