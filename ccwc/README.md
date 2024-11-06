# README for Anthony Nassar's reimplementation of the wc Unix command-line utility

The problem is specified here: https://codingchallenges.substack.com/p/coding-challenge-1

## Scope

The features of `wc` are described here: https://www.gnu.org/software/coreutils/manual/html_node/wc-invocation.html#wc-invocation,
as of 2024-11-05. `--total` was not available in `coreutils` 8.32, from 2020, and I did not implement it. I could have
implemented `--files0-from`, but the burden of integration-testing this option was too onerous.

The source for `wc` is at https://github.com/coreutils/coreutils/blob/master/src/wc.c.

Unfortunately, I chose to develop this code on a Windows laptop, using a version of GnuWin32 from 2005, which does
not even handle UTF-8 properly, so my clever integration-testing scheme (i.e. use a property-based testing framework to 
generate the same command lines for `wc` and `ccwc` and compare the respective output) fell apart immediately
on a document corpus including Latin-1 and Chinese characters in UTF-8 encoding.

### Licensing

Insofar as this is a "derivative" product in every sense of the GNU implementation of `wc`, it 
is offered under the same GNU license. Run `ccwc --version` for details.

### Character Encodings

I have deliberately replicated certain legacy features of `wc`, especially the original equation 
of bytes with characters. On a current Linux system, the encoding with almost certainly be UTF-8 
(i.e. `LANG=C.UTF-8`), which is the default for Go. No provision has yet been made for other character
encodings, or for BOM detection.
