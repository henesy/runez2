# Runez2

The file extension for a runez2 archive is `.rz2`.

## Archive format

The archive is little-endian binary encoded in the form:

	[rune int32]
	…
	[\0 rune]
	[N uint8]
	…

Each rune is a valid utf-8 character, stored by Go as an int32(?).

The `\0` is a null rune separating the preamble from the set of indices.

Each uint8 index `N` is referring to an index position `N` within the ordered table of runes at the beginning of the archive.

## Algorithm

The general conversion looks like:

	αβξαβξ

to

	α β ξ
	\0
	0 1 2 0 1 2

## Restrictions

We assume:

- The whole file is read into memory
- There are no more than `^uint8(0)` unique runes
