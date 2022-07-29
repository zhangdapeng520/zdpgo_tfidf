package zdpgo_tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math"

	"github.com/zhangdapeng520/zdpgo_tfidf/seg"
	"github.com/zhangdapeng520/zdpgo_tfidf/util"
)

// TFIDF tfidf model
type TFIDF struct {
	docIndex  map[string]int         // 训练文档索引在术语频次
	termFreqs []map[string]int       // 每个train文件的术语频率
	termDocs  map[string]int         // train数据中每个术语的文件编号
	n         int                    // train数据中的文件数量
	stopWords map[string]interface{} // 过滤词语
	tokenizer seg.Tokenizer          // tokenizer, 使用空格作为默认值
}

// New 使用默认参数新建模型
func New() *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: &seg.EnTokenizer{},
	}
}

// NewTokenizer 使用指定序列化器新建模型
func NewTokenizer(tokenizer seg.Tokenizer) *TFIDF {
	return &TFIDF{
		docIndex:  make(map[string]int),
		termFreqs: make([]map[string]int, 0),
		termDocs:  make(map[string]int),
		n:         0,
		tokenizer: tokenizer,
	}
}

// initStopWords 初始化停用词
func (f *TFIDF) initStopWords() {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	lines, err := util.ReadLines("../data/stopword")
	if err != nil {
		log.Printf("init stop words with error: %s", err)
	}

	for _, w := range lines {
		f.stopWords[w] = nil
	}
}

// AddStopWords 添加要过滤的停用词
func (f *TFIDF) AddStopWords(words ...string) {
	if f.stopWords == nil {
		f.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		f.stopWords[word] = nil
	}
}

// AddStopWordsFile 添加要过滤的停用词文件，一行一个停用词
func (f *TFIDF) AddStopWordsFile(file string) (err error) {
	lines, err := util.ReadLines(file)
	if err != nil {
		return
	}

	f.AddStopWords(lines...)
	return
}

// AddDocs 添加要train训练的文档
func (f *TFIDF) AddDocs(docs ...string) {
	for _, doc := range docs {
		h := hash(doc)
		if f.docHashPos(h) >= 0 {
			return
		}

		termFreq := f.termFreq(doc)
		if len(termFreq) == 0 {
			return
		}

		f.docIndex[h] = f.n
		f.n++

		f.termFreqs = append(f.termFreqs, termFreq)

		for term := range termFreq {
			f.termDocs[term]++
		}
	}
}

// Cal 计算指定文档的 tf-idf 权重
func (f *TFIDF) Cal(doc string) (weight map[string]float64) {
	weight = make(map[string]float64)

	var termFreq map[string]int

	docPos := f.docPos(doc)
	if docPos < 0 {
		termFreq = f.termFreq(doc)
	} else {
		termFreq = f.termFreqs[docPos]
	}

	docTerms := 0
	for _, freq := range termFreq {
		docTerms += freq
	}
	for term, freq := range termFreq {
		weight[term] = tfidf(freq, docTerms, f.termDocs[term], f.n)
	}

	return weight
}

func (f *TFIDF) termFreq(doc string) (m map[string]int) {
	m = make(map[string]int)

	tokens := f.tokenizer.Seg(doc)
	if len(tokens) == 0 {
		return
	}

	for _, term := range tokens {
		if _, ok := f.stopWords[term]; ok {
			continue
		}

		m[term]++
	}

	return
}

func (f *TFIDF) docHashPos(hash string) int {
	if pos, ok := f.docIndex[hash]; ok {
		return pos
	}

	return -1
}

func (f *TFIDF) docPos(doc string) int {
	return f.docHashPos(hash(doc))
}

func hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func tfidf(termFreq, docTerms, termDocs, N int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(1+N) / (1 + float64(termDocs)))
	return tf * idf
}
