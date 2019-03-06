# MyHTTP

Tool which makes http requests and prints the address of the request along with the MD5 hash of the response or the error if any occured.

**Important**: order of the results is not guaranteed

## Getting started

### Install
```
go get github.com/i1skn/myhttp
```
### Run
```
$ myhttp -parallel 2 sorokin.io google.com golang.org

google.com 941cbb96297677b297ec926c2cc998f5
golang.org fca44a021e4861252e73c05fdbeb3f33
sorokin.io 62a7ec12812740ac2c6305a679c30bc0
```
#### Parameters
* `-parallel` - maximum number of the parallel requests


### Help
```
myhttp -help
```

## License

Apache License v2.0.
