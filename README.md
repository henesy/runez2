# Runez2

Another absolutely awful, but (probably) lossless, compression format for utf-8 text.

This is the successor to the [runez](https://github.com/henesy/runez) archive format.

See [the spec](./spec.md) for implementation details, or read the source ☺.

## Build

	; go build

## Usage

	runez2 [-D] [-c | -d] < input > output

## Examples

	; wc -c mac.txt
	3550 mac.txt
	; 9 wc -r mac.txt
	   1950 mac.txt
	; # ^ Number of runes
	; ./runez2 < mac.txt > mac.rz2
	; wc -c mac.rz2
	2030 mac.rz2
	; ./runez2 -d < mac.rz2 > newmac.txt
	; diff mac.txt newmac.txt
	;

## References

- [runez](https://github.com/henesy/runez)
