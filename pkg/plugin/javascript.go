package plugin

import (
	"github.com/enmand/quarid-go/pkg/adapter"
	"github.com/enmand/quarid-go/pkg/engines"
	"github.com/enmand/quarid-go/pkg/engines/js"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

type javascript struct {
	path   string
	source string
	jsvm   engines.Engine
}

func newJavaScript(path string, source string) (plugin, error) {
	vm := js.New()
	if err := vm.Load(); err != nil {
		return nil, errors.Wrap(err, "unable to load javascript engine")
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
	os := s.(otto.Value)
	obj := os.Object()

	fs := []*otto.Object{}

	for _, k := range obj.Keys() {
		v, err := obj.Get(k)
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse filter keys")
		}
		f, err := v.Object().Get("filters")
		if err != nil {
			return nil, errors.Wrap(err, "unable to find filters")
		}

		fs = append(fs, f.Object())
	}

	return nil, nil
}
