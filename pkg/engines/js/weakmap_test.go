package js

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"

	gm "github.com/onsi/gomega"
)

func TestWeakMapConstructor(t *testing.T) {
	gm.RegisterTestingT(t)

	t.Run("Creates a WeakMap", func(t *testing.T) {
		cases := map[string]string{
			"without constructor": `exports.weakmap = new WeakMap();`,
			"with constructor": `
				it = makeIterator(
					[
						[{"one": "two"}, 1]
					]
				);

				exports.weakmap = new WeakMap(it);
			`,
		}
		for name, code := range cases {
			t.Run(name, func(t *testing.T) {
				exports, jsvm := exportWeakMap(code)

				jsWeakMap := exports.Get("weakmap")
				gm.Expect(jsWeakMap).NotTo(gm.BeNil())
				gm.Expect(goja.IsUndefined(jsWeakMap) && goja.IsNull(jsWeakMap), true)

				nativeWeakMap := jsWeakMap.ToObject(jsvm).Export()
				gm.Expect(nativeWeakMap).To(gm.HaveKey("get"))
				gm.Expect(nativeWeakMap).To(gm.HaveKey("set"))
				gm.Expect(nativeWeakMap).To(gm.HaveKey("clear"))
				gm.Expect(nativeWeakMap).To(gm.HaveKey("delete"))
				gm.Expect(nativeWeakMap).To(gm.HaveKey("has"))
				gm.Expect(nativeWeakMap).ToNot(gm.HaveKey("doesnotexist"))
			})
		}
	})

	t.Run("Get and from a WeakMap", func(t *testing.T) {
		exports, jsvm := exportWeakMap(`
			var obj = {"one": "two"};
			it = makeIterator(
				[
					[obj, 1]
				]
			);

			var wm = new WeakMap(it);
			exports.value = wm.get(obj);
		`)

		e := exports.ToObject(jsvm).Export().(map[string]interface{})
		gm.Expect(e).To(gm.HaveKey("value"))
		gm.Expect(e["value"]).To(gm.BeEquivalentTo(1))
	})

	t.Run("Clears a WeakMap", func(t *testing.T) {
		exports, jsvm := exportWeakMap(`
		var obj = {"one": "two"}
		it = makeIterator(
			[
				[obj, 1]
			]
		);

		var wm = new WeakMap(it)
		exports.oldValue = wm.get(obj)
		exports.cleared = wm.clear()
		exports.value = wm.get(obj)
		`)

		e := exports.ToObject(jsvm).Export().(map[string]interface{})
		gm.Expect(e).To(gm.HaveKey("cleared"))
		gm.Expect(e).To(gm.HaveKey("value"))
		gm.Expect(e["cleared"]).To(gm.BeNil())
		gm.Expect(e["oldValue"]).To(gm.BeEquivalentTo(1))
		gm.Expect(e["value"]).To(gm.BeNil())
	})

	t.Run("Deletes an element from the WeakMap", func(t *testing.T) {
		exports, jsvm := exportWeakMap(`
		var obj = {"one": "two"}
		it = makeIterator(
			[
				[obj, 1]
			]
		);
	
		var wm = new WeakMap(it)
		exports.oldValue = wm.get(obj)
		exports.deleted = wm.delete(obj)
		exports.value = wm.get(obj)	
		`)

		e := exports.ToObject(jsvm).Export().(map[string]interface{})
		gm.Expect(e).To(gm.HaveKey("deleted"))
		gm.Expect(e).To(gm.HaveKey("value"))
		gm.Expect(e["deleted"]).To(gm.BeTrue())
		gm.Expect(e["oldValue"]).To(gm.BeEquivalentTo(1))
		gm.Expect(e["value"]).To(gm.BeNil())
	})

	t.Run("Checks an element in the WeakMap", func(t *testing.T) {
		exports, jsvm := exportWeakMap(`
			var obj = {"one": "two"}
			it = makeIterator(
				[
					[obj, 1]
				]
			);
		
			var wm = new WeakMap(it)
			exports.has = wm.has(obj)
			exports.notHas = wm.has({"no": "object"})
		`)
		e := exports.ToObject(jsvm).Export().(map[string]interface{})
		gm.Expect(e).To(gm.HaveKey("has"))
		gm.Expect(e).To(gm.HaveKey("notHas"))
		gm.Expect(e["has"]).To(gm.BeTrue())
		gm.Expect(e["notHas"]).To(gm.BeFalse())
	})

	t.Run("Sets an element in the WeakMap", func(t *testing.T) {
		exports, jsvm := exportWeakMap(`
			var obj = {"one": "two"}
		
			var wm = new WeakMap()
			wm.set(obj, 1);
			exports.set = wm.get(obj)
		`)
		e := exports.ToObject(jsvm).Export().(map[string]interface{})
		gm.Expect(e).To(gm.HaveKey("set"))
		gm.Expect(e["set"]).To(gm.BeEquivalentTo(1))
	})
}

func exportWeakMap(code string) (*goja.Object, *goja.Runtime) {
	const jsIter = `
	function makeIterator(array) {
		var nextIndex = 0;
	
		return {
		   next: function() {
			   return nextIndex < array.length ?
					{value: array[nextIndex++], done: false} :
					{done: true};
			}
		}
	}`

	jsvm := goja.New()
	wm := NewWeakMap(jsvm)
	wm.Enable()

	src := fmt.Sprintf(`(function(module, exports) {
		%s
		%s
	})`, jsIter, code)

	p, err := jsvm.RunString(src)
	if err != nil {
		panic(fmt.Sprintf("source: %s\n\nerror: %s", src, err.Error()))
	}
	gm.Expect(err).NotTo(gm.HaveOccurred())
	gm.Expect(p).NotTo(gm.BeNil())

	call, ok := goja.AssertFunction(p)
	gm.Expect(ok).To(gm.BeTrue())

	exports := jsvm.NewObject()
	_, err = call(p, nil, exports)
	if err != nil {
		panic(fmt.Sprintf("source: %s\n\nerror: %s", src, err.Error()))
	}
	gm.Expect(err).NotTo(gm.HaveOccurred())

	return exports, jsvm
}
