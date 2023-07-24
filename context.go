package goclitools

import (
	"regexp"
	"strings"
)

type Context struct {
	Flags  map[string]int
	LFlags map[string]string
	Args   map[string]string
}

func getContext(arguments []string, args_names []string) *Context {
	var ctx = Context{
		Flags:  make(map[string]int),
		Args:   make(map[string]string),
		LFlags: make(map[string]string),
	}

	args := make([]string, 0, len(arguments))

	for _, i := range arguments {
		switch true {
		case regexp.MustCompile(`^-[A-z0-9]+$`).MatchString(i):
			ctx.loadFlags(i)
		case regexp.MustCompile(`^--[a-z0-9\-]+(=[A-z0-9\-_]+)?$`).MatchString(i):
			ctx.loadLFlags(i)
		default:
			args = append(args, i)
		}
	}

	for k, v := range args {
		if k < len(args_names) {
			ctx.Args[args_names[k]] = v
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

// determine wether l long flag is present or not
//
// for example HasLFlag("--help")
func (r *Context) HasLFlag(l string) bool {
	_, ok := r.LFlags[l]

	return ok
}

// get a long flag value
//
// for example , in case of "--lflag=something" , GetLFlag("--lflag") returns "something"
func (r *Context) GetLFlag(l string) string {
	v, ok := r.LFlags[l]

	if ok {
		return v
	}

	return ""
}

// scan the long flag to dst
func (r *Context) ScanLflag(l string, dst *string) {
	*dst = r.LFlags[l]
}

// Get an argument by its name , or empty string if the argument doesn't exists
func (ctx *Context) GetArg(name string) string {
	return ctx.Args[name]
}

func (ctx *Context) loadFlags(s string) {
	s = strings.TrimPrefix(s, "-")
	flags := strings.Split(s, "")

	for _, f := range flags {
		ctx.Flags["-"+f] = ctx.Flags["-"+f] + 1
	}
}

func (ctx *Context) loadLFlags(s string) {
	pair := strings.SplitN(s, "=", 2)

	if len(pair) == 2 {
		ctx.LFlags[pair[0]] = pair[1]
	} else {
		ctx.LFlags[pair[0]] = ""
	}
}
