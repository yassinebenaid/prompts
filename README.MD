# Prompts

Build command line prompts with ease, prompts provides several apis to help you create intuitive TUI faster,

## API

### Input

The input api allows you to prompt the user for an input field , it returns the input value

- usage:

```go
value, err := prompts.InputBox(prompts.InputOptions{
	Secure:      false, // hides the user input, very common for passwords
	Label:       "what is your name?",
	Placeholder: "what is your name",
	Required:    true,
	Validator:   func(value string) error { return nil },// will be called when user submit, and returned error will be displayed to the user below the input
})
```

- result
