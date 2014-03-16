gojson.com
==========

The site code to [gojson.com](http://gojson.com)

![](http://i.imgur.com/PtgJH4k.png)

This is a site that takes JSON strings like `{"name":"Hello"}` and produces golang structs for use in json.unmarshal.

example output:

```json

{"name":"Hello"}

```

Tunes into 

```go

package main

type Foo struct {
	Name string `json:"name"`
}


```
