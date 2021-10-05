package doublejump

import (
	"github.com/dgryski/go-jump"
)

type looseHolder struct {
	a []interface{}
	m map[interface{}]int
	f []int
}

func (this *looseHolder) add(obj interface{}) {
	if _, ok := this.m[obj]; ok {
		return
	}

	if nf := len(this.f); nf == 0 {
		this.a = append(this.a, obj)
		this.m[obj] = len(this.a) - 1
	} else {
		idx := this.f[nf-1]
		this.f = this.f[:nf-1]
		this.a[idx] = obj
		this.m[obj] = idx
	}
}

func (this *looseHolder) remove(obj interface{}) {
	if idx, ok := this.m[obj]; ok {
		this.f = append(this.f, idx)
		this.a[idx] = nil
		delete(this.m, obj)
	}
}

func (this *looseHolder) get(key uint64) interface{} {
	na := len(this.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(key, na)
	return this.a[h]
}

func (this *looseHolder) shrink() {
	if len(this.f) == 0 {
		return
	}

	var a []interface{}
	for _, obj := range this.a {
		if obj != nil {
			a = append(a, obj)
			this.m[obj] = len(a) - 1
		}
	}
	this.a = a
	this.f = nil
}

type compactHolder struct {
	a []interface{}
	m map[interface{}]int
}

func (this *compactHolder) add(obj interface{}) {
	if _, ok := this.m[obj]; ok {
		return
	}

	this.a = append(this.a, obj)
	this.m[obj] = len(this.a) - 1
}

func (this *compactHolder) shrink(a []interface{}) {
	for i, obj := range a {
		this.a[i] = obj
		this.m[obj] = i
	}
}

func (this *compactHolder) remove(obj interface{}) {
	if idx, ok := this.m[obj]; ok {
		n := len(this.a)
		this.a[idx] = this.a[n-1]
		this.m[this.a[idx]] = idx
		this.a[n-1] = nil
		this.a = this.a[:n-1]
		delete(this.m, obj)
	}
}

func (this *compactHolder) get(key uint64, nl int) interface{} {
	na := len(this.a)
	if na == 0 {
		return nil
	}

	h := jump.Hash(uint64(float64(key)/float64(nl)*float64(na)), na)
	return this.a[h]
}
