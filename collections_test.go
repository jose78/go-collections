package gocollection

import (
	"fmt"
	"reflect"
	"testing"
)

type testUser struct {
	name       string
	secondName string
	mails      []string
	age        int
	male       bool
}

func generateTestCaseList() []testUser {
	return []testUser{{name: "John", secondName: "Connor", mails: []string{}, age: 10, male: true}, {name: "Sarah", secondName: "Connor", mails: []string{}, age: 43}, {name: "Kyle", secondName: "Risk", male: true, mails: []string{}, age: 43}}
}

func generateTestCaseMap() map[string]testUser {
	result := map[string]testUser{}
	for _, item := range generateTestCaseList() {
		result[item.name] = item
	}
	return result
}

func errorFmt(user testUser) error {
	return fmt.Errorf("KO")
}

// Test functions
func TestMap(t *testing.T) {
	source := []int{1, 2, 3, 4}
	var dest []any
	builder := Map(func(n int) any { return n * 2 }, source, &dest)

	if err := builder.Error(); err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	expected := []any{2, 4, 6, 8}
	for i, v := range dest {
		if v != expected[i] {
			t.Errorf("Map result mismatch at index %d: got %v, want %v", i, v, expected[i])
		}
	}
}

func TestFilter(t *testing.T) {
	source := []int{1, 2, 3, 4}
	var dest []int
	builder := Filter(func(n int) bool { return n%2 == 0 }, source, &dest)

	if err := builder.Error(); err != nil {
		t.Fatalf("Filter failed: %v", err)
	}

	expected := []int{2, 4}
	for i, v := range dest {
		if v != expected[i] {
			t.Errorf("Filter result mismatch at index %d: got %v, want %v", i, v, expected[i])
		}
	}
}

func TestZip(t *testing.T) {
	keys := []string{"a", "b", "c"}
	values := []int{1, 2, 3}
	result := make(map[string]int)
	builder := Zip(keys, values, result)

	if err := builder.Error(); err != nil {
		t.Fatalf("Zip failed: %v", err)
	}

	expected := map[string]int{"a": 1, "b": 2, "c": 3}
	for k, v := range expected {
		if result[k] != v {
			t.Errorf("Zip result mismatch for key %v: got %v, want %v", k, result[k], v)
		}
	}
}

func TestWithErrorMessage(t *testing.T) {
	source := []int{1, 2, 3, 4}
	var dest []any

	mapper := func(n int) any {
		if n == 3 {
			panic("unexpected value")
		}
		return n * 2
	}
	builder := Map(mapper, source, &dest)

	customErrFunc := func(item int) error {
		return fmt.Errorf("custom error for item %d", item)
	}

	err := builder.WithErrorMessage(customErrFunc).Error()
	if err == nil || err.Error() != "custom error for item 3" {
		t.Fatalf("expected custom error for item 3, got %v", err)
	}
}

func TestForEach(t *testing.T) {

	var actionOk Action[testUser] = func(i int, tu testUser) {
		fmt.Printf("index: %d. User: %v", i, tu)
	}

	var actionKO Action[testUser] = func(i int, tu testUser) {
		if i == 0 {
			fmt.Print(tu.mails[3])
		}
		fmt.Printf("index: %d. User: %v", i, tu)
	}

	src := generateTestCaseList()

	type args struct {
		action   Action[testUser]
		errorFmt ErrorFormatter[testUser]
		src      any
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{"Iterate over list of testUser", args{action: actionOk, src: src}, nil},
		{"Iterate and generate and customizable error", args{action: actionKO, errorFmt: errorFmt, src: src}, fmt.Errorf("KO")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ForEach(tt.args.action, tt.args.src)
			if tt.want != nil && !reflect.ValueOf(got).IsZero() {
				err := got.WithErrorMessage(tt.args.errorFmt).Error()
				if err.Error() != tt.want.Error() {
					t.Errorf("Each() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestFilter2(t *testing.T) {

	//femaleResult := []testUser{{name: "Sarah", mails: []string{}, age: 43}}
	parent := map[string]testUser{"Kyle": {name: "Kyle" ,  secondName: "Risk", male: true , mails: []string{}, age: 43}}
	isFemale := func(tu testUser) bool {
		return !tu.male
	}

	type args[T any] struct {
		predicate Predicate[T]
		errorFmt  ErrorFormatter[testUser]
		source    any
		dest      any
	}
	type testsType [T any] struct {
		name      string
		args      args[T]
		want      any
		wantError bool
		err       error
	}	
	
	tests := []testsType
	{
		testsType[testUser]{
			name:      "Prueba",
			args:      args[testUser]{isFemale, errorFmt, generateTestCaseMap(), map[string]testUser{}},
			want:      nil,
			wantError: false,
			err:       nil,
		},
		//{"Filter female from list of test user", args{isFemale, errorFmt, generateTestCaseList(), []testUser{}}, femaleResult, false, nil},
		//testsType[testUser]{name: "Filter dad from map of test user", args: args{isFemale, errorFmt, generateTestCaseMap(), map[string]testUser{}}, want: parent, wantError: false, err:nil},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter2(tt.args.predicate, tt.args.source, tt.args.dest)

			if tt.want != nil && !reflect.ValueOf(got).IsZero() {
				err := got.WithErrorMessage(tt.args.errorFmt).Error()
				if err.Error() != tt.err.Error() {
					t.Errorf("Each() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
