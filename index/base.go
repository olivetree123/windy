package index

import (
	"github.com/olivetree123/sego"
)

// MatchDoc 匹配的文档
type MatchDoc struct {
	DocID    string
	StartPos int
	EndPos   int
}

// Index 索引
type Index struct {
	Frequency int
	Property  string
	MatchDocs []*MatchDoc
}

type Word struct {
	Value string
	Count int
}

// DataMap 数据字典，用户将索引数据全部加载到内存
var DataMap map[string]*Index

var seg sego.Segmenter

// wordsFreq 记录每个词的词频
var wordsFreq map[string]int

func NewWord(value string, count int) Word {
	word := Word{
		Value: value,
		Count: count,
	}
	return word
}

func init() {
	wordsFreq = make(map[string]int, 590000)
	DataMap = make(map[string]*Index)
	seg.LoadDictionary("/var/lib/windy/dictionary.txt")
	for _, token := range seg.Dictionary().GetTokens() {
		wordsFreq[token.Text()] = token.Frequency()
	}
}

// SplitWord 分词
func SplitWord(sentence string) []Word {
	var words []Word
	data := make(map[string]int)
	segments := seg.Segment([]byte(sentence))
	ws := sego.SegmentsToSlice(segments, true)
	for _, w := range ws {
		if _, found := data[w]; found {
			data[w]++
		} else {
			data[w] = 1
		}
	}
	for key, value := range data {
		word := NewWord(key, value)
		words = append(words, word)
	}
	return words
}

func GetWordFreq(word string) int {
	return wordsFreq[word]
}
