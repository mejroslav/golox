package golox

type BreakValue struct {
	Keyword *Token
}

func NewBreakValue(keyword *Token) *BreakValue {
	return &BreakValue{
		Keyword: keyword,
	}
}

func (bv *BreakValue) Error() string {
	return "Break statement encountered"
}
