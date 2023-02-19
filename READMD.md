# wc-go
wc-go is a Go implementation of the wc command-line utility, which counts the number of lines, words, and bytes in a file or input stream.

## Installation
To install wc-go, you must have Go installed on your system. Once Go is installed, you can clone this repository and build the binary yourself using the following commands:

```
git clone https://github.com/andrewkuehne/wc-go.git
cd wc-go
go build
```

## Usage
To use wc-go, run the wc-go command followed by one or more filenames. If no filenames are specified, wc-go will read from standard input. For example:

`wc-go file.txt`
This will output the number of lines, words, and bytes in file.txt.

You can also pipe input into wc-go using the standard Unix pipe syntax. For example:

`echo "hello world" | wc-go`
This will output the number of lines, words, and bytes in the string "hello world".

## License
wc-go is licensed under the Apache License 2.0. See the LICENSE file for details.

## Contributing
Contributions to wc-go are welcome! To contribute, please submit a pull request with your changes.

If you find any bugs or issues, please open an issue in the GitHub repository.

## Authors
`wc-go` was created by Andrew Kuehne.