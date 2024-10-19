[![Go Report Card](https://goreportcard.com/badge/github.com/jose78/go-collections)](https://goreportcard.com/report/github.com/jose78/go-collections)
[![Coverage Status](https://coveralls.io/repos/github/jose78/go-collections/badge.svg?branch=master)](https://coveralls.io/github/jose78/go-collections?branch=master)
[![Go-Collections](https://github.com/jose78/go-collections/actions/workflows/go_collections.yml/badge.svg)](https://github.com/jose78/go-collections/actions/workflows/go_collections.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/jose78/go-collections/v1.svg)](https://pkg.go.dev/github.com/jose78/go-collections/)


# Go-collections <img align="right" width="80" height="100" src="resources/gopher.png">




The `go-collections` package provides a set of utility functions for working with collections (slices and maps) in Go. It includes functions for mapping, filtering, zipping, and performing actions on collections.

Installation
------------

To install the package, use:

```sh
    go get github.com/jose78/go-collections
```

Usage
-----


### Filter

 Filter filters elements of a slice or array based on a predicate and stores the result in dest.

 The `Filter` function takes a predicate, a source (which must be a pointer to a slice or array), and a destination
 (which must also be a pointer to a slice or array). It filters the elements of the source that satisfy the predicate
 and stores them in the destination.

 Parameters:
 - predicate: A function that determines whether an element should be included in the destination.
 - source: A pointer to a slice or array to be filtered. It must be a pointer type.
 - dest: A pointer to a slice or array where the filtered elements will be stored. It must be a pointer type.

 Returns:
 - An error if the source or dest is not a pointer to a slice or array, otherwise returns nil.

 Example usage:
```go
numbers := []int{1, 2, 3, 4, 5}
var evens []int

predicate := func(n int) bool {
    return n%2 == 0
}

err := Filter(predicate, &numbers, &evens)
if err != nil {
    log.Fatal(err)
}

fmt.Println(evens) // Output: [2 4]
```
 Note:
 The `Filter` function does not modify the original slice or array.



### ForEach

 ForEach applies an action to each element of a slice or array.

 The `ForEach` function takes an action and a source (which must be a pointer to a slice or array) and applies
 the action to each element of the source.

 Parameters:
 - action: A function that performs an operation on each element of the source.
 - source: A pointer to a slice or array on which the action will be performed. It must be a pointer type.

 Returns:
 - An error if the source is not a pointer to a slice or array, otherwise returns nil.

 Example usage:

```go
people := []string{"Alice", "Bob", "Charlie"}

action := func(name string) {
    fmt.Println("Hello, " + name)
}

err := ForEach(action, &people)
if err != nil {
    log.Fatal(err)
}
```
 Note:
 The `ForEach` function does not modify the original slice or array.




### GroupBy

 GroupBy groups elements of a slice or array based on a key selector function and stores the result in dest.

 The `GroupBy` function takes a key selector function, a source (which must be a pointer to a slice or array), and a
 destination (which must be either a pointer to a map or a slice). It groups the elements of the source based on the
 keys returned by the key selector function and stores the results in the destination.

 If the destination is a map, the function adds a tuple (key, value) where the key is the result of the key selector
 and the value is the corresponding element from the source.

 Parameters:
 - keySelector: A function that selects a key for each element of the source.
 - source: A pointer to a slice or array to be grouped. It must be a pointer type.
 - dest: A pointer to a map or slice where the grouped elements will be stored. It must be a pointer type.

 Returns:
 - An error if the source or dest is not a pointer to a slice, array, or map, otherwise returns nil.

 Example usage:
```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 30},
}

var groups map[int][]Person

keySelector := func(p Person) int {
    return p.Age
}

err := GroupBy(keySelector, &people, &groups)
if err != nil {
    log.Fatal(err)
}

fmt.Println(groups) // Output: map[25:[{Bob 25}] 30:[{Alice 30} {Charlie 30}]]
```
 Note:
 The `GroupBy` function does not modify the original slice or array.

### Map

 Map applies a mapper function to each element of a slice or array and stores the result in dest.

 The `Map` function takes a mapper function, a source (which must be a pointer to a slice or array), and a destination
 (which must also be a pointer to a slice or array). It applies the mapper function to each element of the source and
 stores the results in the destination.

 Parameters:
 - mapper: A function that transforms an element of the source into an element of the destination.
 - source: A pointer to a slice or array to be transformed. It must be a pointer type.
 - dest: A pointer to a slice or array where the transformed elements will be stored. It must be a pointer type.

 Returns:
 - An error if the source or dest is not a pointer to a slice or array, otherwise returns nil.

 Example usage:
```go
numbers := []int{1, 2, 3, 4, 5}
var squares []int

mapper := func(n int) int {
    return n * n
}

err := Map(mapper, &numbers, &squares)
if err != nil {
    log.Fatal(err)
}

fmt.Println(squares) // Output: [1 4 9 16 25]
```
 Note:
 The `Map` function does not modify the original slice or array.

### SortBy

 SortBy sorts a slice or array based on a provided comparator.

 The `SortBy` function takes a comparator and a source (which must be a pointer to a slice or array) 
 and sorts the elements in the source according to the comparator. The comparator defines the ordering 
 of the elements.

 Parameters:
 - comparator: A function that defines the order of the elements. It should return a negative value 
   if the first argument is less than the second, zero if they are equal, and a positive value if 
   the first argument is greater than the second.
 - source: A pointer to a slice or array to be sorted. It must be a pointer type.

 Returns:
 - An error if the source is not a pointer to a slice or array, otherwise returns nil.

Example usage:
```go
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
}

comparator := func(a, b Person) int {
    return a.Age - b.Age
}

err := SortBy(comparator, &people)
if err != nil {
    log.Fatal(err)
}

fmt.Println(people) // Output: [{Bob 25} {Alice 30} {Charlie 35}]
```
 Note:
 The `SortBy` function modifies the original slice or array.


### ZIP

 Zip combines two slices into a map.

 The `Zip` function takes two slices, `keys` and `values`, and a destination map `dest`.
 It populates the `dest` map with the `keys` and `values` such that each key in the `keys` slice
 maps to the corresponding value in the `values` slice.

 If the number of keys and values are not the same, the function will return an error indicating
 that there is a mismatch in the number of elements.

 Parameters:
 - keys: A slice of keys to be used in the resulting map. Each key must be comparable.
 - values: A slice of values corresponding to the keys. Can be of any type.
 - dest: A map that will be populated with the keys and values. It must be initialized before
   calling this function.

 Returns:
 - An error if there is a mismatch in the length of `keys` and `values`, otherwise returns nil.

 Example usage:

```go
keys := []string{"a", "b", "c"}
values := []int{1, 2, 3}
dest := make(map[string]int)
err := Zip(keys, values, dest)
if err != nil {
    log.Fatal(err)
}

fmt.Println(dest) // Output: map[a:1 b:2 c:3]
```
 Note:
 The `Zip` function will overwrite any existing keys in the `dest` map with new values.