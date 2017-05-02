package ganno_test

import (
	"testing"

	"github.com/brainicorn/ganno"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SimpleParserTestSuite struct {
	suite.Suite
}

// The entry point into the tests
func TestSimpleParserTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(SimpleParserTestSuite))
}

func (suite *SimpleParserTestSuite) SetupSuite() {
}

func (suite *SimpleParserTestSuite) TestNoAnnos() {
	suite.T().Parallel()

	input := `this is some text.
	it can be anything.
	it

	can

	have line breaks.

	// and even comments

	/*
	and
	multi-line
	comments
	*/
	`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 0, len(annos.All()))

}

func (suite *SimpleParserTestSuite) TestSimpleAnno() {
	suite.T().Parallel()

	input := `@simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())
	assert.Equal(suite.T(), 0, len(annos.All()[0].Attributes()))

}

func (suite *SimpleParserTestSuite) TestSimpleAnnoByName() {
	suite.T().Parallel()

	input := `@simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	simpleAnnos := annos.ByName("simpleAnno")
	assert.Equal(suite.T(), "simpleanno", simpleAnnos[0].AnnotationName())
	assert.Equal(suite.T(), 0, len(simpleAnnos[0].Attributes()))

}

func (suite *SimpleParserTestSuite) TestByNameNotFound() {
	suite.T().Parallel()

	input := `@simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	simpleAnnos := annos.ByName("iDonExist")
	assert.Equal(suite.T(), 0, len(simpleAnnos))

}

func (suite *SimpleParserTestSuite) TestSimpleAnnoTrailingWS() {
	suite.T().Parallel()

	input := `@simpleAnno()   `

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestBadAndGoodAnno() {
	suite.T().Parallel()

	input := `@noParenAnno @aGoodAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 1, len(annos.All()))
	assert.Equal(suite.T(), "agoodanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestSimpleAtBreak() {
	suite.T().Parallel()

	input := `@
	simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestSimpleOpenBreak() {
	suite.T().Parallel()

	input := `@simpleAnno
	()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestSimpleCloseBreak() {
	suite.T().Parallel()

	input := `@simpleAnno(

	)`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}
func (suite *SimpleParserTestSuite) TestSimpleOpenCloseBreak() {
	suite.T().Parallel()

	input := `@simpleAnno
	(

	)`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestSimpleIdentBreak() {
	suite.T().Parallel()

	input := `@simple
	Anno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 0, len(annos.All()))

}

func (suite *SimpleParserTestSuite) TestSimpleAllBreak() {
	suite.T().Parallel()

	input := `@
	simple
	Anno(

	)
	`
	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 0, len(annos.All()))

}

// COMMENTED
func (suite *SimpleParserTestSuite) TestCommentedSimpleAnno() {
	suite.T().Parallel()

	input := `// @simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleAnnoTrailingWS() {
	suite.T().Parallel()

	input := ` // @simpleAnno()   `

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedBadAndGoodAnno() {
	suite.T().Parallel()

	input := `// @noParenAnno @aGoodAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 1, len(annos.All()))
	assert.Equal(suite.T(), "agoodanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleAtBreak() {
	suite.T().Parallel()

	input := `// @
	// simpleAnno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleOpenBreak() {
	suite.T().Parallel()

	input := `// @simpleAnno
	// ()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleCloseBreak() {
	suite.T().Parallel()

	input := `// @simpleAnno(

	// )`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}
func (suite *SimpleParserTestSuite) TestCommentedSimpleOpenCloseBreak() {
	suite.T().Parallel()

	input := `// @simpleAnno
	// (
	//
	// )`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), "simpleanno", annos.All()[0].AnnotationName())

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleIdentBreak() {
	suite.T().Parallel()

	input := `// @simple
	// Anno()`

	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 0, len(annos.All()))

}

func (suite *SimpleParserTestSuite) TestCommentedSimpleAllBreak() {
	suite.T().Parallel()

	input := `// @
	// simple
	// Anno(
	//
	// )
	// `
	parser := ganno.NewAnnotationParser()
	annos, _ := parser.Parse(input)

	assert.Equal(suite.T(), 0, len(annos.All()))

}
