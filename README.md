# mft
some tools (golang)
G - generator next id


## Error
You can create error with code and error from error  

Example:  
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

