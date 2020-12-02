package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/kenretto/crane/configurator"
	"testing"
	"time"
)

type user struct {
	Username string
	Age      uint8
	Avatar   []byte
	Password struct {
		Salt  []byte
		Token []byte
	}
}

func TestRBinding(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	t.Log(NewDefaultRBinding(context.Background(), r).Set("test_go", "123456789", 0))

	var out string
	t.Log(NewDefaultRBinding(context.Background(), r).Get("test_go", &out))
	t.Log(out)

	t.Log(NewDefaultRBinding(context.Background(), r).Set("test_go", user{
		Username: "medivh",
		Age:      25,
		Avatar:   []byte("https://www.baidu.com"),
		Password: struct {
			Salt  []byte
			Token []byte
		}{Salt: []byte("123321"), Token: []byte("abc")},
	}, 0))

	var u user
	t.Log(NewDefaultRBinding(context.Background(), r).Get("test_go", &u))
	t.Log(u)
	t.Log(u)
}

func TestCacheGet(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	var ticker = time.NewTicker(1 * time.Second)
	var age = 25
	for {
		<-ticker.C
		var u user
		err := NewDefaultRBinding(context.Background(), r).LoadOrStorage("test_user", &u, func() (expiration time.Duration, data interface{}) {
			return 5 * time.Second, user{
				Username: "medivh",
				Age:      uint8(age),
				Avatar:   []byte("https://www.baidu.com"),
				Password: struct {
					Salt  []byte
					Token []byte
				}{Salt: []byte("123321"), Token: []byte("abc")},
			}
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(u)
		age++
		if age == 35 {
			break
		}
	}
}

func TestRBinding_MGet(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	_ = NewDefaultRBinding(context.Background(), r).Set("nickname", "medivh", 0)
	_ = NewDefaultRBinding(context.Background(), r).Set("age", 25, 0)
	_ = NewDefaultRBinding(context.Background(), r).Set("money", 12.34, 0)
	_ = NewDefaultRBinding(context.Background(), r).Set("user", user{
		Username: "medivh",
		Age:      uint8(123),
		Avatar:   []byte("https://www.baidu.com"),
		Password: struct {
			Salt  []byte
			Token []byte
		}{Salt: []byte("123321"), Token: []byte("abc")},
	}, 0)

	t.Log(NewDefaultRBinding(context.Background(), r).MGet("nickname", "age", "money")["age"].Int64())
	t.Log(NewDefaultRBinding(context.Background(), r).MGet("nickname", "age", "money")["nickname"].String())
	t.Log(NewDefaultRBinding(context.Background(), r).MGet("nickname", "age", "money")["money"].Float64())
	t.Log(NewDefaultRBinding(context.Background(), r).MGet("nickname", "age", "money", "null")["null"].Float64())

	var u user
	_ = NewDefaultRBinding(context.Background(), r).MGet("nickname", "age", "money", "user")["user"].Decode(&u)
	t.Log(u)
	t.Log(string(u.Avatar))
}

func TestRBinding_GetSet(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	var u = user{
		Username: "medivh2",
		Age:      uint8(6),
		Avatar:   []byte("https://www.3.com"),
		Password: struct {
			Salt  []byte
			Token []byte
		}{Salt: []byte("3"), Token: []byte("3")},
	}

	_ = NewDefaultRBinding(context.Background(), r).GetSet("user", &u)
	t.Log(u)
}

func TestRBinding_List(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	NewDefaultRBinding(context.Background(), r).List("list").LPush("a", "b", "c", "d", "e")
	NewDefaultRBinding(context.Background(), r).List("list").LPush("f")
	NewDefaultRBinding(context.Background(), r).List("list").LPush("g")

	var v string
	for {
		err := NewDefaultRBinding(context.Background(), r).List("list").RPop(&v)
		if err != nil {
			if err != redis.Nil {
				t.Error(err)
			}
			break
		}
		t.Log(v)
	}

	NewDefaultRBinding(context.Background(), r).List("members").LPush(user{Username: "aa", Age: 1}, user{Username: "bb", Age: 2}, user{Username: "cc", Age: 3}, user{Username: "dd", Age: 4})
	NewDefaultRBinding(context.Background(), r).List("members").LPush(user{Username: "ee", Age: 5})
	var u user
	for {
		err := NewDefaultRBinding(context.Background(), r).List("members").RPop(&u)
		if err != nil {
			if err != redis.Nil {
				t.Error(err)
			}
			break
		}
		t.Log(u)
	}

}

func TestRBinding_Set(t *testing.T) {
	var r = new(Redis)
	var c, err = configurator.NewConfigurator("testdata/redis.yaml")
	if err != nil {
		t.Log(err)
	}
	c.Add("redis", r)
	t.Log(r.Config)
	NewDefaultRBinding(context.Background(), r).Sets("set_1").SAdd("a", "b", "c", "e")
	NewDefaultRBinding(context.Background(), r).Sets("set_2").SAdd("a", "c")
	var s string
	t.Log(NewDefaultRBinding(context.Background(), r).Sets("set_1").SDiff(s, "set_2"))

	NewDefaultRBinding(context.Background(), r).Sets("set_3").SAdd(user{Username: "aa", Age: 1}, user{Username: "bb", Age: 2}, user{Username: "cc", Age: 3}, user{Username: "dd", Age: 4})
	NewDefaultRBinding(context.Background(), r).Sets("set_4").SAdd(user{Username: "aa", Age: 1}, user{Username: "cc", Age: 3})
	t.Log(NewDefaultRBinding(context.Background(), r).Sets("set_3").SDiff(user{}, "set_4"))
}
