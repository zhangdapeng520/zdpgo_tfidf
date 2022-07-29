package seg

type Tokenizer interface {
	Seg(text string) []string
	Free()
}
