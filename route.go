package goclitools

import (
	"fmt"
	"regexp"
	"strings"
)

type Route struct {
	schema  string
	regex   string
	prefix  string
	Flags   []string
	Args    int
	vars    map[string]string
	LFlags  []string
	Handler func(*Context)
}

type RouteErr struct {
	message string
}

func (e RouteErr) Error() string {
	return e.message
}

func (r *Route) match(c *Context) error {
	var err error

	err = r.matchFlags(c)

	if err != nil {
		return err
	}

	err = r.matchLFlags(c)

	if err != nil {
		return err
	}

	return nil
}

func (r *Route) matchFlags(c *Context) error {
	if len(r.Flags) == 0 && len(c.Flags) != 0 {
		return RouteErr{
			message: fmt.Sprintf("command [%s] does not expect flags,", r.prefix),
		}
	}

	exists := func(f string) bool {
		for _, f2 := range r.Flags {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range c.Flags {
		if !exists(k) {
			return RouteErr{
				message: fmt.Sprintf("flag [%s] does not exists", k),
			}
		}

	}

	return nil
}

func (r *Route) matchLFlags(c *Context) error {
	if len(r.LFlags) == 0 && len(c.LFlags) != 0 {
		return RouteErr{
			message: fmt.Sprintf("command [%s] does not expect flags,", r.prefix),
		}
	}

	exists := func(f string) bool {
		for _, f2 := range r.LFlags {
			if f == f2 {
				return true
			}
		}
		return false
	}

	for k := range c.LFlags {
		if !exists(k) {
			return RouteErr{
				message: fmt.Sprintf("flag [%s] does not exists", k),
			}
		}
	}

	return nil
}

func (r *Route) parseSchema(schema string) error {
	schema = sanitize(schema)

	if !validSchema(schema) {
		return fmt.Errorf("incorrect schema syntax : %s", schema)
	}

	r.prefix = strings.SplitN(schema, " ", 2)[0]
	r.schema = schema
	r.splitUp(schema)

	return nil
}

func validSchema(schema string) bool {
	rx := regexp.MustCompile(`^[a-z0-9]+(\s+\{[a-z\_]+\??\})*(\s+\[-{1,2}[A-z\-]+(\s-{1,2}[A-z\-]+)*\])*(\s+\{[a-z\_]+\??\})*$`)

	return rx.MatchString(schema)
}

func sanitize(s string) string {
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func (r *Route) splitUp(schema string) {
	if r.vars == nil {
		r.vars = make(map[string]string)
	}

	schema = regexp.MustCompile(`\{[a-z\_]+\??\}`).ReplaceAllStringFunc(schema, func(s string) string {
		s = strings.TrimLeft(s, "{")
		s = strings.TrimRight(s, "}")

		if strings.HasSuffix(s, "?") {
			s = strings.TrimRight(s, "?")
			r.vars[s] = ""
			return `[^\s]*`
		}

		r.vars[s] = ""

		return `[^\s]+`
	})

	schema = regexp.MustCompile(`\[(\s*-{1,2}[A-z]+)+\]`).ReplaceAllStringFunc(schema, func(s string) string {
		s = strings.TrimLeft(s, "[")
		s = strings.TrimRight(s, "]")
		fields := strings.Fields(s)
		parts := make([]string, 0)

		for _, f := range fields {
			if strings.HasPrefix(f, "--") {
				r.LFlags = append(r.LFlags, f)
				parts = append(parts, f+`(=[^\s=]+)?`)
			} else {
				r.Flags = append(r.Flags, f)
				parts = append(parts, f)
			}
		}

		return `(\s*(` + strings.Join(parts, `|`) + `)+\s*)*`
	})

	r.regex = "^" + schema + "$"
}
