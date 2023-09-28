package prompts

import (
	"fmt"
	"regexp"
	"strings"
)

type Route struct {
	schema      string
	description string
	regex       string
	prefix      string
	flags       []string
	args        []string
	lflags      []string
	handler     func(*Context)
}

type RouteErr struct {
	message string
}

func (e RouteErr) Error() string {
	return e.message
}

func (route *Route) dispatch(args []string) error {
	ctx := getContext(args, route.args)
	err := route.match(ctx)

	if err != nil {
		return err
	}

	route.handler(ctx)
	return nil
}

func (route *Route) match(ctx *Context) error {
	var err error

	if err = route.matchFlags(ctx); err != nil {
		return err
	}

	if err = route.matchLFlags(ctx); err != nil {
		return err
	}

	if err = route.matchSchema(ctx); err != nil {
		return err
	}

	return nil
}

func (route *Route) matchSchema(ctx *Context) error {
	if !regexp.MustCompile(route.regex).MatchString(formatFields(route.prefix + " " + ctx.path)) {
		return RouteErr{fmt.Sprintf("Usage :  %s", route.schema)}
	}

	return nil
}

func (route *Route) matchFlags(ctx *Context) error {
	if len(route.flags) == 0 && len(ctx.Flags) != 0 {
		return RouteErr{fmt.Sprintf("Usage :  %s", route.schema)}
	}

	exists := func(f string) bool {
		for _, f2 := range route.flags {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range ctx.Flags {
		if !exists(k) {
			return RouteErr{fmt.Sprintf("flag [%s] does not exists", k)}
		}
	}

	return nil
}

func (route *Route) matchLFlags(ctx *Context) error {
	if len(route.lflags) == 0 && len(ctx.LFlags) != 0 {
		return RouteErr{fmt.Sprintf("Usage :  %s", route.schema)}
	}

	exists := func(f string) bool {
		for _, f2 := range route.lflags {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range ctx.LFlags {
		if !exists(k) {
			return RouteErr{fmt.Sprintf("flag [%s] does not exists", k)}
		}
	}

	return nil
}

func (route *Route) parseSchema(schema string) error {
	schema = sanitize(schema)

	if !validSchema(schema) {
		return RouterErr{fmt.Sprintf("incorrect schema syntax : %s", schema)}
	}

	route.prefix = strings.SplitN(schema, " ", 2)[0]
	route.schema = schema
	route.splitUp(schema)

	return nil
}

func validSchema(schema string) bool {
	rx := regexp.MustCompile(`^[a-z0-9\:]+(\s+\<[a-z\_]+\>)*(\s+\[-{1,2}[A-z\-]+(\s*-{1,2}[A-z\-]+)*\])*(\s+\<[a-z\_]+\>)*(\s+\<[a-z\_]+\?\>)*$`)

	return rx.MatchString(schema)
}

func sanitize(s string) string {
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func (route *Route) splitUp(schema string) {
	if route.args == nil {
		route.args = make([]string, 0)
	}

	schema = regexp.MustCompile(`\<[a-z\_]+\??\>`).ReplaceAllStringFunc(schema, func(s string) string {
		return route.parseArg(s)
	})

	schema = regexp.MustCompile(`\[(\s*-{1,2}[A-z]+)+\]`).ReplaceAllStringFunc(schema, func(s string) string {
		return route.parseFlags(s)
	})

	schema = strings.ReplaceAll(schema, " ", "")

	route.regex = "^" + schema + "$"
}

func formatFields(f string) string {
	fs := strings.Fields(f)
	var fields string

	for _, v := range fs {
		fields += " " + regexp.MustCompile(`^-{1}[A-z]+$`).ReplaceAllStringFunc(v, func(s string) string {
			tmp := strings.TrimLeft(s, "-")
			s = ""
			for _, el := range tmp {
				s += " -" + string(el)
			}

			return strings.TrimSpace(s)
		})
	}

	return strings.TrimSpace(fields)
}

func (route *Route) Description(description string) *Route {
	route.description = description
	return route
}

func (route *Route) parseArg(arg string) string {
	arg = strings.TrimLeft(arg, "<")
	arg = strings.TrimRight(arg, ">")

	if strings.HasSuffix(arg, "?") {
		arg = strings.TrimRight(arg, "?")
		route.args = append(route.args, arg)
		return `(\s+` + delimiter + `[^` + delimiter + `]+` + delimiter + `)?`
	}

	route.args = append(route.args, arg)

	return `(\s+` + delimiter + `[^` + delimiter + `]+` + delimiter + `)`
}

func (route *Route) parseFlags(flags string) string {
	flags = strings.TrimLeft(flags, "[")
	flags = strings.TrimRight(flags, "]")
	fields := strings.Fields(flags)
	parts := make([]string, 0)

	for _, f := range fields {
		if strings.HasPrefix(f, "--") {
			route.lflags = append(route.lflags, f)
			parts = append(parts, f+`(=[^\s=]+)?`)
		} else {
			route.flags = append(route.flags, f)
			parts = append(parts, f)
		}
	}

	return `(\s+(` + strings.Join(parts, `|`) + `)+)*`
}
