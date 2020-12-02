package redis

import (
	"encoding/gob"
	"reflect"
	"strconv"
)

// Value gob 转一些基础类型
type Value struct {
	decoder *gob.Decoder
}

// Decode 解析 redis 原始值, 调用此方法需要明确的知道在该值存入 redis 是以什么类型序列化存入的, 否则会解析失败
func (v Value) Decode(val interface{}) error {
	if v.decoder == nil {
		return nil
	}
	return v.decoder.Decode(val)
}

// DecodeValue 同 Decode, 不过这里的 val 是传入的类型的 reflect.Value
func (v Value) DecodeValue(val reflect.Value) error {
	if v.decoder == nil {
		return nil
	}
	return v.decoder.DecodeValue(val)
}

// String 将 redis 的原始值直接以 string 返回
func (v Value) String() (val string, err error) {
	err = v.Decode(&val)
	return
}

// Int 将 redis 的原始值转为 int
func (v Value) Int() (val int, err error) {
	err = v.Decode(&val)
	return
}

// Int8 将 redis 的原始值转为 int8
func (v Value) Int8() (val int8, err error) {
	err = v.Decode(&val)
	return
}

// Int16 将 redis 的原始值转为 int16
func (v Value) Int16() (val int16, err error) {
	err = v.Decode(&val)
	return
}

// Int32 将 redis 的原始值转为 int32
func (v Value) Int32() (val int32, err error) {
	err = v.Decode(&val)
	return
}

// Int64 将 redis 的原始值转为 int64
func (v Value) Int64() (val int64, err error) {
	err = v.Decode(&val)
	return
}

// Uint 将 redis 的原始值转为 uint
func (v Value) Uint() (val uint, err error) {
	err = v.Decode(&val)
	return
}

// Uint8 将 redis 的原始值转为 uint8
func (v Value) Uint8() (val uint8, err error) {
	err = v.Decode(&val)
	return
}

// Uint16 将 redis 的原始值转为 uint16
func (v Value) Uint16() (val uint16, err error) {
	err = v.Decode(&val)
	return
}

// Uint32 将 redis 的原始值转为 uint32
func (v Value) Uint32() (val uint32, err error) {
	err = v.Decode(&val)
	return
}

// Uint64 将 redis 的原始值转为 uint64
func (v Value) Uint64() (val uint64, err error) {
	err = v.Decode(&val)
	return
}

// Float32 将 redis 的原始值转为 float32
func (v Value) Float32() (val float32, err error) {
	err = v.Decode(&val)
	return
}

// Float64 将 redis 的原始值转为 float64
func (v Value) Float64() (val float64, err error) {
	err = v.Decode(&val)
	return
}

// Bool 与其他类型转换不同, Bool 直接使用redis原始string值做 ParseBool 操作
func (v Value) Bool() (val bool, err error) {
	s, _ := v.String()
	return strconv.ParseBool(s)
}
