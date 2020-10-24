# mft
some tools (golang)

## RWCMutex
RWCMutex is RW mutex with `TryLock` method.  
TryLock try locks mutex until context done. TryLock returns true when success and false when not.  
LockD try locks mutex until during input duration.  

## PMutex
PMutex is RW mutex with `TryLock` method and you can `Promote` RLock to Lock and `Reduce` from Lock to RLock.  
PMutex works slowly then RWCMutex  

PMutex.Lock method returns lock key  
You must not lose this key, because you need it to unlock.  
You can use `TryUnlock` multiple times.

```
	k := mx.Lock()
    ...
	mx.TryUnlock(k)
    ...
	mx.TryUnlock(k)
```
```
	k := mx.Lock()
    ...
	mx.Unlock(k)
```


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

