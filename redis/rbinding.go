package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/go-redis/redis/v8"
	"reflect"
	"time"
)

// RBinding redis go type binding
//  使用 gob 序列化, gob 是 golang 提供的序列化 golang 的一种操作, 但是无法对 interface{} 做序列化操作, 例如你拥有以下结构体,
//   type User struct {
//       Username string
//       FieldAny interface{}
//   }
// 那么实际上使用 gob 对该结构体对象做序列化反序列化操作时都将返回错误, 因为 gob 不会自动去帮助你处理类型问题, 但是 gob 提供了两个方法,
// gob.Register 和 gob.RegisterName , 这两个方法是将你的类型注册给 gob 全局的一个操作, 建议使用 gob.RegisterName , 虽然两者作用相差不大,
// gob.RegisterName 内部使用的 sync.Map 管理注册的类型, 所以其实可以不用太关系在什么地方注册的问题以及并发的问题, 问题都不大, 所以你可以选择在控制器
// 或者任何地方去注册, 但是记得在你使用序列化操作和反序列化操作之前注册进去了进行,
// 需要注意的是, gob.RegisterName 不会接受同一个 name 前后注册的类型不一样这种问题, 比如 gob.RegisterName("a", A{}); gob.RegisterName("a", B{}),
// 这样操作的话将会得到异常
type RBinding struct {
	client *Redis
	log    ILogger
	ctx    context.Context
}

// NewDefaultRBinding 获取一个默认的 RBinding 操作绑定
func NewDefaultRBinding(ctx context.Context, r *Redis) *RBinding {
	return &RBinding{
		client: r,
		log:    r.logger,
		ctx:    ctx,
	}
}

// Set redis set
func (r *RBinding) Set(key string, val interface{}, expiration time.Duration) (err error) {
	var buffer = new(bytes.Buffer)
	err = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(val))
	r.client.Instance().Set(r.ctx, key, buffer.Bytes(), expiration)
	return
}

// Get redis get
func (r *RBinding) Get(key string, val interface{}) (err error) {
	var result = r.client.Instance().Get(r.ctx, key)
	data, err := result.Bytes()
	if err != nil {
		return err
	}

	err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(reflect.ValueOf(val))
	return
}

// GetSet redis getset
//  将 val 的写入到 redis , 并将 redis 之前此 key 的值绑定到 val, val 必须是一个引用
func (r *RBinding) GetSet(key string, val interface{}) (err error) {
	var buffer = new(bytes.Buffer)
	var rValue = reflect.ValueOf(val)
	err = gob.NewEncoder(buffer).EncodeValue(rValue)
	if err != nil {
		return err
	}

	result := r.client.Instance().GetSet(r.ctx, key, buffer.Bytes())
	data, err := result.Bytes()
	if err != nil {
		if err == redis.Nil {
			if rValue.Elem().CanSet() {
				rValue.Elem().Set(reflect.Zero(rValue.Elem().Type()))
			}
		}
		return err
	}

	var oldValue = reflect.New(rValue.Elem().Type())
	err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(oldValue)
	if rValue.Elem().CanSet() {
		rValue.Elem().Set(oldValue.Elem())
	}
	return
}

// Del redis del
func (r *RBinding) Del(key string) {
	r.client.Instance().Del(r.ctx, key)
}

// SetNX set nx
func (r *RBinding) SetNX(key string, val interface{}, expiration time.Duration) (rs bool, err error) {
	var buffer = new(bytes.Buffer)
	err = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(val))
	if err != nil {
		return false, err
	}

	setNx := r.client.Instance().SetNX(r.ctx, key, buffer.Bytes(), expiration)
	rs, err = setNx.Result()
	return
}

// LoadOrStorage 获取指定key的值,如果值不存在,就执行f方法将返回值存入redis
func (r *RBinding) LoadOrStorage(key string, val interface{}, f func() (expiration time.Duration, data interface{})) error {
	var v = reflect.ValueOf(val)
	err := r.Get(key, val)
	if err != nil && err == redis.Nil {
		expiration, val := f()
		v.Elem().Set(reflect.ValueOf(val))
		return r.Set(key, val, expiration)
	}

	if r.client == nil || r.client.Instance().Ping(r.ctx).Val() != "PONG" {
		_, val := f()
		v.Elem().Set(reflect.ValueOf(val))
		r.log.Warn("redis client lose connection or client not initiate")
	}

	return err
}

// MGet redis mGet
func (r *RBinding) MGet(keys ...string) map[string]Value {
	rs, err := r.client.Instance().MGet(r.ctx, keys...).Result()
	if err != nil {
		r.log.Warn(err)
	}

	var data = make(map[string]Value)
	for i := range rs {
		var val Value
		if rs[i] == nil {
			val = Value{decoder: nil}
		} else {
			val = Value{decoder: gob.NewDecoder(bytes.NewReader([]byte(rs[i].(string))))}
		}
		data[keys[i]] = val
	}

	return data
}

// List 返回一个指定 key 的 list 操作
//  返回的 push 和 pop 操作, 需要保证 val 为同一类型
func (r *RBinding) List(key string) *List {
	return &List{client: r.client.Instance(), key: key, ctx: r.ctx}
}

// Sets redis set 结构
func (r *RBinding) Sets(key string) *Set {
	return &Set{client: r.client.Instance(), key: key, logger: r.log, ctx: r.ctx}
}

// List redis list 数据结构
type List struct {
	client redis.Cmdable
	key    string
	ctx    context.Context
}

// Len list len
func (l *List) Len() int64 {
	return l.client.LLen(l.ctx, l.key).Val()
}

// LPush redis Left Push
func (l *List) LPush(values ...interface{}) {
	var val = make([]interface{}, len(values))
	for i := range values {
		var buffer = new(bytes.Buffer)
		_ = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(values[i]))
		val[i] = buffer.Bytes()
	}
	l.client.LPush(l.ctx, l.key, val...)
}

// LPop redis Left pop
func (l *List) LPop(val interface{}) error {
	data, err := l.client.LPop(l.ctx, l.key).Bytes()
	if err != nil {
		return err
	}

	err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(reflect.ValueOf(val))
	return err
}

// RPush redis right Push
func (l *List) RPush(values ...interface{}) {
	var val = make([]interface{}, len(values))
	for i := range values {
		var buffer = new(bytes.Buffer)
		_ = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(values[i]))
		val[i] = buffer.Bytes()
	}
	l.client.RPush(l.ctx, l.key, val...)
}

// RPop redis right pop
func (l *List) RPop(val interface{}) error {
	data, err := l.client.RPop(l.ctx, l.key).Bytes()
	if err != nil {
		return err
	}

	err = gob.NewDecoder(bytes.NewReader(data)).DecodeValue(reflect.ValueOf(val))
	return err
}

// Set redis set 操作
type Set struct {
	client redis.Cmdable
	logger ILogger
	key    string
	ctx    context.Context
}

// SAdd 向集合添加一个或多个成员
func (set *Set) SAdd(members ...interface{}) {
	var val = make([]interface{}, len(members))
	for i := range members {
		var buffer = new(bytes.Buffer)
		_ = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(members[i]))
		val[i] = buffer.Bytes()
	}

	set.client.SAdd(set.ctx, set.key, val...)
}

// SCard 向集合添加一个或多个成员
func (set *Set) SCard() int64 {
	return set.client.SCard(set.ctx, set.key).Val()
}

// SDiff 返回第一个集合与其他集合之间的差异。
//  container 是一个结构, 比如这个 set 的值是 string, container 就是 string, 值是struct,container就传struct, 不需要是个引用
func (set *Set) SDiff(container interface{}, diffKeys ...string) interface{} {
	var keys []string
	keys = append(keys, set.key)
	keys = append(keys, diffKeys...)

	var values = set.client.SDiff(set.ctx, keys...).Val()
	var typ = reflect.TypeOf(container)

	newInstance := reflect.MakeSlice(reflect.SliceOf(typ), len(values), len(values))
	items := reflect.New(newInstance.Type())
	items.Elem().Set(newInstance)

	var val = reflect.New(typ)
	for i := range values {
		err := gob.NewDecoder(bytes.NewReader([]byte(values[i]))).DecodeValue(val)
		if err != nil {
			set.logger.Warn(err)
		}
		items.Elem().Index(i).Set(val.Elem())
	}

	return items.Elem().Interface()
}

// SInter 返回给定所有集合的交集
func (set *Set) SInter(container interface{}, diffKeys ...string) interface{} {
	var keys []string
	keys = append(keys, set.key)
	keys = append(keys, diffKeys...)

	var values = set.client.SInter(set.ctx, keys...).Val()
	var typ = reflect.TypeOf(container)

	newInstance := reflect.MakeSlice(reflect.SliceOf(typ), len(values), len(values))
	items := reflect.New(newInstance.Type())
	items.Elem().Set(newInstance)

	var val = reflect.New(typ)
	for i := range values {
		err := gob.NewDecoder(bytes.NewReader([]byte(values[i]))).DecodeValue(val)
		if err != nil {
			set.logger.Warn(err)
		}
		items.Elem().Index(i).Set(val.Elem())
	}

	return items.Elem().Interface()
}

// SMembers 返回集合中的所有成员
func (set *Set) SMembers(container interface{}) interface{} {
	var values = set.client.SMembers(set.ctx, set.key).Val()
	var typ = reflect.TypeOf(container)

	newInstance := reflect.MakeSlice(reflect.SliceOf(typ), len(values), len(values))
	items := reflect.New(newInstance.Type())
	items.Elem().Set(newInstance)

	var val = reflect.New(typ)
	for i := range values {
		err := gob.NewDecoder(bytes.NewReader([]byte(values[i]))).DecodeValue(val)
		if err != nil {
			set.logger.Warn(err)
		}
		items.Elem().Index(i).Set(val.Elem())
	}

	return items.Elem().Interface()
}

// SMove 将 member 元素从 source 集合移动到 destination 集合
func (set *Set) SMove(destination string, member interface{}) error {
	var buffer = new(bytes.Buffer)
	_ = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(member))
	return set.client.SMove(set.ctx, set.key, destination, buffer.Bytes()).Err()
}

// SPop 移除并返回集合中的一个或多个随机元素
func (set *Set) SPop(container interface{}, n ...int64) interface{} {
	var count int64 = 1
	if len(n) == 1 {
		count = n[1]
	}

	values := set.client.SPopN(set.ctx, set.key, count).Val()
	var typ = reflect.TypeOf(container)

	newInstance := reflect.MakeSlice(reflect.SliceOf(typ), len(values), len(values))
	items := reflect.New(newInstance.Type())
	items.Elem().Set(newInstance)

	var val = reflect.New(typ)
	for i := range values {
		err := gob.NewDecoder(bytes.NewReader([]byte(values[i]))).DecodeValue(val)
		if err != nil {
			set.logger.Warn(err)
		}
		items.Elem().Index(i).Set(val.Elem())
	}

	return items.Elem().Interface()
}

// SRem 移除集合中一个或多个成员
func (set *Set) SRem(members ...interface{}) int64 {
	var val = make([]interface{}, len(members))
	for i := range members {
		var buffer = new(bytes.Buffer)
		_ = gob.NewEncoder(buffer).EncodeValue(reflect.ValueOf(members[i]))
		val[i] = buffer.Bytes()
	}
	return set.client.SRem(set.ctx, set.key, val...).Val()
}

// SUnion 返回所有给定集合的并集
func (set *Set) SUnion(container interface{}, diffKeys ...string) interface{} {
	var keys []string
	keys = append(keys, set.key)
	keys = append(keys, diffKeys...)

	var values = set.client.SUnion(set.ctx, keys...).Val()
	var typ = reflect.TypeOf(container)

	newInstance := reflect.MakeSlice(reflect.SliceOf(typ), len(values), len(values))
	items := reflect.New(newInstance.Type())
	items.Elem().Set(newInstance)

	var val = reflect.New(typ)
	for i := range values {
		err := gob.NewDecoder(bytes.NewReader([]byte(values[i]))).DecodeValue(val)
		if err != nil {
			set.logger.Warn(err)
		}
		items.Elem().Index(i).Set(val.Elem())
	}

	return items.Elem().Interface()
}
