package collections

import (
	"errors"
)

// GenerateMapEmpty create an empty MapType
func GenerateMapEmpty() MapType {
	result := MapType{}
	return result
}

// GenerateMap is the default item
func GenerateMap(a, b interface{}) MapType {
	result := MapType{}
	result[a] = b
	return result
}

// GenerateMapFromTuples is the default item
func GenerateMapFromTuples(tuples ListType) MapType {
	result := MapType{}
	for _, item := range tuples {
		tuple := item.(Tuple)
		result[tuple.a] = tuple.b
	}
	return result
}

// GenerateMapFromZip is the default item
func GenerateMapFromZip(keys, values []interface{}) MapType {
	tuples, _ := Zip(keys, values)
	return GenerateMapFromTuples(tuples)
}

//Foreach is the default
func (mapType MapType) Foreach(fn func(interface{}, interface{}, int)) {
	index := 0
	for key, value := range mapType {
		fn(key, value, index)
		index++
	}
}

func callbackMapTypeMap(index int, key, value interface{}, fnInternal MapperMapType) (item interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {

			// find out exactly what the error was and set err
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	item = fnInternal(key, value, index)
	return item, err
}


// MapperMapType define the function to be used in Map function
type MapperMapType func(interface{}, interface{}, int) interface{}

//Map is the default
func (mapType MapType) Map(fn MapperMapType) (ListType , error){
	result := ListType{}
	index := 0
	for key, value := range mapType {
		if item, err := callbackMapTypeMap(index, key, value, fn); err != nil{
			return nil, err
		}else{
			result = append(result, item)
		}
		index++
	}	
	return result, nil
}

//FilterAll is the default
func (mapType MapType) FilterAll(fn func(interface{}, interface{}) bool) MapType {
	result := GenerateMapEmpty()
	for key, value := range mapType {
		if fn(key, value) {
			result[key] = value
		}
	}
	return result
}

//ListValues obtains a ListType of the values in this Map.
func (mapType MapType) ListValues() ListType {
	list := ListType{}
	for _, value := range mapType {
		list = append(list, value)
	}
	return list
}

//ListKeys obtains a ListType of the keys in this Map.
func (mapType MapType) ListKeys() ListType {
	list := ListType{}
	for key, _ := range mapType {
		list = append(list, key)
	}
	return list
}
