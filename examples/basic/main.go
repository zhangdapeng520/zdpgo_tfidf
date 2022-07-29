package main

import (
	"fmt"
	"github.com/zhangdapeng520/zdpgo_tfidf"
	"github.com/zhangdapeng520/zdpgo_tfidf/seg"
	"github.com/zhangdapeng520/zdpgo_tfidf/similarity"
)

func main() {

	f := zdpgo_tfidf.New()

	// 计算两个英文的相似度
	f.AddDocs("how are you", "are you fine", "how old are you", "are you ok", "i am ok", "i am file")
	t1 := "it is so cool"
	w1 := f.Cal(t1)
	fmt.Printf("weight of %s is %+v.\n", t1, w1)
	t2 := "you are so beautiful"
	w2 := f.Cal(t2)
	fmt.Printf("weight of %s is %+v.\n", t2, w2)
	sim := similarity.Cosine(w1, w2)
	fmt.Printf("cosine between %s and %s is %f .\n", t1, t2, sim)
	fmt.Println("===================")

	// 创建jieba分词对象
	tokenizer := seg.NewJieba()
	defer tokenizer.Free()

	// 创建序列化器
	f = zdpgo_tfidf.NewTokenizer(tokenizer)

	// 添加文档
	f.AddDocs("重庆大学", "上海市复旦大学", "上海交通大学", "重庆理工大学")

	// 计算权重
	t1 = "重庆市西南大学"
	w1 = f.Cal(t1)
	fmt.Printf("weight of %s is %+v.\n", t1, w1)

	// 计算权重
	t2 = "重庆市重庆大学"
	w2 = f.Cal(t2)
	fmt.Printf("weight of %s is %+v.\n", t2, w2)

	// 计算相似度
	sim = similarity.Cosine(w1, w2)
	fmt.Printf("cosine between %s and %s is %f .\n", t1, t2, sim)
}
