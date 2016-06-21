# godispatch
> Simple golang package that dispatches work to workers 
> Inspired by [http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/)

## Features
* Master Worker Map - Each Worker (Value) is tagged to a Master (Key) as a Key-Value pair in a map
* Provides sequential processing 

## Example
* Refer to test file