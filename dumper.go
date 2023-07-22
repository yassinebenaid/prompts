package goclitools

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

func DD(value any) {
	fmt.Print(getfl(2), "\n")
	dump(value)
	fmt.Println()
}

func dumpString(s string) {
	Print(SprintRGB(`"`+s+`"`, 50, 170, 10))
}

type Num interface {
	int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint | float32 | float64 | complex64 | complex128
}

func dumpNumber[N Num](n N) {
	v := fmt.Sprintf("%v", n)
	fmt.Print(SprintRGB(v, 10, 90, 150))
}

func dumpNumSlice[N Num](s []N) {
	fmt.Print(typeString(s, fmt.Sprintf(":%d", len(s))))
	PrintRGB(" [\n", 150, 90, 10)
	for _, v := range s {
		fmt.Print("   ")
		dumpNumber[N](v)
		fmt.Println()
	}
	PrintRGB("]\n", 150, 90, 10)
}

func dumpStrSlice(s []string) {
	fmt.Print(typeString(s, fmt.Sprintf(":%d", len(s))))
	PrintRGB(" [\n", 150, 90, 10)
	for _, v := range s {
		fmt.Print("   ")
		dumpString(v)
		fmt.Println()
	}
	PrintRGB("]\n", 150, 90, 10)
}

// func dumpStrSlice2(s []string) {
// 	Print(fmt.Sprintf("[]int:%d {\n", len(s)), T_Blue)
// 	for _, v := range s {
// 		fmt.Print("   ")
// 		dumpString(v)
// 		Print(": ", T_BrightBlack)
// 		dumpString(v)
// 		fmt.Println()
// 	}
// 	Print("}\n", T_Blue)
// }

func getfl(skip int) string {
	_, f, l, ok := runtime.Caller(skip)

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

func dump(value any) {
	switch value := value.(type) {
	case string:
		dumpString(value)
	case int:
		dumpNumber[int](value)
	case int8:
		dumpNumber[int8](value)
	case int16:
		dumpNumber[int16](value)
	case uint:
		dumpNumber[uint](value)
	case uint8:
		dumpNumber[uint8](value)
	case uint16:
		dumpNumber[uint16](value)
	case uint32:
		dumpNumber[uint32](value)
	case uint64:
		dumpNumber[uint64](value)
	case float32:
		dumpNumber[float32](value)
	case float64:
		dumpNumber[float64](value)
	case complex64:
		dumpNumber[complex64](value)
	case complex128:
		dumpNumber[complex128](value)
	case []int:
		dumpNumSlice[int](value)
	case []int8:
		dumpNumSlice[int8](value)
	case []int16:
		dumpNumSlice[int16](value)
	case []uint:
		dumpNumSlice[uint](value)
	case []uint8:
		dumpNumSlice[uint8](value)
	case []uint16:
		dumpNumSlice[uint16](value)
	case []uint32:
		dumpNumSlice[uint32](value)
	case []uint64:
		dumpNumSlice[uint64](value)
	case []float32:
		dumpNumSlice[float32](value)
	case []float64:
		dumpNumSlice[float64](value)
	case []complex64:
		dumpNumSlice[complex64](value)
	case []complex128:
		dumpNumSlice[complex128](value)
	case []string:
		dumpStrSlice(value)
	}
}

func typeString(value any, with string) string {
	return SprintRGB(reflect.TypeOf(value).String()+with, 10, 90, 150)
}
