package ganno

import "github.com/brainicorn/goblex"

const (
	beginLineComment      string = "//"
	beginMultiLineComment        = "/*"
	endMultiLineComment          = "*/"
	leftBracket                  = "["
	rightBracket                 = "]"
	atSymbol                     = "@"
	openParen                    = "("
	closeParen                   = ")"
	equalSign                    = "="
	comma                        = ","
	doubleQuote                  = "\""
)

const (
	tokenTypeStartAnno goblex.TokenType = 1 + iota
	tokenTypeKey
	tokenTypeValue
	tokenTypeEndAnno
)

var (
	comments = []string{beginLineComment, beginMultiLineComment, endMultiLineComment}
)
