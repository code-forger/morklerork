# Standard Library

The standard library is composed of modules, All functions within those libraries that are meant for use are prefixed with the module name surrounded in `$` such as `$string$length` being a function in the `string` module

Note: since all functions and variables are exported from a file in MorkleRork, the library modules will also export other functions, you should avoid calling undocumented functions unless you really understand them.

## The `$heap$` module

The $heap$ module allows you to pass some portion of the morklerork heap to be 'managed'. This allows you to call $heap$new to allocate a new contiguous portion of memory, and $heap$free to de-allocate it.

Note: The heap has no memory protection system, so if you pass an incorrect address to any of its function, it will likely stomp memory, and maybe get itself in an infinite loop.

Furthermore, if you write into a cell not allocated to you, you risk destroying the heaps internal structure

### Exported Variables
:NULL_PTR

A value the $heap$ module uses to indicate invalid addresses

### Exported Functions

The standard library comes with 3 functions for creating and using a 'C like' heap manager

#### $heap$init
```morkleRork
call $heap$init <Int SingleExpression | heapStartAddress> <Int SingleExpression | heapSize>
# No Return
# heapStartAddress: The first address you want the managed heap to use, useful if you are using some of the cells for your own global storage
# heapSize: The number of cells to give to the managed heap
```

This function will initialize a set of cells as a 'heap'.

Note: init immediately consumes 6 cells, and in general every allocation you perform in the heap will consume 3 more cells. This is the overhead of a dynamic memory management system

This function must be called before $heap$new or $heap$free

#### $heap$new
```morkleRork
call <Int | Address> $heap$new <Int SingleExpression | heapStartAddress> <Int SingleExpression | size>
# returns: A heap address of the first cell allocated to you
# heapStartAddress: The same start address as passed to $heap$init, This allows you to init and interact with multiple managed heaps
# size: the number of contiguous cells you want
```

This function will allocate some contiguous memory to you, and return it to you

#### $heap$free
```morkleRork
call $heap$new <Int SingleExpression | heapStartAddress> <Int SingleExpression | address>
# No Return
# heapStartAddress: The same start address as passed to $heap$init, This allows you to init and interact with multiple managed heaps
# address: The address that was returned to you from $heap$new
```

This function will free some memory, allowing it to be re-allocated later

## The `$string$` module

This module provides helper functions for working with strings

### Exported Functions

#### $string$length
```morkleRork
call <Int | Length> $heap$init <String | Input String>
# returns the length of the string
```

This function returns the length of the int

#### $string$toInt
```morkleRork
call <Int | Int from string> $string$toInt <String | Input String>
# returns the string converted to an int, non-numerical characters will be replaced with 0's
```

This function parses ints in the string, non-numerical characters will be replaced with 0's

### Future plans the Standard Library
MorkleRork has a few more tricks up its sleeve that are coming soon, such as:
* Receive stdin as a string
* Receive the cli args as a string
* Read a file as a string

* Maybe a set of functions to implement common collections, such as Stacks, Queues, Maps, etc:
  * `call :address $stack$create`
  * `call $stack$push :address :value`
  * `call :value $stack$pop :address`