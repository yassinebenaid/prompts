# Prompts

Build command line prompts with ease, prompts provides several apis to help you create intuitive TUI faster,

 ![Screenshot from 2023-09-30 19-38-45](https://github.com/yassinebenaid/prompts/assets/101285507/7ef5edb3-b13c-4e64-b03d-74f627165982)


## API

### Input

The input api allows you to prompt the user for an input field , it returns the input value

- **Usage**:

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

- **Result**:

![image](https://github.com/yassinebenaid/prompts/assets/101285507/de482c92-7ab3-4a36-a68a-422f3f74de02)
![image](https://github.com/yassinebenaid/prompts/assets/101285507/025942de-025d-4e5d-a476-172c3311cf5e)

#

### Password Input

The password input api is just normal input but with `Secure` option set to `true` ,

- **Usage**:

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

- **Result**:
  
 ![image](https://github.com/yassinebenaid/prompts/assets/101285507/825f5688-9d55-4e2d-acf2-c5b308dfc4a5)


#

### Confirmation Input

The confirmation api can be used to prompt the user for confirmation , it returns a boolean ,

- **Usage**:

```go
// [...]

const DEFAULT = true
value, err := prompts.ConfirmBox("are you sure ?", DEFAULT)

if err != nil {
	log.Fatal(err)
}

fmt.Println("answer : ", value)
```

- **Result**:

  ![image](https://github.com/yassinebenaid/prompts/assets/101285507/8a6c51cb-847d-415a-b11e-96cbf0a2beeb)


#

### Radio Input

The radio api can be used to prompt the user to choose one of several options , it returns a the index of the checked option ,

- **Usage**:

```go
// [...]

genders := []string{"male", "female"}
value, err := prompts.RadioBox("Choose your gender : ", genders)

if err != nil {
	log.Fatal(err)
}

fmt.Println("gender : ", genders[value])
```

- **Result**:

  ![image](https://github.com/yassinebenaid/prompts/assets/101285507/52e82922-844c-4c3a-89b7-a9d8800cd5b0)


#

### Select Box

The select box api can be used to prompt the user to choose between several options , it returns a slice of selected  indexes,

- **Usage**:

```go
// [...]

hobbies := []string{"swimming", "coding", "gaming", "playing"}
value, err := prompts.SelectBox("Choose your hobbies : ", hobbies)

if err != nil {
	log.Fatal(err)
}

fmt.Println("gender : ", value)
```
- **Result**:

  ![image](https://github.com/yassinebenaid/prompts/assets/101285507/192424db-c9fe-480d-9f9b-4f8d6cdff56b)

#

### Alerts
these are helper apis you can use for better alerts and messages.
```go
prompts.Info("Info alert")
prompts.Error("Error alert")
prompts.Warning("Warning alert")
prompts.Success("Success alert")
```

- **Result**:
  
![Screenshot from 2023-09-29 18-18-21](https://github.com/yassinebenaid/prompts/assets/101285507/d380f45b-4db4-4a9d-b4a2-8edf5f518579)



  
```go
prompts.InfoMessage("Info message")
prompts.ErrorMessage("Error message")
prompts.WarningMessage("Warning message")
prompts.SuccessMessage("Success message")
```

![Screenshot from 2023-09-29 18-19-17](https://github.com/yassinebenaid/prompts/assets/101285507/90fc7d67-f0fe-4608-9f57-b5a0852129f7)
