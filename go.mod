module github.com/dcilke/hz

go 1.18

require (
	github.com/francoispqt/gojay v1.2.13
	github.com/jessevdk/go-flags v1.5.0
	github.com/mattn/go-colorable v0.1.12
)

require (
	github.com/mattn/go-isatty v0.0.14 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

// Replace default gojay with github.com/dcilke/gojay fork for Decoder.ReadByte() goodness
replace github.com/francoispqt/gojay v1.2.13 => github.com/dcilke/gojay v1.2.14
