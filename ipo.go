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
	T_Black   Style = "\033[30m"
	T_Red     Style = "\033[31m"
	T_Green   Style = "\033[32m"
	T_Yellow  Style = "\033[33m"
	T_Blue    Style = "\033[34m"
	T_Magenta Style = "\033[35m"
	T_Cyan    Style = "\033[36m"
	T_White   Style = "\033[37m"
	T_Default Style = "\033[39m"

	// Background colors
	BG_Black   Style = "\033[40m"
	BG_Red     Style = "\033[41m"
	BG_Green   Style = "\033[42m"
	BG_Yellow  Style = "\033[43m"
	BG_Blue    Style = "\033[44m"
	BG_Magenta Style = "\033[45m"
	BG_Cyan    Style = "\033[46m"
	BG_White   Style = "\033[47m"
	BG_Default Style = "\033[49m"

	// High-intensity foreground text colors
	T_BrightBlack   Style = "\033[90m"
	T_BrightRed     Style = "\033[91m"
	T_BrightGreen   Style = "\033[92m"
	T_BrightYellow  Style = "\033[93m"
	T_BrightBlue    Style = "\033[94m"
	T_BrightMagenta Style = "\033[95m"
	T_BrightCyan    Style = "\033[96m"
	T_BrightWhite   Style = "\033[97m"

	// High-intensity background colors
	BG_BrightBlack   Style = "\033[100m"
	BG_BrightRed     Style = "\033[101m"
	BG_BrightGreen   Style = "\033[102m"
	BG_BrightYellow  Style = "\033[103m"
	BG_BrightBlue    Style = "\033[104m"
	BG_BrightMagenta Style = "\033[105m"
	BG_BrightCyan    Style = "\033[106m"
	BG_BrightWhite   Style = "\033[107m"

	Tab   Style = "  "
	UP    Style = "\x1b[1A"
	Clear Style = "\x1b[2K"
)
