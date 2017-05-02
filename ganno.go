package ganno

import "strings"

// Annotation is the interface for a single annotation instance. This is what is returned by an
// AnnotationFactory. Consumers can implement their own custom annotations using this interface.
//
// Custom implementations can/should provide strongly typed accessors for their specific attributes.
type Annotation interface {

	// AnnotationName is the name of the annotation between the @ symbol and the (.
	// The value needs to be a valid ident and comparisons are case-insensitive
	AnnotationName() string

	// Attributes holds any key/value pairs inside of the ().
	// The map key is the (loer-cased) key found in the key/val pair. The value is a slice of strings.
	// The value is a slice even if the annotation has a single value to support multi-value k/v pairs.
	Attributes() map[string][]string
}

type basicAnnotation struct {
	AnnoName string              `json:"name"`
	Attrs    map[string][]string `json:"attributes"`
}

func (ba *basicAnnotation) AnnotationName() string {
	return ba.AnnoName
}

func (ba *basicAnnotation) Attributes() map[string][]string {
	return ba.Attrs
}

// Annotations is the interface for holding the collection of annotations found during parsing.
// This interface provides helper methods for retrieving all annotations or annotations by name.
//
// For the most part consumers do not need to implement this unless they implement a custom
// AnnotationParser and have a specific need to augment annotation retrieval.
type Annotations interface {
	// All returns a slice of all found Annotation objects
	All() []Annotation

	// ByName retrieves a slice of Annotation objects whose name matches name. The name comparison
	// should use strings.ToLower before comparing.
	ByName(name string) []Annotation
}

type defaultAnnotations struct {
	all   []Annotation
	named map[string][]Annotation
}

func (da *defaultAnnotations) All() []Annotation {
	return da.all
}

func (da *defaultAnnotations) ByName(name string) []Annotation {
	annos, found := da.named[strings.ToLower(name)]

	if !found {
		annos = make([]Annotation, 0)
	}

	return annos
}

func (da *defaultAnnotations) addAnnotation(anno Annotation) {
	da.all = append(da.all, anno)
	da.named[strings.ToLower(anno.AnnotationName())] = append(da.named[strings.ToLower(anno.AnnotationName())], anno)
}

// AnnotationFactory is the interface consumers can implement to provide custom Annotation creation.
//
// AnnotationFactories can be registered with the parser by name.
type AnnotationFactory interface {
	// ValidateAndCreate is caled by parsers when they find an annotation whose (lower-case) name
	// matches the name of the registered factory. This is the same value passed in as name to this
	// method which allows registering the same factory multiple times under different names.
	//
	// attrs is the map of k/v pairs found when parsing the annotation. This method should use the attrs
	// to validate that the annotation has all expected/required parameters.
	//
	// Any validation error should be returned as error.
	ValidateAndCreate(name string, attrs map[string][]string) (Annotation, error)
}

type basicAnnotationFactory struct {
}

func (af *basicAnnotationFactory) ValidateAndCreate(name string, attrs map[string][]string) (Annotation, error) {
	return &basicAnnotation{
		AnnoName: name,
		Attrs:    attrs,
	}, nil
}
