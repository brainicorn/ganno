package ganno

import "github.com/brainicorn/goblex"

// LexBegin is the entry point LexFn for lexing java style annotations.
// Parsers should pass this function as the begin parameter when calling goblex.NewLexer
func LexBegin(lexer *goblex.Lexer) goblex.LexFn {

	if lexer.CaptureUntil(true, atSymbol) {
		lexer.SkipCurrentToken(true)
		return lexAtSymbol
	}

	return nil
}

func lexAtSymbol(lexer *goblex.Lexer) goblex.LexFn {
	if lexer.CaptureIdent() {
		if lexer.CurrentTokenIs(openParen) {
			lexer.Emit(tokenTypeStartAnno)
			return lexOpenParen
		}
	}

	return LexBegin
}

func lexOpenParen(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, openParen)
	lexer.SkipCurrentToken(true)

	if lexer.CurrentTokenIs(closeParen) {
		return lexCloseParen
	}

	return lexKey
}

func lexCloseParen(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, closeParen)
	lexer.SkipCurrentToken(true)
	lexer.Emit(tokenTypeEndAnno)
	return LexBegin
}

func lexKey(lexer *goblex.Lexer) goblex.LexFn {

	if lexer.CaptureIdent() {
		if lexer.CurrentTokenIs(equalSign) {
			lexer.Emit(tokenTypeKey)
			return lexEqualSign
		}
	}

	lexer.Errorf("error parsing parameter key")
	return LexBegin
}

func lexEqualSign(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, equalSign)
	lexer.SkipCurrentToken(true)

	return lexValue
}

func lexValue(lexer *goblex.Lexer) goblex.LexFn {
	if lexer.CurrentTokenIs(leftBracket) {
		return lexLeftBracket
	}

	return lexSingleValue
}

func lexSingleValue(lexer *goblex.Lexer) goblex.LexFn {
	if lexer.CurrentTokenIs(doubleQuote) {
		return lexSingleQuotedValue
	}

	if tkn := lexer.CaptureUntilOneOf(true, comma, closeParen); tkn != "" {
		lexer.Emit(tokenTypeValue)

		switch tkn {
		case comma:
			return lexSingleValueComma

		case closeParen:
			return lexCloseParen
		}
	}

	lexer.Errorf("error parsing single value: comma or close paren missing")
	return LexBegin
}

func lexMultiValue(lexer *goblex.Lexer) goblex.LexFn {

	if lexer.CurrentTokenIs(doubleQuote) {
		return lexMultiQuotedValue
	}

	if tkn := lexer.CaptureUntilOneOf(true, comma, rightBracket); tkn != "" {
		lexer.Emit(tokenTypeValue)

		switch tkn {
		case comma:
			return lexMultiValueComma

		case rightBracket:
			return lexRightBracket
		}
	}

	lexer.Errorf("error multi value: comma or rbracket missing")
	return LexBegin
}

func lexSingleQuotedValue(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, doubleQuote)
	lexer.SkipCurrentToken(true)
	lexer.RemoveIgnoreTokens(comments...)
	if lexer.CaptureUntil(false, doubleQuote) {
		lexer.Emit(tokenTypeValue)
		lexer.AddIgnoreTokens(comments...)
		lexer.SkipCurrentToken(true)
		yup, tkn := lexer.CurrentTokenIsOneOf(comma, closeParen)

		if yup {
			switch tkn {
			case comma:
				return lexSingleValueComma

			case closeParen:
				return lexCloseParen
			}
		}
	}

	lexer.AddIgnoreTokens(comments...)
	lexer.Errorf("error parsing single quoted value: comma or close paren missing")
	return LexBegin
}

func lexMultiQuotedValue(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, doubleQuote)
	lexer.SkipCurrentToken(true)
	lexer.RemoveIgnoreTokens(comments...)
	if lexer.CaptureUntil(false, doubleQuote) {
		lexer.Emit(tokenTypeValue)
		lexer.AddIgnoreTokens(comments...)
		lexer.SkipCurrentToken(true)

		yup, tkn := lexer.CurrentTokenIsOneOf(comma, rightBracket)

		if yup {
			switch tkn {
			case comma:
				return lexMultiValueComma

			case rightBracket:
				return lexRightBracket
			}
		}
	}

	lexer.AddIgnoreTokens(comments...)
	lexer.Errorf("error parsing multi quoted value: comma or Rbracket missing")
	return LexBegin
}

func lexSingleValueComma(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, comma)
	lexer.SkipCurrentToken(true)
	return lexKey
}

func lexMultiValueComma(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, comma)
	lexer.SkipCurrentToken(true)
	return lexMultiValue
}

func lexLeftBracket(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, leftBracket)
	lexer.SkipCurrentToken(true)

	return lexMultiValue
}

func lexRightBracket(lexer *goblex.Lexer) goblex.LexFn {
	lexer.CaptureUntil(true, rightBracket)
	lexer.SkipCurrentToken(true)

	yup, tkn := lexer.CurrentTokenIsOneOf(comma, closeParen)

	if yup {
		switch tkn {
		case comma:
			return lexSingleValueComma

		case closeParen:
			return lexCloseParen
		}
	}

	lexer.Errorf("error parsing array value")
	return LexBegin
}
