package goutils

import "sync"

// FreeAble 释放实例
type FreeAble interface {
	Free()
}

type Reset interface {
	Reset()
}

// WaitFreeAble 等待合适时机自动释放
type WaitFreeAble interface {
	WaitFree()
}

type PoolElem interface {
	SetPool(p Pool)
}

type Pool interface {
	Get() interface{}
	Put(interface{})
}

type pool struct {
	p *sync.Pool
}

func NewPool(new func() interface{}) Pool {
	return &pool{
		p: &sync.Pool{
			New: new,
		},
	}
}

func (p *pool) Get() interface{} {
	i := p.p.Get()
	if x, ok := i.(Reset); ok {
		x.Reset()
	}
	if x, ok := i.(PoolElem); ok {
		x.SetPool(p)
	}
	return i
}

func (p *pool) Put(i interface{}) {
	p.p.Put(i)
}

func Free(x ...interface{}) {
	for _, o := range x {
		if o == nil {
			continue
		}
		if r, ok := o.(FreeAble); ok {
			r.Free()
		}
	}
}
