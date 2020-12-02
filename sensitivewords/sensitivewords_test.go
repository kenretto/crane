package sensitivewords

import (
	"testing"
)

func TestNewSensitiveWords(t *testing.T) {
	var store = NewSensitiveWords()
	store.Add("打球")
	store.Add("打球真好玩")
	store.Add("看球赛")
	store.Add("球迷")

	t.Log(store.HasKeywords("想跟我一起打球吗"))
	t.Log(store.HasKeywords("打球可真好玩儿了"))
	t.Log(store.HasKeywords("那跟我一起去看球赛"))
	t.Log(store.HasKeywords("好啊，我可是苍老师的铁杆球迷了"))

	t.Log(store.KeywordsList("想跟我一起打球吗"))
	t.Log(store.KeywordsList("打球可真好玩儿了"))
	t.Log(store.KeywordsList("那跟我一起去看球赛"))
	t.Log(store.KeywordsList("好啊，我可是苍老师的铁杆球迷了"))

	t.Log(store.Filter("想跟我一起打球吗", '*'))
	t.Log(store.Filter("打球可真好玩儿了", '*'))
	t.Log(store.Filter("那跟我一起去看球赛", '*'))
	t.Log(store.Filter("好啊，我可是苍老师的铁杆球迷了", '*'))

	t.Log(store.Filter("那跟我一起去看球赛", '-'))
}
