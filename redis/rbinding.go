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
//  Using gob serialization, gob is an operation provided by golang to serialize golang, but interface{} cannot be serialized. For example, you have the following structure,
//   type User struct {
//       Username string
//       FieldAny interface{}
//   }
// In fact, when using gob to serialize and deserialize the structure object, an error will be returned, because gob will not automatically help you deal with type problems, but gob provides two methods,
// gob.Register and gob.RegisterName , These two methods are an operation to register your type to gob global, which is recommended gob.RegisterName  Although there is little difference between them,
// gob.RegisterName internal use sync.Map manage the type of registration, So in fact, it doesn't have to do with where to register and concurrency,
// So you can choose to register in the controller or anywhere, but remember to register before you use serialization and deserialization,
// It should be noted that, gob.RegisterName will not accept the problem that the type of registration before and after the same name is different, such as gob.RegisterName("a", A{}); gob.RegisterName("a", B{}),
// If you do this, will panic
type RBinding struct {
	client *Redis
	log    ILogger
	ctx    context.Context
}

// NewDefaultRBinding Gets a default binding for the rbinding operation
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
