package index

import (
	"github.com/huichen/sego"
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
	Content string
	Count   int
	Freq    int
}

// DataMap 数据字典，用户将索引数据全部加载到内存
var DataMap map[string]*Index

var seg sego.Segmenter

// wordsFreq 记录每个词的词频
var wordsFreq map[string]int

func init() {
	wordsFreq = make(map[string]int, 590000)
	DataMap = make(map[string]*Index)
	seg.LoadDictionary("/Users/gao/code/gowork/src/github.com/huichen/sego/data/dictionary.txt")
	for _, token := range seg.Dictionary().GetTokens() {
		wordsFreq[token.Text()] = token.Frequency()
	}
}

// SplitWord 分词
func SplitWord(sentence string) []Word {
	var words []Word
	segments := seg.Segment([]byte(sentence))
	//rs := sego.SegmentsToSlice(segments, true)
	for _, r := range segments {
		var found bool
		for _, word := range words {
			if word.Content == r.Token().Text() {
				word.Count++
				found = true
				break
			}
		}
		if !found {
			word := Word{
				Content: r.Token().Text(),
				Count:   1,
				Freq:    r.Token().Frequency(),
			}
			words = append(words, word)
		}
	}
	return words
}

func GetWordFreq(word string) int {
	return wordsFreq[word]
}
