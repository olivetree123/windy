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

// DataMap 数据字典，用户将索引数据全部加载到内存
var DataMap map[string]*Index

var seg sego.Segmenter

func init() {
	DataMap = make(map[string]*Index)
	seg.LoadDictionary("/Users/gao/code/gowork/src/github.com/huichen/sego/data/dictionary.txt")
}

// SplitWord 分词
func SplitWord(sentence string) []string {
	segments := seg.Segment([]byte(sentence))
	words := sego.SegmentsToSlice(segments, true)
	return words
}
