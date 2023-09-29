# Prompts

Build command line prompts with ease, prompts provides several apis to help you create intuitive TUI faster,

## API

### Input

The input api allows you to prompt the user for an input field , it returns the input value

- usage:

```go
// [...]

value, err := prompts.InputBox(prompts.InputOptions{
	Secure:      false, // hides the user input, very common for passwords
	Label:       "what is your name?",
	Placeholder: "what is your name",
	Required:    true,
	Validator:   func(value string) error {// will be called when user submit, and returned error will be displayed to the user below the input
		if len(value) < 3{
			return fmt.Errorf("minimum len is 3")
		}
		return nil
	},
})

if err != nil{
	log.Fatal(err)
}

fmt.Println("selected " + value)
```

- result

![Screencast from 09-28-2023 08_00_26 PM](https://github.com/yassinebenaid/prompts/assets/101285507/5e4e8c68-5e6a-4cb1-8ca0-169203ca5f6c)

### Password Input

The password input api is just normal input but with `Secret` option set to `true` ,

- usage:

```go
// [...]

value, err := prompts.InputBox(prompts.InputOptions{
	Secure:      true, // set password mode
	Label:       "what is your password?",
	Placeholder: "what is your password",
	Required:    true,
	Validator:   func(value string) error {// will be called when user submit, and returned error will be displayed to the user below the input
		if len(value) < 3{
			return fmt.Errorf("minimum len is 3")
		}
		return nil
	},
})

if err != nil{
	log.Fatal(err)
}

fmt.Println("password : " + value)
```

- result
  
  ![Screencast from 09-29-2023 09_59_36 AM (1)](https://github.com/yassinebenaid/prompts/assets/101285507/c3c54db1-5964-41b6-90c6-9f5614f28448)


### Confirmation Input

The confirmation api can be used to prompt the user for confirmation , it returns a boolean ,

- usage:

```go
// [...]

const DEFAULT = true
value, err := prompts.ConfirmBox("are you sure ?", DEFAULT)

if err != nil {
	log.Fatal(err)
}

fmt.Println("answer : ", value)
```

- result
  
  ![Screencast from 09-29-2023 05_51_27 PM](https://github.com/yassinebenaid/prompts/assets/101285507/83536f0d-a551-4963-ae5e-66d9eff54d89)

### Radio Input

The radio api can be used to prompt the user to choose one of several options , it returns a the index number ,

- usage:

```go
// [...]

genders := []string{"male", "female"}
value, err := prompts.RadioBox("Choose your gender : ", genders)

if err != nil {
	log.Fatal(err)
}

fmt.Println("gender : ", genders[value])
```

- result

   ![Screencast from 09-29-2023 06_01_03 PM](https://github.com/yassinebenaid/prompts/assets/101285507/bb264973-7112-4faa-a19d-c1bbc9fa1a1e)

