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

type SearchResult struct {
	Data     []*DocumentEntity `json:"data"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Total    int               `json:"total"`
}

func NewDocumentScore(docID string, score float32) *DocumentScore {
	return &DocumentScore{
		DocID: docID,
		Score: score,
	}
}

func NewSearchResult(data []*DocumentEntity, page int, pageSize int, total int) *SearchResult {
	r := SearchResult{
		Data:     data,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	return &r
}
