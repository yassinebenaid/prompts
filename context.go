package goclitools

import (
	"regexp"
	"strings"
)

type Context struct {
	Flags   map[string]int
	Args    []string
	Options map[string]string
}

func getContext(Args []string) *Context {
	var ctx = Context{
		Flags:   make(map[string]int),
		Args:    make([]string, 0, len(Args)),
		Options: make(map[string]string),
	}

	flag := regexp.MustCompile(`^-[A-z0-9]+$`)
	opt := regexp.MustCompile(`^--[a-z0-9\-]+(=[A-z0-9\-_]+)?$`)

	for _, i := range Args {
		switch true {
		case flag.MatchString(i):
			fs := strings.Split(strings.TrimPrefix(i, "-"), "")

			for _, f := range fs {
				ctx.Flags["-"+f] = ctx.Flags["-"+f] + 1
			}
		case opt.MatchString(i):
			kv := strings.SplitN(i, "=", 2)
			if len(kv) == 2 {
				ctx.Options[kv[0]] = kv[1]
			} else {
				ctx.Options[kv[0]] = ""
			}
		default:
			ctx.Args = append(ctx.Args, i)
		}
	}

	return &ctx
}

// determine wether f flag is present or not
//
// for example HasFlag("-h")
func (r *Context) HasFlag(f string) bool {
	_, ok := r.Flags[f]
	return ok
}

// determine wether f flag is present or not
//
// for example HasFlag("-h")
func (r *Context) GetFlagCount(f string) int {
	return r.Flags[f]
}

// determine wether opt option is present or not
//
// for example HasOption("--help")
func (r *Context) HasOption(opt string) bool {
	_, ok := r.Options[opt]

	return ok
}

// get an option value
//
// for example , in case of "--path=some/path/here" , GetOption("--path") returns "/some/path/here"
func (r *Context) GetOption(opt string) string {
	v, ok := r.Options[opt]

	if ok {
		return v
	}

	return ""
}

// scan the option and seve it to dst
func (r *Context) ScanOption(opt string, dst *string) {
	*dst = r.Options[opt]
}

// Get an argument by its index , or "" if doesn't exists
//
// this functions returns all args that are not a flag , and not options
func (r *Context) GetArg(index int) string {

	if len(r.Args) < index || index < 1 {
		return ""
	}

	return r.Args[index-1]
}
