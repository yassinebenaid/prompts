package goclitools

type Style string

const (
	Reset        Style = "\033[0m"
	Bold         Style = "\033[1m"
	Dim          Style = "\033[2m"
	Italic       Style = "\033[3m"
	Underline    Style = "\033[4m"
	BlinkSlow    Style = "\033[5m"
	BlinkRapid   Style = "\033[6m"
	ReverseVideo Style = "\033[7m"
	Concealed    Style = "\033[8m"
	CrossedOut   Style = "\033[9m"

	// Foreground text colors
	TextBlack   Style = "\033[30m"
	TextRed     Style = "\033[31m"
	TextGreen   Style = "\033[32m"
	TextYellow  Style = "\033[33m"
	TextBlue    Style = "\033[34m"
	TextMagenta Style = "\033[35m"
	TextCyan    Style = "\033[36m"
	TextWhite   Style = "\033[37m"
	TextDefault Style = "\033[39m"

	// Background colors
	BackgroundBlack   Style = "\033[40m"
	BackgroundRed     Style = "\033[41m"
	BackgroundGreen   Style = "\033[42m"
	BackgroundYellow  Style = "\033[43m"
	BackgroundBlue    Style = "\033[44m"
	BackgroundMagenta Style = "\033[45m"
	BackgroundCyan    Style = "\033[46m"
	BackgroundWhite   Style = "\033[47m"
	BackgroundDefault Style = "\033[49m"

	// High-intensity foreground text colors
	TextBrightBlack   Style = "\033[90m"
	TextBrightRed     Style = "\033[91m"
	TextBrightGreen   Style = "\033[92m"
	TextBrightYellow  Style = "\033[93m"
	TextBrightBlue    Style = "\033[94m"
	TextBrightMagenta Style = "\033[95m"
	TextBrightCyan    Style = "\033[96m"
	TextBrightWhite   Style = "\033[97m"

	// High-intensity background colors
	BackgroundBrightBlack   Style = "\033[100m"
	BackgroundBrightRed     Style = "\033[101m"
	BackgroundBrightGreen   Style = "\033[102m"
	BackgroundBrightYellow  Style = "\033[103m"
	BackgroundBrightBlue    Style = "\033[104m"
	BackgroundBrightMagenta Style = "\033[105m"
	BackgroundBrightCyan    Style = "\033[106m"
	BackgroundBrightWhite   Style = "\033[107m"

	Tab Style = "  "
)
