package ganno

import (
	"fmt"
	"strconv"
)

// PetAnno is a custom Annotation type for @pet() annotations
type PetAnno struct {
	Attrs map[string][]string
}

func (a *PetAnno) AnnotationName() string {
	return "pet"
}

func (a *PetAnno) Attributes() map[string][]string {
	return a.Attrs
}

// HasFur is a strongly-typed accessor for the "hasfur" attribute
func (a *PetAnno) Hasfur() bool {
	b, _ := strconv.ParseBool(a.Attrs["hasfur"][0])

	return b
}

// Name is an accessor for the "name" attribute
func (a *PetAnno) Name() string {
	return a.Attrs["name"][0]
}

// PetAnnoFactory is our custom factory that can validate @pet attributes and return new PetAnnos
type PetAnnoFactory struct{}

func (f *PetAnnoFactory) ValidateAndCreate(name string, attrs map[string][]string) (Annotation, error) {
	furs, furOk := attrs["hasfur"]
	_, nameOk := attrs["name"]

	// require the hasfur attr
	if !furOk {
		return nil, fmt.Errorf("pet annotation requires the attribute %q", "hasfur")
	}

	// require the name attr
	if !nameOk {
		return nil, fmt.Errorf("pet annotation requires the attribute %q", "name")
	}

	// make sure hasfur is a parsable boolean
	if _, err := strconv.ParseBool(furs[0]); err != nil {
		return nil, fmt.Errorf("pet annotation attribute %q must have a boolean value", "hasfur")
	}

	return &PetAnno{
		Attrs: attrs,
	}, nil
}

func Example() {

	parser := NewAnnotationParser()

	//register our factory
	parser.RegisterFactory("pet", &PetAnnoFactory{})

	input := `my @pet(name="fluffy buns", hasFur=true) is soooo cute!`

	annos, errs := parser.Parse(input)

	if len(errs) > 0 {
		panic(errs[0])
	}

	// get our one pet annotation and cast it to a PetAnno
	mypet := annos.ByName("pet")[0].(*PetAnno)

	fmt.Printf("%s is fluffy? %t\n", mypet.Name(), mypet.Hasfur())
	// Output: fluffy buns is fluffy? true
}
