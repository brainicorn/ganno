package ganno

import "github.com/brainicorn/goblex"

const (
	beginLineComment      string = "//"
	beginMultiLineComment string = "/*"
	endMultiLineComment   string = "*/"
	leftBracket           string = "["
	rightBracket          string = "]"
	atSymbol              string = "@"
	openParen             string = "("
	closeParen            string = ")"
	equalSign             string = "="
	comma                 string = ","
	doubleQuote           string = "\""
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
