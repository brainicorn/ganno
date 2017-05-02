//Package ganno implements java-style annotations in Go.
//This library includes a lexer (based on https://github.com/brainicorn/goblex) and a pluggable
//annotation parser.
//
//Annotations can be written in the following formats:
//	@simpleAnnotation()
//
//	@simplewithParam(mykey=myval)
//
//	@quotedVal(somekey="my value")
//
//	@multipleParams(magic="wizards", awesome="unicorns")
//
//	@multipleVals(mypets=["dog", "kitty cat"])
//
//Annotations may also be split across multiple lines:
//	@stuffILike(
//		instrument="drums"
//		,mypets=[
//			"dog"
//			,"kitty cat"
//		]
//		,food="nutritional units"
//	)
//
//Although this library can be used out of the box as-is, consumer can also "plug-in" custom annotation
//types that can validate discovered annotations and provide strongly-type attribute accessors.
//
//Custom Annotations and their factories can be re-used/distributed as libraries and registered with
//parsers at runtime.
//
//see example below
package ganno
