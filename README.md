# mft
some tools (golang)
## ID Generator  
struct `*mft.G` - generator next id  
field `G.AddValue` - adds value that adds to each generated id use it with `*G.RvGetPart()` method  
valid values for `G.AddValue`: 10 - 9990 step 10  
method `*G.RvGetPart()` - returns ids as time.Now().UnixNano() / 10000 * 10000 + `AddValue`[10-9990; s10] + tail[0-9]  
method `*G.RvGet()` - returns ids as time.Now().UnixNano() / 10000 * 10000 + `AddValue` + tail[unbounded]  
`*G.RvGet()` faster then `*G.RvGetPart()` however does not provide uniqueness between generators when using different `G.AddValue`  

#### Global generator  
`mft.GlobalGenerator.AddValue = 100` sets global AddValue as 100  
`mft.RvGet()` generates next id  
`mft.RvGetPart()` generates next id uniqueness between generators when using different `mft.GlobalGenerator.AddValue`  


## Error
You can create error with code and error from error  
#### Flags
`mft.FillCallStack = true` - enables call stack  

#### Error codes
Use `mft.AddErrorsCodes(m)` for append error codes  
Use `mft.GenerateError(code, params...)` for create error with code  
Use `mft.GenerateErrorE(m)` for create error with code from another error  
Use `mft.GenerateErrorSubList(m)` for create error with code from errors list  

#### Example:  
```
package main
import (
	"fmt"
	"github.com/myfantasy/mft"
)
func main() {
	var e error
	e = mft.ErrorCS(7, "Bond Error")
	fmt.Println(e)
	e = mft.ErrorCSE(12, "M Error", e)
	fmt.Println(e)
}
```
```
[7] Bond Error
[12] M Error    [7] Bond Error
```

