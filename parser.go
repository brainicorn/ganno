package ganno

import (
	"fmt"
	"strings"

	"github.com/brainicorn/goblex"
)

// AnnotationParser is the interface for parsing a string and returning Annotations.
// The returned Annnotations object should be populated with Annotation objects created using the
// registered AnnotationFactory objects.
//
// Consumers of this library shouldn't need to implement this interface unless they have very specific
// custom parsing needs.
type AnnotationParser interface {
	// RegisterFactory registers an AnnotationFactory with the supplied name. The name will be
	// lower-case compared with the names of discovered annotations to choose the proper factory for
	// creation.
	//
	// If name is blank or a factory with the same name has already been registered an error will be
	// retured.
	RegisterFactory(name string, factory AnnotationFactory) error

	// Parse lexes all of the tokens in the input string and returns an Annotations object which holds
	// all of the valid discovered annotations. The annotations may be of various types/implementations
	// based on the registered factories and can be cast to those specific types.
	//
	// If a validation error is returned by the factory during creation, the annotation will not be
	// added to the Annotations and the error will be put into the returned errors slice.
	Parse(input string) (Annotations, []error)
}

type defaultAnnotationParser struct {
	registry map[string]AnnotationFactory
}

// NewAnnotationParser creates an AnnotationParser that can be used to discover annotations.
func NewAnnotationParser() AnnotationParser {
	return &defaultAnnotationParser{
		registry: make(map[string]AnnotationFactory),
	}
}

// RegisterFactory implements AnnotationParser
func (p *defaultAnnotationParser) RegisterFactory(name string, factory AnnotationFactory) error {
	factoryName := strings.ToLower(strings.TrimSpace(name))
	if factoryName == "" {
		return fmt.Errorf("cannot register annotation factory with blank name")
	}

	_, exists := p.registry[factoryName]
	if exists {
		return fmt.Errorf("annotation factory with name '%s' is already registered", factoryName)
	}

	p.registry[factoryName] = factory

	return nil
}

// Parse implements AnnotationParser
func (p *defaultAnnotationParser) Parse(input string) (Annotations, []error) {
	output := &defaultAnnotations{named: make(map[string][]Annotation)}
	var errs = make([]error, 0)
	var token goblex.Token

	currentAttrs := make(map[string][]string)
	currentParamKey := ""
	currentAnnoName := ""

	l := goblex.NewLexer("someFile", input, LexBegin)
	l.AddIgnoreTokens(comments...)
	// l.Debug = true
	for {

		if l.IsEOF() {
			break
		}

		token = l.NextEmittedToken()

		switch token.Type() {

		case tokenTypeStartAnno:
			currentAnnoName = strings.ToLower(strings.TrimSpace(token.String()))
			currentAttrs = make(map[string][]string)
			currentParamKey = ""

		case tokenTypeValue:
			currentAttrs[currentParamKey] = append(currentAttrs[currentParamKey], token.String())

		case tokenTypeKey:
			currentParamKey = strings.ToLower(strings.TrimSpace(token.String()))

		case tokenTypeEndAnno:
			factory, found := p.registry[currentAnnoName]

			if !found {
				factory = &basicAnnotationFactory{}
			}

			anno, err := factory.ValidateAndCreate(currentAnnoName, currentAttrs)

			if err == nil {
				output.addAnnotation(anno)
			} else {
				errs = append(errs, err)
			}

			currentAttrs = make(map[string][]string)
			currentParamKey = ""

		case goblex.TokenTypeError:
			fmt.Println("error: ", token.String())
			errs = append(errs, fmt.Errorf("%s", token.String()))
		}
	}

	return output, errs

}
