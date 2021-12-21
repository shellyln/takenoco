package parser

import (
	"reflect"
	"unsafe"
)

// An interface that abstracts memory address.
type Box interface {
	// Get value as any
	GetAny() interface{}
	// Get value as int64
	GetInt() int64
	// Get value as uint64
	GetUint() uint64
	// Get value as float64
	GetFloat() float64
	// Get value as bool
	GetBool() bool
	// Get value as string
	GetString() string
	// Get value as []byte
	GetBytes() []byte
	// Get value as unsafe.Pointer
	GetPointer() unsafe.Pointer

	// Set any
	SetAny(v interface{})
	// Set int64 value
	SetInt(v int64)
	// Set uint64 value
	SetUint(v uint64)
	// Set float64 value
	SetFloat(v float64)
	// Set bool value
	SetBool(v bool)
	// Set string value
	SetString(v string)
	// Set []byte value
	SetBytes(v []byte)
	// Set unsafe.Pointer
	SetPointer(v unsafe.Pointer)

	// Get a specified element that is associated to the index
	Index(i int) Box
	// Get a specified element that is associated to the key
	MapIndex(k string) Box
	// Get a specified element that is associated to the key
	ComplexMapIndex(k interface{}) Box
}

// Implements the interface Box. Abstraction of reflect.Value.
type ReflectionBox struct {
	Val reflect.Value
}

//
func (s ReflectionBox) GetAny() interface{} {
	return s.Val.Interface()
}

//
func (s ReflectionBox) GetInt() int64 {
	return s.Val.Int()
}

//
func (s ReflectionBox) GetUint() uint64 {
	return s.Val.Uint()
}

//
func (s ReflectionBox) GetFloat() float64 {
	return s.Val.Float()
}

//
func (s ReflectionBox) GetBool() bool {
	return s.Val.Bool()
}

//
func (s ReflectionBox) GetString() string {
	return s.Val.String()
}

//
func (s ReflectionBox) GetBytes() []byte {
	return s.Val.Bytes()
}

//
func (s ReflectionBox) GetPointer() unsafe.Pointer {
	return unsafe.Pointer(s.Val.Pointer())
}

//
func (s ReflectionBox) SetAny(v interface{}) {
	s.Val.Set(reflect.ValueOf(v))
}

//
func (s ReflectionBox) SetInt(v int64) {
	s.Val.SetInt(v)
}

//
func (s ReflectionBox) SetUint(v uint64) {
	s.Val.SetUint(v)
}

//
func (s ReflectionBox) SetFloat(v float64) {
	s.Val.SetFloat(v)
}

//
func (s ReflectionBox) SetBool(v bool) {
	s.Val.SetBool(v)
}

//
func (s ReflectionBox) SetString(v string) {
	s.Val.SetString(v)
}

//
func (s ReflectionBox) SetBytes(v []byte) {
	s.Val.SetBytes(v)
}

//
func (s ReflectionBox) SetPointer(v unsafe.Pointer) {
	s.Val.SetPointer(v)
}

//
func (s ReflectionBox) Index(i int) Box {
	return ReflectionBox{Val: s.Val.Index(i)}
}

//
func (s ReflectionBox) MapIndex(k string) Box {
	return ReflectionBox{Val: s.Val.MapIndex(reflect.ValueOf(k))}
}

//
func (s ReflectionBox) ComplexMapIndex(k interface{}) Box {
	return ReflectionBox{Val: s.Val.MapIndex(reflect.ValueOf(k))}
}

// Implements the interface Box. Abstraction of Map.
type MapContainerReflectionBox struct {
	Container reflect.Value
	Key       reflect.Value
}

//
func (s MapContainerReflectionBox) GetAny() interface{} {
	return s.Container.MapIndex(s.Key).Interface()
}

//
func (s MapContainerReflectionBox) GetInt() int64 {
	return s.Container.MapIndex(s.Key).Int()
}

//
func (s MapContainerReflectionBox) GetUint() uint64 {
	return s.Container.MapIndex(s.Key).Uint()
}

//
func (s MapContainerReflectionBox) GetFloat() float64 {
	return s.Container.MapIndex(s.Key).Float()
}

//
func (s MapContainerReflectionBox) GetBool() bool {
	return s.Container.MapIndex(s.Key).Bool()
}

//
func (s MapContainerReflectionBox) GetString() string {
	return s.Container.MapIndex(s.Key).String()
}

//
func (s MapContainerReflectionBox) GetBytes() []byte {
	return s.Container.MapIndex(s.Key).Bytes()
}

//
func (s MapContainerReflectionBox) GetPointer() unsafe.Pointer {
	return unsafe.Pointer(s.Container.MapIndex(s.Key).Pointer())
}

//
func (s MapContainerReflectionBox) SetAny(v interface{}) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetInt(v int64) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetUint(v uint64) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetFloat(v float64) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetBool(v bool) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetString(v string) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetBytes(v []byte) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) SetPointer(v unsafe.Pointer) {
	s.Container.SetMapIndex(s.Key, reflect.ValueOf(v))
}

//
func (s MapContainerReflectionBox) Index(i int) Box {
	return ReflectionBox{Val: reflect.ValueOf(nil)}
}

//
func (s MapContainerReflectionBox) MapIndex(k string) Box {
	return MapContainerReflectionBox{Container: s.Container, Key: reflect.ValueOf(k)}
}

//
func (s MapContainerReflectionBox) ComplexMapIndex(k interface{}) Box {
	return MapContainerReflectionBox{Container: s.Container, Key: reflect.ValueOf(k)}
}

const msgReferencedNonInitializedVariable = "Error: Non initialized variable is referenced."

// Implements the interface Box. Abstraction of Map.
// If you try to get the value before it is set, you will get an error.
type NotInitializedMapContainerReflectionBox struct {
	Container   reflect.Value
	Key         reflect.Value
	Initialized bool
}

//
func (p *NotInitializedMapContainerReflectionBox) GetAny() interface{} {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Interface()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetInt() int64 {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Int()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetUint() uint64 {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Uint()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetFloat() float64 {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Float()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetBool() bool {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Bool()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetString() string {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).String()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetBytes() []byte {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return p.Container.MapIndex(p.Key).Bytes()
}

//
func (p *NotInitializedMapContainerReflectionBox) GetPointer() unsafe.Pointer {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return unsafe.Pointer(p.Container.MapIndex(p.Key).Pointer())
}

//
func (p *NotInitializedMapContainerReflectionBox) SetAny(v interface{}) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetInt(v int64) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetUint(v uint64) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetFloat(v float64) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetBool(v bool) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetString(v string) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetBytes(v []byte) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) SetPointer(v unsafe.Pointer) {
	p.Container.SetMapIndex(p.Key, reflect.ValueOf(v))
	p.Initialized = true
}

//
func (p *NotInitializedMapContainerReflectionBox) Index(i int) Box {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return ReflectionBox{Val: reflect.ValueOf(nil)}
}

//
func (p *NotInitializedMapContainerReflectionBox) MapIndex(k string) Box {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return MapContainerReflectionBox{Container: p.Container, Key: reflect.ValueOf(k)}
}

//
func (p *NotInitializedMapContainerReflectionBox) ComplexMapIndex(k interface{}) Box {
	if !p.Initialized {
		panic(msgReferencedNonInitializedVariable)
	}
	return MapContainerReflectionBox{Container: p.Container, Key: reflect.ValueOf(k)}
}
