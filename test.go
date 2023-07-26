package wind

import "fmt"

func Test() {
	router := NewRouter()
	log := Log{}

	// router.Add("do", func(ctx *Context) {
	// 	Success("Can run simple commands")
	// })

	router.Add("do [-a]", func(ctx *Context) {
		if ctx.HasFlag("-a") {
			Success("Can run simple commands with options")
		} else {
			Warning("Can run simple commands with options but flag is not detected")
		}

	})

	router.Fallback(func(ctx *Context) {
		fmt.Println("fallback")
	})

	cmds := []string{
		"do ",
		// "do2 -a",
	}

	for _, i := range cmds {
		log.Info("running command : [" + i + "]")
		if err := router.Test(i); err != nil {
			log.Fatal(err)
		}
	}

}
