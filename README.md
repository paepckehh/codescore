# OVERVIEW

[paepche.de/codescore](https://paepcke.de/codescore)

- Calculate a code quality score
- Based on number of golang (idomatic) code style violations per lines of code
- 100 % golang, minimal (external) imports, use as app or api (see api.goo)

# INSTALL

```
go install paepcke.de/codescore/cmd/codescore@latest
```

### DOWNLOAD (prebuild)

[github.com/paepckehh/codescore/releases](https://github.com/paepckehh/codescore/releases)

# SHOWTIME 

## Summary one-liner!

```Shell 
codescore .
[score:100/100] [loc:669] [err:0] -> [/usr/store/dev/codescore]

```

## Details why?

```Shell 
codescore  --verbose .
[score:099/100] [loc:1752] [err:40] -> [/usr/store/dev/asn2pf]
/usr/store/dev/asn2pf/asnfetch/api.go:48:4: if c { ... } else { ... continue } can be simplified to if !c { ... continue } ...
/usr/store/dev/asn2pf/asnfetch/api.go:71:6: explicit call to the garbage collector
[...]

```

## More?

```Shell 
codescore  --help
syntax: codescore [options] <file|directory>

--file [-f]
		create .go.codescore info file

--file-full [-F]
		create .go.codescore info file and dump all details into the file

--score-only [-s]
		print only the score to stdout

--enable-hidden-files [-e]
		enable scanning hidden files and directories

--exclude [-e]
		exclude all directories matching any of the keywords
		this option can be specified several times

--verbose [-v]
--silent [-q]
--debug [-d]
--help [-h]
```

# DOCS

[pkg.go.dev/paepcke.de/codescore](https://pkg.go.dev/paepcke.de/codescore)

# CONTRIBUTION

Yes, Please! PRs Welcome! 
