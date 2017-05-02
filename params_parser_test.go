package ganno_test

import (
	"fmt"
	"testing"

	"github.com/brainicorn/ganno"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParamsParserTestSuite struct {
	suite.Suite
}

// The entry point into the tests
func TestParamsParserTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ParamsParserTestSuite))
}

func (suite *ParamsParserTestSuite) SetupSuite() {
}

type paramsAnno struct {
	attrs map[string][]string
}

func (a *paramsAnno) AnnotationName() string {
	return "paramsAnno"
}

func (a *paramsAnno) Attributes() map[string][]string {
	return a.attrs
}

func (a *paramsAnno) GetAttr(name string) []string {
	attrs, ok := a.attrs[name]
	if !ok {
		attrs = make([]string, 0)
	}

	return attrs
}

func (a *paramsAnno) LenAttr(name string) int {
	attrs, ok := a.attrs[name]
	if !ok {
		attrs = make([]string, 0)
	}

	return len(attrs)
}

type paramsAnnoFactory struct {
}

func (f *paramsAnnoFactory) ValidateAndCreate(name string, attrs map[string][]string) (ganno.Annotation, error) {
	return &paramsAnno{
		attrs: attrs,
	}, nil
}

type erroringAnnoFactory struct {
}

func (f *erroringAnnoFactory) ValidateAndCreate(name string, attrs map[string][]string) (ganno.Annotation, error) {
	return nil, fmt.Errorf("invalid params?")
}

func (suite *ParamsParserTestSuite) TestRegisterBlankName() {
	suite.T().Parallel()

	parser := ganno.NewAnnotationParser()
	err := parser.RegisterFactory("", &paramsAnnoFactory{})

	assert.Error(suite.T(), err)
}

func (suite *ParamsParserTestSuite) TestAlreadyRegistered() {
	suite.T().Parallel()

	parser := ganno.NewAnnotationParser()
	err := parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	assert.NoError(suite.T(), err)

	err = parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	assert.Error(suite.T(), err)
}

func (suite *ParamsParserTestSuite) TestValidationError() {
	suite.T().Parallel()

	input := `@errorAnno()`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("erroranno", &erroringAnnoFactory{})

	annos, errs := parser.Parse(input)

	assert.Equal(suite.T(), 1, len(errs))
	assert.EqualError(suite.T(), errs[0], "invalid params?")
	assert.Equal(suite.T(), 0, len(annos.All()))
}

func (suite *ParamsParserTestSuite) TestSingleParam() {
	suite.T().Parallel()

	input := `@paramsAnno(moo=cow)`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "cow", anno.GetAttr("moo")[0])

}

func (suite *ParamsParserTestSuite) TestBadEqualParam() {
	suite.T().Parallel()

	input := `@paramsAnno(moo!=cow)`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, errs := parser.Parse(input)

	assert.Equal(suite.T(), 1, len(errs))
	assert.EqualError(suite.T(), errs[0], "error parsing parameter key")
	assert.Equal(suite.T(), 0, len(annos.All()))

}

func (suite *ParamsParserTestSuite) TestMultiParam() {
	suite.T().Parallel()

	input := `@paramsAnno(foo="bar",dog="cat")`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 2, len(anno.Attributes()))
	assert.Equal(suite.T(), "bar", anno.GetAttr("foo")[0])
	assert.Equal(suite.T(), "cat", anno.GetAttr("dog")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineParam() {
	suite.T().Parallel()

	input := `@paramsAnno(quack=duck,
 		neigh="horse or mare"
		 )`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 2, len(anno.Attributes()))
	assert.Equal(suite.T(), "duck", anno.GetAttr("quack")[0])
	assert.Equal(suite.T(), "horse or mare", anno.GetAttr("neigh")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineCommentParam() {
	suite.T().Parallel()

	input := `// @paramsAnno(quack=duck,
// 		neigh="horse or mare"
//		 )`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 2, len(anno.Attributes()))
	assert.Equal(suite.T(), "duck", anno.GetAttr("quack")[0])
	assert.Equal(suite.T(), "horse or mare", anno.GetAttr("neigh")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineOpen() {
	suite.T().Parallel()

	input := `@paramsAnno
// 		(a="b")
// `

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "b", anno.GetAttr("a")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineOpenParams() {
	suite.T().Parallel()

	input := `@paramsAnno
// 		(
//		c="d")`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "d", anno.GetAttr("c")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineOpenParamsClose() {
	suite.T().Parallel()

	input := `//
// @paramsanno
// 		(
//		e="f"
//	)
//
//`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "f", anno.GetAttr("e")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineOpenEqualClose() {
	suite.T().Parallel()

	input := `// @ParamsAnno
// 		(	g
//	="h"
//	)
//`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "h", anno.GetAttr("g")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineOpenValueClose() {
	suite.T().Parallel()

	input := `// @paramsAnno
// 		(	i=
//		"j"
//	)`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), "j", anno.GetAttr("i")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineMutliValue() {
	suite.T().Parallel()

	input := `//@paramsanno
// 		(	i=
//		"j"
//	,k="l"
//	)
//`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 2, len(anno.Attributes()))
	assert.Equal(suite.T(), "j", anno.GetAttr("i")[0])
	assert.Equal(suite.T(), "l", anno.GetAttr("k")[0])
}

func (suite *ParamsParserTestSuite) TestMultiMultiLineMutliValueCommaSplit() {
	suite.T().Parallel()

	input := `// @paramsanno
// 		(	i=
//		"j"
//	,
//	m="n"
//	)`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("paramsanno", &paramsAnnoFactory{})

	annos, _ := parser.Parse(input)
	anno := annos.All()[0].(*paramsAnno)
	assert.Equal(suite.T(), "paramsAnno", anno.AnnotationName())
	assert.Equal(suite.T(), 2, len(anno.Attributes()))
	assert.Equal(suite.T(), "j", anno.GetAttr("i")[0])
	assert.Equal(suite.T(), "n", anno.GetAttr("m")[0])
}
