package collections

import (
        "errors"
        "fmt"
        "reflect"
        "sort"
)



// Touple represents a key-value pair with a generic key and value.
// The key must be of a comparable type, allowing it to be used in maps or
// other dest structures that require comparison operations.
// The value can be of any type.
//
// K - the type of the key, which must be comparable.
// V - the type of the value, which can be any type.
type Touple struct {
        Key   any // Key is the key of the key-value pair.
        Value any // Value is the value of the key-value pair.
}

// builderError struct with an error and the item that caused the error
type builderError[T any] struct {
        err   error
        index int
        item  T
}

// Method to retrieve the error from the builder
func (b *builderError[K]) Error() error {
        return b.err
}

type ErrorFormatter[T any] func(int, T) (err error) // New type for error formatting

// Method to set a custom error message in the builder using a function
func (b *builderError[T]) WithErrorMessage(fn ErrorFormatter[T]) *builderError[T] {
        if fn != nil {
                b.err = fn(b.index, b.item)
        }
        return b
}

// Action is a function type that takes an index and a value of type T.
// This function is used to perform an action on each element in a collection.
type Action[T any] func(int, T)

// ForEach applies the action function to each element in the source collection.
// Parameters:
//   - action: a function that takes an index and a value of type T and performs an action.
//   - source: the collection of elements to iterate over. Must be a slice or array.
//
// Returns:
//   - error: an error if the source is not of the appropriate type or if any other problem occurs during the operation.
func ForEach[K any](action Action[K], src any) (err error) {
        builderError := iterate(action, src)
        if builderError != nil {
                currentError := builderError.Error()
                errorFormatter := func(index int, item K) error {
                        return fmt.Errorf("error processing item %v at index %d: %w", item, index, currentError)
                }

                err = builderError.WithErrorMessage(errorFormatter).Error()
        }
        return
}

func iterate[K any](action Action[K], src any) *builderError[K] {
        var errBuilder *builderError[K]
        evaluate := func(index int, internaParam any) {
                defer func(index int, item any) {
                        if err := recover(); err != nil {
                                valueParametrized := item.(K)
                                errBuilder = &builderError[K]{
                                        item:  valueParametrized,
                                        index: index,
                                        err:   err.(error),
                                }
                        }
                }(index, internaParam)
                action(index, internaParam.(K))
        }
        if IsMap(src) {
                val := reflect.ValueOf(src)
                count := -1
                for _, key := range val.MapKeys() {
                        count++
                        value := val.MapIndex(key)
                        touple := Touple{key.Interface(), value.Interface()}
                        evaluate(count, touple)
                }
        } else {
                for index, item := range src.([]K) {
                        if !reflect.ValueOf(errBuilder).IsZero() {
                                break
                        }
                        if !reflect.ValueOf(errBuilder).IsNil() {
                                break
                        }
                        evaluate(index, item)
                }
        }
        return errBuilder
}

// Zip combines two slices into a map, using elements from the keys slice as keys and elements from the values slice as values.
// K - the type of the keys, which must be comparable.
// V - the type of the values.
// keys - the slice of keys.
// values - the slice of values.
// result - the map where the keys and values are combined.
// Returns an error if the operation fails, such as when the lengths of keys and values do not match.
func Zip[K comparable, V any](keys []K, values []V, result map[K]V) *builderError[K] {
        b := &builderError[K]{}
        if len(keys) != len(values) {
                b.err = errors.New("keys and values slices must have the same length")
                return b
        }
        for i, key := range keys {
                defer func(k K, v V) {
                        if r := recover(); r != nil {
                                b.err = fmt.Errorf("error zipping item: %v", r)
                                b.item = k
                        }
                }(key, values[i])
                result[key] = values[i]
        }
        return b
}

// IsMap checks if the given element is of map type.
// elements - the element to check.
// Returns true if the element is a map, false otherwise.
func IsMap(elements any) bool {
        t := reflect.TypeOf(elements)
        return reflect.Map == t.Kind()
}

// store inserts data into the destination collection, which can be either a map or a slice.
// data - the data to be inserted. If dest is a map, data should be of type Touple with Key and Value fields.
// dest - the destination collection where the data will be stored; should be a map or a pointer to a slice.
// If dest is a map, data.(Touple).Key is used as the key and data.(Touple).Value is used as the value.
// If dest is a slice, data is appended to the slice.
func store(data any, dest any) {
        if IsMap(dest) {
                val := reflect.ValueOf(dest)
                keyVal := reflect.ValueOf(data.(Touple).Key)
                valueVal := reflect.ValueOf(data.(Touple).Value)
                if val.MapIndex(keyVal).IsValid() {
                        existingValue := val.MapIndex(keyVal)
                        if existingValue.Kind() == reflect.Slice && valueVal.Kind() == reflect.Slice {
                                mergedValue := reflect.AppendSlice(existingValue, valueVal)
                                val.SetMapIndex(keyVal, mergedValue)
                        } else {
                                val.SetMapIndex(keyVal, valueVal)
                        }
                } else {
                        val.SetMapIndex(keyVal, valueVal)
                }
        } else {
                sliceVal := reflect.ValueOf(dest).Elem()
                elemVal := reflect.ValueOf(data)
                result := reflect.Append(sliceVal, elemVal)
                sliceVal.Set(result)
        }
}

// Predicate is a function type that takes a value of type T and returns a boolean.
// This function is used to test whether an input value satisfies a condition.
// If the destination is of type map, the input value must be of type Tuple.
type Predicate[T any] func(T) bool

// Filter filters elements from the source using a predicate function and stores the results in the destination.
// The destination must be a list (slice or array).
// Parameters:
//   - predicate: a function that takes a value of type T and returns a boolean indicating whether the value satisfies the condition.
//   - source: the collection of elements to be filtered.
//   - dest: the destination where the results will be stored. Must be a list (slice or array).
//
// Returns:
//   - error: an error if the destination is not of the appropriate type or if any other problem occurs during the operation.
func Filter[T any](predicate Predicate[T], source any, dest any) (err error) {
        var action Action[T] = func(index int, item T) {
                if predicate(item) {
                        store(item, dest)
                }
        }

        builderError := iterate(action, source)
        if builderError != nil {
                currentError := builderError.Error()
                errorFormatter := func(index int, item T) error {
                        return fmt.Errorf("error processing item %v at index %d: %w", item, index, currentError)
                }

                err = builderError.WithErrorMessage(errorFormatter).Error()
        }
        return
}

// Mapper is a function type that takes a value of type T and returns a value of any type.
// This function is used to transform an input value.
// If the destination is of type map, the input value must be of type Tuple.
type Mapper[T any] func(T) any

// Map applies the mapper function to each element in the source and stores the results in the destination.
// The destination must be either a map or a list (slice or array).
// Parameters:
//   - mapper: a function that takes a value of type T and returns a transformed value.
//   - source: the collection of elements to be mapped.
//   - dest: the destination where the results will be stored. Must be a map or a list (slice or array).
//
// Returns:
//   - error: an error if the destination is not of the appropriate type or if any other problem occurs during the operation.
func Map[T any](mapper Mapper[T], source any, dest any) (err error) {
        var action Action[T] = func(index int, item T) {
                result := mapper(item)
                store(result, dest)
        }

        builderError := iterate(action, source)
        if builderError != nil {
                currentError := builderError.Error()
                errorFormatter := func(index int, item T) error {
                        return fmt.Errorf("error processing item %v at index %d: %w", item, index, currentError)
                }

                err = builderError.WithErrorMessage(errorFormatter).Error()
        }
        return
}

// KeySelector is a function type that takes a value of type K and returns a value of any type.
// This function is used to select a key from an input value.
// If the destination is of type map, the input value must be of type Tuple.
type KeySelector[K any] func(K) any

// GroupBy groups elements from the source using a key selection function (keySelector).
// The results are stored in the destination (dest), which can only be of type map or a pointer to a list (slice or array).
// Parameters:
//   - keySelector: a function that takes a value of type T and returns a grouping key.
//   - source: the collection of elements to be grouped.
//   - dest: the destination where the results will be stored. Must be a map or a pointer to a list (slice or array).
//
// Returns:
//   - error: an error if the destination is not of the appropriate type or if any other problem occurs during the operation.
func GroupBy[T any](keySelector KeySelector[T], source any, dest any) (err error) {
        var action Action[T] = func(index int, item T) {
                result := keySelector(item)
                touple := Touple{result, []T{item}}
                store(touple, dest)
        }

        builderError := iterate(action, source)
        if builderError != nil {
                currentError := builderError.Error()
                errorFormatter := func(index int, item T) error {
                        return fmt.Errorf("error processing item %v at index %d: %w", item, index, currentError)
                }

                err = builderError.WithErrorMessage(errorFormatter).Error()
        }
        return
}

// Comparator is a function type that takes two values of type T and returns an integer.
// The return value should be negative if the first value is less than the second,
// zero if they are equal, and positive if the first value is greater than the second.
type Comparator[T any] func(T, T) int

// SortBy sorts the elements in the source using the provided comparator function.
// The source must be a pointer to a list (array or slice).
// Parameters:
//   - comparator: a function that takes two values of type T and returns an integer
//     indicating their order.
//   - source: a pointer to the list (array or slice) to be sorted.
//
// Returns:
//   - error: an error if the source is not of the appropriate type or if any other problem occurs during the operation.
func SortBy[T any](comparator Comparator[T], source any) error {

        if !IsListUpdatable(source) {
                return fmt.Errorf("the provided source is not an updatable list (pointer to list): %v", source)
        }

        lst := source.(*[]T)
        sort.Slice(*lst, func(i, j int) bool {
                return comparator((*lst)[i], (*lst)[j]) < 0
        })

        return nil
}

func IsListUpdatable(source any) bool {
        sourceType := reflect.TypeOf(source)

        return reflect.Ptr == sourceType.Kind() &&
                (reflect.Slice == sourceType.Elem().Kind() ||
                        reflect.Array == sourceType.Elem().Kind())
}
