# Acme text editor JSON value inspector

TL;DR: `cat example.json | go run .`

This Acme plugin reads a JSON value either from standard input or from the first command line argument, if provided. It shows an interactive window representing that value. For arrays, indicies are shown. For objects, the keys are shown. Right clicking on them opens a window with the representation of that part of the value. File path represents the location of the value relative to the value that was given initially. Type of the value is provided in the tag as well, to disambiguate strings from other JSON data types.
