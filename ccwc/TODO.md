# Requirements and Bright Ideas

## Encodings

Handle locale/encodings: https://pkg.go.dev/github.com/delthas/go-localeinfo.
`wc` _does_ consider `LC_CTYPE`, e.g. `LC_CTYPE=en_US.UTF-8`. 

UTF-16?

https://stackoverflow.com/questions/36550038/in-utf-16-utf-16be-utf-16le-is-the-endian-of-utf-16-the-computers-endianness

## Buffer Overflow

At the moment, a "line" of text longer than 64K bytes will overflow the buffer used by Go's `bufio.ReadString()`.