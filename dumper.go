package goclitools

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func DD(v any) {
	t := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		dstr(t.String())
	case reflect.Int:
		dnum[int](int(t.Int()))
	case reflect.Float64, reflect.Float32:
		dnum[float64](t.Float())
	case reflect.Complex128, reflect.Complex64:
		dnum[complex128](t.Complex())
	}
}

func dstr(s string) {
	dmp := Sprint(SprintRGB(`"`+s+`"`, 100, 200, 1), Bold)
	fmt.Println(dmp + string(Tab) + getfl())
}

type Num interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint | float32 | float64 | complex64 | complex128
}

func dnum[N Num](n N) {
	v := fmt.Sprintf("%v", n)
	fmt.Println(Sprint(SprintRGB(v, 200, 110, 0), Bold) + string(Tab) + getfl())
}

func getfl() string {
	_, f, l, ok := runtime.Caller(3)

	if !ok {
		return ""
	}

	fl := Sprint(f+":"+fmt.Sprint(l), T_BrightBlack)

	wd, err := os.Getwd()

	if err != nil {
		return fl
	}

	return strings.Replace(fl, wd, "", 1)
}
