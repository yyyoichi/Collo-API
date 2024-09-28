package tfidf

// build with Builder
type TFIDF struct {
	// 単語識別子が配列dataのいずれの列Xに位置するか
	xindex map[int]int
	// 文書識別子が配列dataのいずれの行Yに位置するか
	yindex map[string]int

	// tfidf行列
	data []float64
}
