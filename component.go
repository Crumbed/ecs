package ecs

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

func partition[T cmp.Ordered](arr []T, lo, hi uint64) uint64 {
    pivot := arr[hi]
    i := lo - 1
    for j := lo; j <= hi - 1; j += 1 {
        if arr[j] >= pivot { continue }
        i += 1
        arr[i], arr[j] = arr[j], arr[i]
    }

    arr[i + 1], arr[hi] = arr[hi], arr[i + 1]
    return i + 1
}

func QuickSort[T cmp.Ordered](arr []T, lo, hi uint64) {
    if lo >= hi { return }
    pi := partition(arr, lo, hi)
    QuickSort(arr, lo, pi - 1)
    QuickSort(arr, pi + 1, hi)
}

// Hashed slice of `ComponentHandle`s because you cant use slices in maps
type ComponentHash [32]byte
func CreateComponentHash(handles ...ComponentHandle) ComponentHash {
    QuickSort(handles, 0, uint64(len(handles) - 1))
    return CreateComponentHashSorted(handles...)
}

// Creates a `ComponentQuery` without checking if `handles` is sorted.
// Only use this if you are 100% sure that the handles passed in are already sorted.
func CreateComponentHashSorted(handles ...ComponentHandle) ComponentHash {
    var buffer bytes.Buffer
    for _, h := range handles {
        buffer.WriteString(strconv.Itoa(int(h)))
        buffer.WriteByte(' ')
    }

    return ComponentHash(sha256.Sum256([]byte(buffer.String())))
}

type ComponentType struct {
    Handle  ComponentHandle
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

func CreateComponentType(c Component, handle ComponentHandle) ComponentType {
    raw := *(*rawInterface)(unsafe.Pointer(&c))
    return ComponentType {
        Handle: handle,
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

func (list *ComponentList) Get(i EntityHandle) Component {
    p := list.GetPtr(uintptr(i))
    return MakeComponent(list.ptr_t, p)
}

func (list *ComponentList) Set(i uint64, c unsafe.Pointer) {
    p := list.GetPtr(uintptr(i))
    memcpy(p, c, list.size_t)
}

func (list *ComponentList) Add() {
    if list.len == list.cap {
        //slice := unsafe.Slice((*uint8)(list.array), list.len * uint64(list.size_t))
        newCap := uint64(float64(list.cap) * 1.5) // * uint64(list.size_t)
        slice := make([]uint8, newCap * uint64(list.size_t))
        newPtr := unsafe.Pointer(&slice[0])
        memcpy(newPtr, list.array, uintptr(list.cap) * list.size_t)
        list.cap = newCap
        list.array = newPtr
    }

    list.len += 1
}

func (list *ComponentList) Remove(i uint64) {
    if i == list.len - 1 { list.Pop(); return }
    // how much data needs to be moved
    len := uintptr(list.len - 1 - i) * list.size_t
    dest_p := list.GetPtr(uintptr(i))
    dest := unsafe.Slice((*uint8)(dest_p), len)
    // data to be copied
    src_p := list.GetPtr(uintptr(i + 1))
    src := unsafe.Slice((*uint8)(src_p), len)
    copy(dest, src)
    list.len -= 1
}

func (list *ComponentList) Pop() Component {
    list.len -= 1
    return list.Get(EntityHandle(list.len))
}
