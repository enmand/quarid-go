package python

// #cgo pkg-config: python-2.7
// #cgo CFLAGS: -Wno-unused-variable -Wno-infinite-recursion
// #include "Python.h"
import "C"

import (
	"context"
	"runtime"
	"sync"

	"errors"
	"fmt"
	"github.com/enmand/quarid-go/pkg/vm"
	"github.com/liamzdenek/go-pthreads"
	py "github.com/sbinet/go-python"
)

func init() {
	runtime.LockOSThread()

	if C.Py_IsInitialized() == 0 {
		C.Py_SetProgramName(C.CString("quarid"))

		err := py.Initialize()
		if err != nil {
			panic("unable to initialize Python")
		}
	}

	if C.PyEval_ThreadsInitialized() == 0 {
		C.PyEval_InitThreads()
	}

	state := C.PyGILState_GetThisThreadState()
	C.PyEval_ReleaseThread(state)
}

type python struct {
	lock sync.Mutex
}

func New() vm.VM {
	return &python{
		lock: sync.Mutex{},
	}
}

func (vm *python) Run(code string, args ...interface{}) error {
	ctx := context.Background()

	return vm.run(ctx, code, "")
}

func (vm *python) RunFile(path string, args ...interface{}) error {
	return nil
}

func (vm *python) RunIn(code, entry string, args ...interface{}) error {
	return nil
}

func (vm *python) RunInFile(path, entry string, args ...interface{}) error {
	return nil
}

func (vm *python) run(ctx context.Context, code, entry string, args ...interface{}) error {
	vm.lock.Lock()
	defer vm.lock.Unlock()

	_, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
	}()

	threadCh := make(chan error)
	thread := pthread.Create(func() {
		state := py.PyGILState_Ensure()
		defer func() {
			py.PyGILState_Release(state)
		}()

		obj := py.PyMarshal_ReadObjectFromString(code)

		err := obj.Call(nil, nil)
		if py.PyErr_ExceptionMatches(err) {
			threadCh <- errors.New(py.PyString_AsString(err))
		}

		threadCh <- nil
	})

	defer func() {
		thread.Kill()
	}()

	return <-threadCh
}
