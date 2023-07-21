package goclitools

import (
	"regexp"
	"strings"
)

type Context struct {
	Flags  map[string]int
	LFlags map[string]string
	Args   []string
}

func getContext(Args []string) *Context {
	var ctx = Context{
		Flags:  make(map[string]int),
		Args:   make([]string, 0, len(Args)),
		LFlags: make(map[string]string),
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
				ctx.LFlags[kv[0]] = kv[1]
			} else {
				ctx.LFlags[kv[0]] = ""
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

// Get an argument by its index , or "" if doesn't exists
//
// this functions execluds the flags
func (r *Context) GetArg(index int) string {

	if len(r.Args) < index || index < 1 {
		return ""
	}

	return r.Args[index-1]
}
