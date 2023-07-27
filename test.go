package wind

import "fmt"

func Test() {
	router := NewRouter()
	log := Log{}

	router.Add("test1", func(ctx *Context) {
		Success("Can run simple commands")
	})

	router.Add("test2 [-a]", func(ctx *Context) {
		if ctx.HasFlag("-a") {
			Success("Can run simple commands with flags")
		} else {
			Warning("Can run simple commands with flags but flag is not detected")
		}
	})

	router.Add("test3 [-a -s]", func(ctx *Context) {
		if ctx.HasFlag("-a") && ctx.HasFlag("-s") {
			Success("Can run simple commands with many flags")
		} else {
			Warning("Can run simple commands with many flags but flags are not detected")
		}
	})

	router.Add("test4 [-a -s] <name>", func(ctx *Context) {
		name := ctx.GetArg("name")
		if name != "" {
			Success("Can run commands with arguments : name = " + name)
		} else {
			Warning("Cannot run commands with arguments")
		}
	})

	router.Add("test5 [-a -s] <name?>", func(ctx *Context) {
		name := ctx.GetArg("name")

		Success("Can run commands with optional arguments : name = " + name)

	})

	router.Fallback(func(ctx *Context) {
		fmt.Println("fallback")
	})

	cmds := []string{
		"test1",
		"test2 -a",
		"test3 -a -s",
		"test3 -as",
		"test3 -aaa -sss",
		"test4 somename",
		"test5 something",
		"test5 ",
	}

	for _, i := range cmds {
		log.Info("running command : " + i)
		_, err := router.Test(i)

		if err != nil {
			log.Fatal(err)
		}
	}

}
