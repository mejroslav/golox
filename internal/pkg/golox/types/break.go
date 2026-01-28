package types

import "github.com/mejroslav/golox/internal/pkg/golox/token"

type BreakValue struct {
	Keyword *token.Token
}

func NewBreakValue(keyword *token.Token) *BreakValue {
	return &BreakValue{
		Keyword: keyword,
	}
}

func (bv *BreakValue) Error() string {
	return "Break statement encountered"
}
