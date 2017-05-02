package ganno_test

import (
	"testing"

	"github.com/brainicorn/ganno"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MultivalParserTestSuite struct {
	suite.Suite
}

// The entry point into the tests
func TestMultivalParserTestSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(MultivalParserTestSuite))
}

func (suite *MultivalParserTestSuite) SetupSuite() {
}

type mvAnno struct {
	attrs map[string][]string
}

func (a *mvAnno) AnnotationName() string {
	return "multiVal"
}

func (a *mvAnno) Attributes() map[string][]string {
	return a.attrs
}

func (a *mvAnno) GetAttr(name string) []string {
	attrs, ok := a.attrs[name]
	if !ok {
		attrs = make([]string, 0)
	}

	return attrs
}

func (a *mvAnno) LenAttr(name string) int {
	attrs, ok := a.attrs[name]
	if !ok {
		attrs = make([]string, 0)
	}

	return len(attrs)
}

type mvAnnoFactory struct {
}

func (f *mvAnnoFactory) ValidateAndCreate(name string, attrs map[string][]string) (ganno.Annotation, error) {
	return &mvAnno{
		attrs: attrs,
	}, nil
}

func (suite *MultivalParserTestSuite) TestMVSimple() {
	suite.T().Parallel()

	input := `@multiVal(pets=[dog,cat])`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "cat", pets[1])
}

func (suite *MultivalParserTestSuite) TestMVLBreak() {
	suite.T().Parallel()

	input := `// @multiVal(pets=[
	// dog,cat])`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "cat", pets[1])
}

func (suite *MultivalParserTestSuite) TestMVLValBreak() {
	suite.T().Parallel()

	input := `// @multiVal(pets=[
	// dog
	,cat])`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "cat", pets[1])
}

func (suite *MultivalParserTestSuite) TestMVLValRBreak() {
	suite.T().Parallel()

	input := `// @multiVal(pets=[
	// dog
	// ,cat
	// ])`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "cat", pets[1])
}

func (suite *MultivalParserTestSuite) TestMVBreakAll() {
	suite.T().Parallel()

	input := `// @multiVal(
	// pets=[
	// 		"dog"
	// 		,"cat"
	// 	]
	// )`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 1, len(anno.Attributes()))
	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "cat", pets[1])
}

func (suite *MultivalParserTestSuite) TestMVMixed() {
	suite.T().Parallel()

	input := `// @multiVal(
	// dwelling="house"
	// ,pets=[
	// 		"dog"
	// 		,"kitty cat"
	// 	]
	// ,vehicle="car"
	// )`

	parser := ganno.NewAnnotationParser()
	parser.RegisterFactory("multival", &mvAnnoFactory{})

	annos, _ := parser.Parse(input)

	anno := annos.All()[0].(*mvAnno)
	pets := anno.GetAttr("pets")
	dwelling := anno.GetAttr("dwelling")
	vehicle := anno.GetAttr("vehicle")

	assert.Equal(suite.T(), "multiVal", anno.AnnotationName())
	assert.Equal(suite.T(), 3, len(anno.Attributes()))

	assert.Equal(suite.T(), 2, len(pets))
	assert.Equal(suite.T(), "dog", pets[0])
	assert.Equal(suite.T(), "kitty cat", pets[1])

	assert.Equal(suite.T(), 1, len(dwelling))
	assert.Equal(suite.T(), "house", dwelling[0])

	assert.Equal(suite.T(), 1, len(vehicle))
	assert.Equal(suite.T(), "car", vehicle[0])
}
