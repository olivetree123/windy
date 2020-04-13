package entity

type QueryParam struct {
	Name  string
	Value string
}

//type Field struct {
//	Name  string
//	Value string
//}

type DocumentScore struct {
	DocID string
	Score float32
}

func NewDocumentScore(docID string, score float32) *DocumentScore {
	return &DocumentScore{
		DocID: docID,
		Score: score,
	}
}
