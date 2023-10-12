package patcher

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrPatcherInvalidType is an error that returns when the patcher receives an invalid type
	ErrPatcherInvalidType = errors.New("patcher: invalid type")
)

// Patch handles the patching of a map[string]any into a struct via reflection
// - keys in the patch are mapped to the struct via field name or field tag
// - key found - not mapped cases:
//  - the field is not exported
//  - the field type is not assignable
func Patch(ptr any, patch map[string]any) (err error) {
	// rValue is the reflect value of the pointer
	rValue:= reflect.ValueOf(ptr).Elem()
	// rType is the reflect type of the pointer
	rType := reflect.TypeOf(ptr).Elem()


	// iterate over the rValue
	for i := 0; i < rValue.NumField(); i++ {
		// get the field info
		fieldName := rType.Field(i).Name
		fieldType := rType.Field(i)
		fieldTag  := fieldType.Tag.Get("patcher")

		// get the field value
		fieldValue := rValue.FieldByName(fieldName)
		
		// if the field is not exported, skip it
		if !fieldValue.CanSet() {
			continue
		}

		// check if the field is in the patch
		var exists bool; var key string
		for k := range patch {
			if k == fieldTag || k == fieldName {
				exists = true
				key = k
				break
			}
		}
		if !exists {
			continue
		}

		// get the patch value
		patchValue := patch[key]

		// check assignable case
		assignableCase := checkCase(patchValue, fieldValue.Interface())
		if assignableCase == NotAssignable {
			err = fmt.Errorf("%w - fieldName %s - fieldTag %s", ErrPatcherInvalidType, fieldName, fieldTag)
			return
		}

		// handle the assignable case
		switch assignableCase {
		case AssignableJSON:
			// convert the patch value to int
			patchValue = int(patchValue.(float64))
			fieldValue.Set(reflect.ValueOf(patchValue))
		default:
			fieldValue.Set(reflect.ValueOf(patchValue))
		}
	}

	return
}

// Assignament is a type that represents the assignable case
type Assignament int
const (
	NotAssignable Assignament = iota
	Assignable
	AssignableJSON
	// ...
)

// checkCase is a function that checks the assignable case
// - By default, the assignable case is NotAssignable
func checkCase(patchValue any, floatValue any) (a Assignament) {
	// check if the patch value is assignable to the field value
	if reflect.TypeOf(patchValue).AssignableTo(reflect.TypeOf(floatValue)) {
		a = Assignable
	}

	// check if is assignable by JSON
	// check if field type is int and patch value is float64
	if reflect.TypeOf(patchValue).Kind() == reflect.Float64 && reflect.TypeOf(floatValue).Kind() == reflect.Int {
		// check if the patch value does not have decimals
		if patchValue.(float64) == float64(int(patchValue.(float64))) {
			// set the assignable case
			a = AssignableJSON
		}
	}

	return
}