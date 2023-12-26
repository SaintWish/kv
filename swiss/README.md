# swiss
Fork of the swiss map [here](https://github.com/dolthub/swiss) with a few minor changes to work with some of my packages.

List of said minor changes.
* CHANGED ``(Map) Get(key any) (value any)`` - Now only returns a value.
* ADDED ``(Map) GetHas(key any) (ok bool, value any)`` - Returns value and ok if said key exists.
* CHANGED ``(Map) Set(key any, value any)`` - Renamed method Put to Set.
* CHANGED ``(Map) Delete(key any) (ok bool, value any)`` - Now returns the old value if successful.
* ADDED ``(Map) MaxCapacity() (int)`` - Returns the max capacity before the map needs to resize.