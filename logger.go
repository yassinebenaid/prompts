package wind

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

type Logger interface {
	Info(m string, args ...any)
	Debug(m string, args ...any)
	Fatal(m string, args ...any)
	Warn(m string, args ...any)
	Error(m string, args ...any)
}

type Log struct {
	Writer io.Writer
}

func (log *Log) Info(m string, args ...any) {
	log.log(Sprint(" INFO  ", T_BrightCyan, Bold), m, args...)
}

func (log *Log) Debug(m string, args ...any) {
	log.log(Sprint(" DEBUG ", T_BrightMagenta, Bold), m, args...)

}

func (log *Log) Warn(m string, args ...any) {
	log.log(Sprint(" WARN  ", T_BrightYellow, Bold), m, args...)

}

func (log *Log) Error(m string, args ...any) {
	log.log(Sprint(" ERROR ", T_BrightRed, Bold), m, args...)

}

func (log *Log) Fatal(m string, args ...any) {
	label := SprintRGB(" FATAL ", 255, 0, 150)
	log.log(Sprint(label, Bold), m, args...)
	os.Exit(1)
}

func (log *Log) log(label string, m string, args ...any) {
	time := time.Now().Format(time.DateTime)
	fmt.Println(time, label, m, log.kvpair(args))
	if log.Writer != nil {
		log.Writer.Write(log.sanitize(time, label, m, log.kvpair(args)))
	}
}

func (log *Log) kvpair(kv []any) string {
	if len(kv) >= 2 && len(kv)%2 != 0 {
		return ""
	}

	var out string
	var index int

	for index < len(kv) {
		out += Sprint(fmt.Sprint(kv[index])+"=", T_BrightBlack)

		switch true {
		case isNumeric(kv[index+1]) || isBool(kv[index+1]):
			out += Sprint(fmt.Sprint(kv[index+1])+" ", Dim)
		default:
			out += Sprint(`"`+fmt.Sprint(kv[index+1])+`" `, Dim)
		}
		index += 2
	}

	return out
}

func isNumeric(v any) bool {
	t := fmt.Sprintf("%T", v)
	return strings.HasPrefix(t, "int") || strings.HasPrefix(t, "uint") || strings.HasPrefix(t, "float") || strings.HasPrefix(t, "complex")
}

func isBool(v any) bool {
	return fmt.Sprintf("%T", v) == "bool"
}

func (log *Log) sanitize(s ...string) []byte {
	str := strings.Join(s, " ")
	str = regexp.MustCompile(`((\x1b)|(\033))\[[^m]+m`).ReplaceAllString(str, "")
	return []byte(str + "\n")
}
