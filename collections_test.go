package collections


import (
        "fmt"
        "reflect"
        "strings"
        "testing"
)

var (
        john  = testUser{name: "John", secondName: "Connor", mails: []string{}, age: 10, male: true}
        sarah = testUser{name: "Sarah", secondName: "Connor", mails: []string{}, age: 43}
        kyle  = testUser{name: "Kyle", secondName: "Risk", male: true, mails: []string{}, age: 43}
)

type testUser struct {
        name       string
        secondName string
        mails      []string
        age        int
        male       bool
}

func generateTestCaseList() []testUser {
        return []testUser{john, sarah, kyle}
}

func generateTestCaseMap() map[string]testUser {
        result := map[string]testUser{}
        for _, item := range generateTestCaseList() {
                result[item.name] = item
        }
        return result
}

func errorFmtAsString(user string) error {
        return fmt.Errorf("KO")
}

func errorFmtAsTouple(user Touple) error {
        return fmt.Errorf("KO")
}

func errorFmt(user testUser) error {
        return fmt.Errorf("KO")
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
                action Action[testUser]
                src    any
        }
        tests := []struct {
                name      string
                args      args
                want      error
                wantError bool
        }{
                {"Iterate over list of testUser", args{action: actionOk, src: src}, nil, false},
                {"Iterate and generate and customizable error", args{action: actionKO, src: src}, fmt.Errorf("KO"), true},
        }
        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        err := ForEach(tt.args.action, tt.args.src)

                        if tt.wantError && err == nil {
                                t.Errorf("Each() = %v, want %v", err, tt.want)
                        }

                })
        }
}

type argsFilter[T any] struct {
        predicate Predicate[T]
        source    any
        dest      any
}
type testsTypeFilter[T any] struct {
        name      string
        args      argsFilter[T]
        want      any
        wantError bool
        err       error
}

func (tt testsTypeFilter[T]) runTestFilter(testRunner *testing.T) {
        testRunner.Run(tt.name, func(t *testing.T) {
                got := Filter(tt.args.predicate, tt.args.source, tt.args.dest)
                if !tt.wantError && got != nil {
                        t.Errorf("filter() = %v, wantError %v", got, tt.wantError)
                }
                if !tt.wantError && !reflect.DeepEqual(tt.want, tt.args.dest) {
                        t.Errorf("Filter() = %v, want %v", tt.args.dest, tt.want)
                }
        },
        )
}

func isMale(tu testUser) bool {
        return tu.male
}

func isMaleAsTouple(tu Touple) bool {
        testUser := tu.Value.(testUser)
        flag := testUser.male
        return flag
}

func TestFilter(t *testing.T) {

        resultListFiltered := []testUser{
                john,
                kyle,
        }
        resultMapFiltered := map[string]testUser{"John": john, "Kyle": kyle}

        lstUsers := []testUser{}
        mapUsers := map[string]testUser{}

        testsTypeFilter[testUser]{
                name:      "Filter male from list of test user",
                args:      argsFilter[testUser]{isMale, generateTestCaseList(), &lstUsers},
                want:      &resultListFiltered,
                wantError: false,
                err:       nil}.runTestFilter(t)

        testsTypeFilter[Touple]{
                name:      "Filter male from map of test user",
                args:      argsFilter[Touple]{isMaleAsTouple, generateTestCaseMap(), mapUsers},
                want:      resultMapFiltered,
                wantError: false,
                err:       nil,
        }.runTestFilter(t)
}

type argsMap[T any] struct {
        mapper Mapper[T]
        source any
        dest   any
}
type testsTypeMap[T any] struct {
        name      string
        args      argsMap[T]
        want      any
        wantError bool
        err       error
}

func compareTestUser(userA, userB string) int {
        return strings.Compare(userA, userB)
}
func (tt testsTypeMap[T]) runTestMap(testRunner *testing.T) {
        testRunner.Run(tt.name, func(t *testing.T) {
                got := Map(tt.args.mapper, tt.args.source, tt.args.dest)

                if !IsMap(tt.args.dest) {
                        SortBy(compareTestUser, tt.args.dest)
                }

                if !tt.wantError && got != nil {
                        t.Errorf("Map() KO = %v, wantError %v", got, tt.wantError)
                }
                if !tt.wantError && !reflect.DeepEqual(tt.want, tt.args.dest) {
                        t.Errorf("Map() = %v, want %v", tt.args.dest, tt.want)
                }
        },
        )
}

var mapperToNamesFromMap Mapper[Touple] = func(s Touple) any {
        user := s.Value.(testUser)
        return fmt.Sprintf("%s %s", user.name, user.secondName)
}

var mapperToNamesFromList Mapper[testUser] = func(s testUser) any {
        return fmt.Sprintf("%s %s", s.name, s.secondName)
}

var mapperSplitName Mapper[string] = func(s string) any {
        nameSplited := strings.Split(s, " ")
        return Touple{nameSplited[0], nameSplited[1]}
}

// Test functions
func TestMap(t *testing.T) {
        result := []string{"John Connor", "Kyle Risk", "Sarah Connor"}
        lstUsersFromList := []string{}
        lstUsersFromMap := []string{}
        names := map[string]string{}

        resulNames := map[string]string{"Sarah": "Connor", "Kyle": "Risk", "John": "Connor"}

        testsTypeMap[testUser]{
                name:      "Geneare a list of names from Map.",
                args:      argsMap[testUser]{mapperToNamesFromList, generateTestCaseList(), &lstUsersFromList},
                want:      &result,
                wantError: false,
                err:       nil}.runTestMap(t)

        testsTypeMap[Touple]{
                name:      "Generate a new Map taken the name as key and the struct as value.",
                args:      argsMap[Touple]{mapperToNamesFromMap, generateTestCaseMap(), &lstUsersFromMap},
                want:      &result,
                wantError: false,
                err:       nil,
        }.runTestMap(t)

        testsTypeMap[string]{
                name:      "Generate a new Map, taken as keyy the first name, and value the second name.",
                args:      argsMap[string]{mapperSplitName, lstUsersFromMap, names},
                want:      resulNames,
                wantError: false,
                err:       nil,
        }.runTestMap(t)
}

type argsGroupBy[T any] struct {
        groupBy KeySelector[T]
        source  any
        dest    any
}
type testsTypeGroupBy[T any] struct {
        name      string
        args      argsGroupBy[T]
        want      any
        wantError bool
        err       error
}

func (tt testsTypeGroupBy[T]) runTestGroupBy(testRunner *testing.T) {
        testRunner.Run(tt.name, func(t *testing.T) {
                got := GroupBy(tt.args.groupBy, tt.args.source, tt.args.dest)
                if !tt.wantError && got != nil {
                        t.Errorf("Map() KO = %v, wantError %v", got, tt.wantError)
                }
                if !tt.wantError && !reflect.DeepEqual(tt.want, tt.args.dest) {
                        t.Errorf("Map() = %v, want %v", tt.args.dest, tt.want)
                }
        },
        )
}

var keySelectorBySex KeySelector[testUser] = func(tu testUser) any {
        result := "female"
        if tu.male {
                result = "male"
        }
        return result
}

func TestGroupBy(t *testing.T) {
        result := map[string][]testUser{"male": {john, kyle}, "female": {sarah}}
        testsTypeGroupBy[testUser]{
                name:      "Filter male from list of test user",
                args:      argsGroupBy[testUser]{keySelectorBySex, generateTestCaseList(), map[string][]testUser{}},
                want:      result,
                wantError: false,
                err:       nil}.runTestGroupBy(t)
}

func TestSortBy(t *testing.T) {
        src := []int{8, 2, 809, 40, 43, 32838, 2, 67}
        want := []int{2, 2, 8, 40, 43, 67, 809, 32838}
        compareInt := func(a, b int) int {
                return a - b
        }
        type args[T any] struct {
                comparator Comparator[int]
                source     any
        }
        type testsType[T any] struct {
                name      string
                args      args[T]
                want      any
                wantError bool
        }

        tests := []testsType[int]{
                {"Sort int slice", args[int]{comparator: compareInt, source: &src}, &want, false},
                {"Should generate and error Sort int slice", args[int]{comparator: compareInt, source: src}, nil, true},
        }
        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        got := SortBy(tt.args.comparator, tt.args.source)

                        if tt.wantError && got == nil {
                                t.Errorf("sortBy() KO = %v, wantError %v", got, tt.wantError)
                        } else if !tt.wantError && !reflect.DeepEqual(tt.want, tt.args.source) {
                                t.Errorf("List() = %v, want %v", tt.args.source, tt.want)
                        }
                })
        }
}
