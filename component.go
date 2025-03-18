package ecs

import (
    "reflect"
    "unsafe"
    "math"
)

type ComponentType struct {
    typ     reflect.Type
    ptr_t   unsafe.Pointer
}
type ComponentHandle uint64;
const InvalidComponentHandle ComponentHandle = math.MaxUint64

type Component interface {
    // When defining a component you must also declare a global variable
    // of type `ComponentHandle`. This function must return the value of said global variable.
    // While adding this component to the ECS using `WithComponentType(*ComponentHandle, components)`,
    // the `*ComponentHandle` pointer will be to this global value.
	GetComponentHandle() ComponentHandle
}
type rawInterface struct {
    typ     unsafe.Pointer
    data    unsafe.Pointer
}

func CreateComponentType(c Component) ComponentType {
    raw := *(*rawInterface)(unsafe.Pointer(&c))
    return ComponentType {
        typ: reflect.TypeOf(c).Elem(),
        ptr_t: raw.typ,
    }
}

func MakeComponent(t unsafe.Pointer, data unsafe.Pointer) Component {
    empty := &rawInterface {
        typ: t,
        data: data,
    }
    return *(*Component)(unsafe.Pointer(empty))
}

func memcpy(dest unsafe.Pointer, src unsafe.Pointer, len uintptr) unsafe.Pointer {
    cnt := len >> 3
    var i uintptr
    for i = 0; i < cnt; i++ {
        pdest := (*uint64)(unsafe.Pointer(uintptr(dest) + i * 8))
        psrc := (*uint64)(unsafe.Pointer(uintptr(src) + i * 8))
        *pdest = *psrc
    }
    left := len & 7
    for i = 0; i < left; i++ {
        pdest := (*uint8)(unsafe.Pointer(uintptr(dest) + 8 * cnt + i))
        psrc := (*uint8)(unsafe.Pointer(uintptr(src) + 8 * cnt + i))
        *pdest = *psrc
    }

    return dest
}

type ComponentList struct {
    // void* of structs that implement Component
    array   unsafe.Pointer
    ptr_t   unsafe.Pointer
    size_t  uintptr
    len     uint64
    cap     uint64
}

func NewComponentList(t ComponentType) ComponentList {
    size_t := t.typ.Size()
    array := make([]uint8, size_t * 10);

    return ComponentList {
        array: unsafe.Pointer(&array[0]),
        ptr_t: t.ptr_t,
        size_t: size_t,
        len: 0,
        cap: 10,
    }
}

func (list *ComponentList) GetPtr(i uintptr) unsafe.Pointer {
    return unsafe.Pointer(uintptr(list.array) + i * list.size_t)
}

func (list *ComponentList) Get(i uint64) Component {
    p := list.GetPtr(uintptr(i))
    return MakeComponent(list.ptr_t, p)
}

func (list *ComponentList) Set(i uint64, c unsafe.Pointer) {
    p := list.GetPtr(uintptr(i))
    memcpy(p, c, list.size_t)
}

func (list *ComponentList) Add() Component {
    if list.len == list.cap {
        //slice := unsafe.Slice((*uint8)(list.array), list.len * uint64(list.size_t))
        newCap := uint64(float64(list.cap) * 1.5) * uint64(list.size_t)
        slice := make([]uint8, newCap)
        newPtr := unsafe.Pointer(&slice[0])
        memcpy(newPtr, list.array, uintptr(list.cap) * list.size_t)
        list.cap = newCap
        list.array = newPtr
    }

    i := list.len
    list.len += 1
    return list.Get(i)
}

func (list *ComponentList) Pop() Component {
    list.len -= 1
    return list.Get(list.len)
}
