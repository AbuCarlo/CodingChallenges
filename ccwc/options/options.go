package options

type WcOptions struct {
	Chars   bool   `short:"m" long:"chars" description:"print the character counts"`
	Bytes   bool   `short:"c" long:"bytes" description:"print the byte counts"`
	Lines   bool   `short:"l" long:"lines" description:"print the newline counts"`
	Width   bool   `short:"L" long:"max-line-length" description:"print the maximum display width"`
	Words   bool   `short:"w" long:"words" description:"print the word count"`
	Help    bool   `long:"help" description:"display this help and exit"`
	Version bool   `long:"version" description:"output version information and exit"`
}

func (o WcOptions) IsDefault() bool {
	return !o.Bytes && !o.Chars && !o.Lines && !o.Width && !o.Words
}