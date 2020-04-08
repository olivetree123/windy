package index

import "github.com/huichen/sego"

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
}

// DataMap 数据字典，用户将索引数据全部加载到内存
var DataMap map[string]*Index

var seg sego.Segmenter

func init() {
	DataMap = make(map[string]*Index)
	seg.LoadDictionary("/Users/gao/code/gowork/src/github.com/huichen/sego/data/dictionary.txt")
}

// SplitWord 分词
func SplitWord(sentence string) []Word {
	var words []Word
	segments := seg.Segment([]byte(sentence))
	rs := sego.SegmentsToSlice(segments, true)
	for _, r := range rs {
		var found bool
		for _, word := range words {
			if word.Content == r {
				word.Count++
				found = true
				break
			}
		}
		if !found {
			word := Word{
				Content: r,
				Count:   1,
			}
			words = append(words, word)
		}
	}
	return words
}
