# shortuuid

The go edtion of [skorokithakis‘s shortuuid](https://github.com/skorokithakis/shortuuid).

## Installation

You can use `go get` to install or `git clone`

```shell
go get -u github.com/zhengxiaowai/shortuuid
```

## Usage

First, you can import this package in your project, like so:

```golang
import "github.com/zhengxiaowai/shortuuid"
```

And, you can use `NewShortUUID()` to create shortuuid instance, which call `UUID()` to get a short uuid.

```golang
su := shortuuid.NewShortUUID()
fmt.Println(su.UUID()) // eMRLKa8f2jqFwRxqH7HjV2
```

Support UUID Version 5. if you provide URL or DNS，you will call `GetUUIDWithNameSpace`， like so:

```golang
su := shortuuid.NewShortUUID()
fmt.Println(su.GetUUIDWithNameSpace("example.com"))
fmt.Println(su.GetUUIDWithNameSpace("https://example"))
```

If the default 22 digits are too long for you, `Random()` method can truncate by length param：

 ```golang
 su := shortuuid.NewShortUUID()
 fmt.Println(su.Random(5)) // KLFHn
 ```
The secure random string is from `crypto/rand`'s Read function. 

If you want to use your own alphabet to generate UUIDs, use SetAlphabet():

> default alphabet is "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

```golang
su := shortuuid.NewShortUUID()
su.SetAlphabet("0123456789")
fmt.Println(su.UUID()) // 083961560630356262206399931541464041895
```

## Example

```golang
package main

import (
	"fmt"
	"github/zhengxiaowai/shortuuid"
)

func main() {
	su := shortuuid.NewShortUUID()
	fmt.Println(su.UUID())
	fmt.Println(su.Random(5))
	fmt.Println(su.GetUUIDWithNameSpace("example.com"))
	fmt.Println(su.GetUUIDWithNameSpace("https://example"))

	su.SetAlphabet("0123456789")
	fmt.Println(su.UUID())
	fmt.Println(su.Random(5))
	fmt.Println(su.GetUUIDWithNameSpace("example.com"))
	fmt.Println(su.GetUUIDWithNameSpace("https://example"))
}
```

## Warning

shortuuid only generate unique id in stand-alone.

if provide length(< 22) for `Random` function, The IDs won't be universally unique any longer, but the probability of a collision will still be very low. 

## License

BSD License