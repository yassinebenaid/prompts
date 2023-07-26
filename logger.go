package wind

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = 0 + iota
	WarnLevel
	InfoLevel
)

type Log struct {
	Writer     io.Writer
	WithCaller bool
	Format     string
	Level      LogLevel
	slc        bool
}

func (log *Log) silence(l LogLevel) {
	switch true {
	case log.Level == DebugLevel:
		log.slc = false
	case log.Level == WarnLevel && l >= WarnLevel:
		log.slc = false
	case log.Level == InfoLevel && l >= InfoLevel:
		log.slc = false
	default:
		log.slc = true
	}
}

// log info message to the console and optionaly in the log file
//
// Info log messages regardless of log level
func (log *Log) Info(m any, args ...any) {
	log.silence(InfoLevel)
	ms := fmt.Sprintf("%s", m)
	log.log(Sprint("INFO ", T_BrightCyan, Bold), ms, args...)
}

// log debug message to the console and optionaly in the log file
//
// debug log messages only of log level is DebugLevel
func (log *Log) Debug(m any, args ...any) {
	log.silence(DebugLevel)
	ms := fmt.Sprintf("%s", m)
	log.log(Sprint("DEBUG", T_BrightMagenta, Bold), ms, args...)

}

// log warning message to the console and optionaly in the log file
//
// warn log messages only of log level is WarnLevel
func (log *Log) Warn(m any, args ...any) {
	log.silence(WarnLevel)
	ms := fmt.Sprintf("%s", m)
	log.log(Sprint("WARN ", T_BrightYellow, Bold), ms, args...)
}

// log error message to the console and optionaly in the log file
//
// Error log messages only of log level is DebugLevel
func (log *Log) Error(m any, args ...any) {
	log.silence(DebugLevel)
	ms := fmt.Sprintf("%s", m)
	log.log(Sprint("ERROR", T_BrightRed, Bold), ms, args...)
}

// log fatal errors to the console and optionaly in the log file
//
// fatal log messages only of log level is DebugLevel
func (log *Log) Fatal(m any, args ...any) {
	log.silence(DebugLevel)
	ms := fmt.Sprintf("%s", m)
	label := SprintRGB("FATAL", 255, 0, 150)
	log.log(Sprint(label, Bold), ms, args...)
	os.Exit(1)
}

func (log *Log) log(label string, m string, args ...any) {
	var ts string
	if log.Format != "" {
		ts = time.Now().Format(log.Format)
	} else {
		ts = time.Now().Format(time.DateTime)
	}
	caller := ""

	if log.WithCaller {
		caller = log.getfl()
	}

	if !log.slc {
		fmt.Println(ts, label, caller, m, log.kvpair(args))
	}

	if log.Writer != nil {
		log.Writer.Write(log.sanitize(ts, label, caller, m, log.kvpair(args)))
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
		case isNumeric(kv[index+1]):
			out += SprintRGB(fmt.Sprint(kv[index+1])+" ", 20, 160, 200)
		case isBool(kv[index+1]):
			out += SprintRGB(fmt.Sprint(kv[index+1])+" ", 225, 110, 30)
		default:
			out += SprintRGB(`"`+fmt.Sprint(kv[index+1])+`" `, 100, 230, 80)
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

func (log *Log) getfl() string {
	_, f, l, ok := runtime.Caller(3)

	if !ok {
		return ""
	}

	wd, err := os.Getwd()

	if err != nil {
		return ""
	}

	f = strings.Replace(f, wd+"/", "", 1)
	f = strings.Replace(f, wd, "", 1)

	return Sprint("<"+f+":"+fmt.Sprint(l)+">", T_BrightBlack)
}
