package golox

type FunctionType int

const (
	FT_NONE FunctionType = iota
	FT_FUNCTION
	FT_INITIALIZER
	FT_METHOD
)
