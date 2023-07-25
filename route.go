package wind

import (
	"fmt"
	"regexp"
	"strings"
)

type Route struct {
	schema  string
	path    string
	regex   string
	prefix  string
	flags   []string
	vars    []string
	lflags  []string
	handler func(*Context)
}

type RouteErr struct {
	message string
}

func (e RouteErr) Error() string {
	return e.message
}

func (route *Route) match(ctx *Context) error {
	var err error

	if err = route.matchFlags(ctx); err != nil {
		return err
	}

	if err = route.matchLFlags(ctx); err != nil {
		return err
	}

	if err = route.matchSchema(); err != nil {
		return err
	}

	return nil
}

func (route *Route) matchSchema() error {
	if !regexp.MustCompile(route.regex).MatchString(route.path) {
		return RouteErr{fmt.Sprintf("Usage :  %s", route.schema)}
	}

	return nil
}

func (route *Route) matchFlags(c *Context) error {
	if len(route.flags) == 0 && len(c.Flags) != 0 {
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

	for k := range c.Flags {
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
	rx := regexp.MustCompile(`^[a-z0-9\:]+(\s+\<[a-z\_]+\>)*(\s+\[-{1,2}[A-z\-]+(\s*-{1,2}[A-z\-]+)+\])*(\s+\<[a-z\_]+\>)*(\s+\<[a-z\_]+\?\>)*$`)

	return rx.MatchString(schema)
}

func sanitize(s string) string {
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func (route *Route) splitUp(schema string) {
	if route.vars == nil {
		route.vars = make([]string, 0)
	}

	schema = regexp.MustCompile(`\<[a-z\_]+\??\>`).ReplaceAllStringFunc(schema, func(s string) string {
		s = strings.TrimLeft(s, "<")
		s = strings.TrimRight(s, ">")

		if strings.HasSuffix(s, "?") {
			s = strings.TrimRight(s, "?")
			route.vars = append(route.vars, s)
			return `(\s+[^\s\-]+)?`
		}

		route.vars = append(route.vars, s)

		return `(\s+[^\s\-]+)`
	})

	schema = regexp.MustCompile(`\[(\s*-{1,2}[A-z]+)+\]`).ReplaceAllStringFunc(schema, func(s string) string {
		s = strings.TrimLeft(s, "[")
		s = strings.TrimRight(s, "]")
		fields := strings.Fields(s)
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
	})
	schema = strings.ReplaceAll(schema, " ", "")

	route.regex = "^" + schema + "$"
}
