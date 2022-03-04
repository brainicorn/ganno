![Build Status](https://github.com/brainicorn/ganno/actions/workflows/build.yaml/badge.svg)
[![codecov](https://codecov.io/gh/brainicorn/ganno/branch/main/graph/badge.svg)](https://codecov.io/gh/brainicorn/ganno)
[![Go Report Card](https://goreportcard.com/badge/github.com/brainicorn/ganno)](https://goreportcard.com/report/github.com/brainicorn/ganno)
[![GoDoc](https://godoc.org/github.com/brainicorn/ganno?status.svg)](https://godoc.org/github.com/brainicorn/ganno)

# ganno

Package ganno implements java-style annotations in Go.

If your unfamiliar with java's annotation system, you can look at their [annotations basics guide](https://docs.oracle.com/javase/tutorial/java/annotations/basics.html)

This library includes a lexer (based on [the goblex lexer library](https://github.com/brainicorn/goblex)) and a pluggable
annotation parser.

**API Documentation:** [https://godoc.org/github.com/brainicorn/ganno](https://godoc.org/github.com/brainicorn/ganno)

[Issue Tracker](https://github.com/brainicorn/ganno/issues)

### Java-Style Annotations

Java-Style anotations are written using the following format:

`@<ident>(<ident>=<valueOrValues>,...)`

examples:

```go

@simpleAnnotation()

@simplewithParam(mykey=myval)

@quotedVal(somekey="my value")

@multipleParams(magic="wizards", awesome="unicorns")

@multipleVals(mypets=["dog", "kitty cat"])

```

Annotations may also be split across multiple lines even within single-line comments:

```go

// @stuffILike(
// 	instrument="drums"
// 	,mypets=[
// 		"dog"
// 		,"kitty cat"
// 	]
// 	,food="nutritional units"
// )

```

## Features

- Highly tested
- Ready to use out of the box
- Anotations can be parsed from:
  - single line text
  - multi-line text
  - multi-line comment blocks
  - single-line comments
  - multiple single-line comment blocks
- Runtime pluggable
  - Supports custom Annotation types
    - Custom annotations can provide strongly-typed sttribute accessors
  - Custom annotation factories can be registered at runtime
  - Factories can enforce strict validation of attributes
  - Factories and custom Annotations can be re-used/distributed as libraries
- Extensible
  - Provides interfaces for implementing custom AnnotationParsers and or Annotations collections

## Basic Use

Out of the box, ganno is ready to use, but it's good to understand a couple of things:

- The parser returns an Annotations object as well as any validation errors while parsing
  - Validation errors are returned in a slice and the annotation that errored is discarded
  - The Annotations object provides accessors for All() annotations as well as ByName(name)
- The default annotation object returned provides an Attributes() method which returns a
  `map[string][]string` where the map key is the attribute name and the value is a slice of strings. This
  is to support multi-value atrributes and single-value attributes are a slice of 1.

Here's a simple example of parsing an @pets annotation:

```go
input := `my @pet(name="fluffy buns", hasFur=true) is soooo cute!`

parser := ganno.NewAnnotationParser()

annos, errs := parser.Parse(input)

//there were some validation errors, we could inspect them, but let's just panic the first one
if len(errs) > 0 {
	panic(errs[0])
}

//since we know there's only one...
ourPet := annos.All()[0]

// get our pet's name
var name := "unknown"
if nameSlice, nameOk := ourPet.Attributes()["name"]; nameOk {
	name = nameSlice[0]
}

var hasFur := false
// get our pet's fur status
furSlice, furOk := ourPet.Attributes()["hasfur"]; if furOk {
	if hasFur, err := strconv.ParseBool(furSlice[0]); err != nil {
		hasFur = false
	}
}

fmt.Printf("%s is fluffy? %t\n", name, hasFur)
// Output: fluffy buns is fluffy? true

```

While the above example certainly works, there's a lot of validation that's going on that's not very
reusable and is error prone. Let's see how we can use a plugin to solve this...

## Advanced Use

By creating a custom annotation type and a factory for it, we can encapsulate all of the above validation
logic and provide a cleaner interface when dealing with @pet annotations...

First off, let's create the "plugin" bits that could be put into their own package and/or library if
desired....

```go
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
```

With all of that in place, our dealings with @pet annotations is simplfied:

```go
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
```

**API Documentation:** [https://godoc.org/github.com/brainicorn/ganno](https://godoc.org/github.com/brainicorn/ganno)

[Issue Tracker](https://github.com/brainicorn/ganno/issues)

## Contributors

Pull requests, issues and comments welcome. For pull requests:

- Add tests for new features and bug fixes
- Follow the existing style
- Separate unrelated changes into multiple pull requests

See the existing issues for things to start contributing.

For bigger changes, make sure you start a discussion first by creating
an issue and explaining the intended change.

## License

Apache 2.0 licensed, see [LICENSE.txt](LICENSE.txt) file.
