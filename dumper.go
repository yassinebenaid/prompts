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

func dumpNumStrMap[N Num](s map[N]string) {
	fmt.Print(typeString(s, fmt.Sprintf(":%d", len(s))))
	PrintRGB(" [\n", 150, 90, 10)
	for k, v := range s {
		fmt.Print("   ")
		dumpNumber[N](k)
		PrintRGB(": ", 150, 90, 10)
		dumpString(v)
		fmt.Println()
	}
	PrintRGB("]\n", 150, 90, 10)
}

func dumpStrNumMap[N Num](s map[string]N) {
	fmt.Print(typeString(s, fmt.Sprintf(":%d", len(s))))
	PrintRGB(" [\n", 150, 90, 10)
	for k, v := range s {
		fmt.Print("   ")
		dumpString(k)
		PrintRGB(": ", 150, 90, 10)
		dumpNumber[N](v)
		fmt.Println()
	}
	PrintRGB("]\n", 150, 90, 10)
}

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
	case map[int]string:
		dumpNumStrMap[int](value)
	case map[int8]string:
		dumpNumStrMap[int8](value)
	case map[int16]string:
		dumpNumStrMap[int16](value)
	case map[uint]string:
		dumpNumStrMap[uint](value)
	case map[uint8]string:
		dumpNumStrMap[uint8](value)
	case map[uint16]string:
		dumpNumStrMap[uint16](value)
	case map[uint32]string:
		dumpNumStrMap[uint32](value)
	case map[uint64]string:
		dumpNumStrMap[uint64](value)
	case map[float32]string:
		dumpNumStrMap[float32](value)
	case map[float64]string:
		dumpNumStrMap[float64](value)
	case map[complex64]string:
		dumpNumStrMap[complex64](value)
	case map[complex128]string:
		dumpNumStrMap[complex128](value)
	case map[string]int:
		dumpStrNumMap[int](value)
	case map[string]int8:
		dumpStrNumMap[int8](value)
	case map[string]int16:
		dumpStrNumMap[int16](value)
	case map[string]uint:
		dumpStrNumMap[uint](value)
	case map[string]uint8:
		dumpStrNumMap[uint8](value)
	case map[string]uint16:
		dumpStrNumMap[uint16](value)
	case map[string]uint32:
		dumpStrNumMap[uint32](value)
	case map[string]uint64:
		dumpStrNumMap[uint64](value)
	case map[string]float32:
		dumpStrNumMap[float32](value)
	case map[string]float64:
		dumpStrNumMap[float64](value)
	case map[string]complex64:
		dumpStrNumMap[complex64](value)
	case map[string]complex128:
		dumpStrNumMap[complex128](value)
	}
}

func typeString(value any, with string) string {
	return SprintRGB(reflect.TypeOf(value).String()+with, 10, 90, 150)
}
