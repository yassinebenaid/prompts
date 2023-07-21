package goclitools

import (
	"fmt"
)

type Route struct {
	prefix  string
	Flags   []string
	Args    int
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
