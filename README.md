# MorkleRork

MorkleRork is a programming language invented to be very very easy to parse, lex, and execute, while having _some_ modern features.

It was invented so I could learn golang, and therefore this repo contains a MorkleRork interpreter written in golang

## Language Introduction

MorkelRork programs are constructed of lines of code

Each line can contain exactly one **Command**, and each **Command** is contained entirely on one line of code

Each **Command** is made up of **Symbols**

**Symbols** are always separated by 1 space

Each **Command** starts with 1 **CommandSymbol**, then any number of other **Symbols** relating to its operation

This means **Commands** can be lexed and parsed line by line, and commands can by lexed and parsed symbol by symbol

Note: The one exception to this rule is string literals, which can have spaces in them. For example "Hello World!" should be lexed as a single string literal

Note: A blank line is also allowed in MorkelRork, and should simply be discarded by the Lexer

For example, take the **CommandSymbol** `log`, which prints to stdout the result of the provided **Expression**

Given the **IntLiteral** **Symbol** 5 like so:

```morklerork
log 5
```

This **Command** will print the number 5 to stdout

**Expressions** can be composed of chains of **Expression Symbols** and **Operators**

For example, to log the number 5 followed by a new line you would have:

```morklerork
log '' + 5 + '\n''
```

since the `log` **Command** supports one **Expression**, it will first evaluate the **Expression**, then output the resulting Value to stdout.

Finally, in MorkleRork indentation defines **Blocks** of code. Some commands expect a **Block** after them.

**Blocks** are defined until the indentation reduces back to its original level.

For example, this if **Command** is associated with a **Block** comprised of two **Commands**

```morklerork
if 5 == :num
  log 'Hello! :num is 5\n'
  log 'I\'m still in this block\n'
log 'I\'m out of the block\n'
```

Note: For simplicity of lexing, you can consider the leading whitespace (that defines **Blocks**) to be a Symbol of its own

## Symbols

As mentioned, each symbol can be ascertained by splitting a line on spaces (Taking care of string literals as the 1 special case)

### Built In Symbols

These symbols are 'reserved' by the language, they are all discussed below in their relevant sections

MorkleRork has 8 **CommandSymbols** (and therefore only 8 possible **Commands**), 9 **OperatorSymbols**, three types of **LiteralSymbol**, and two types of **UserDefinedSymbols**

#### CommandSymbols
`log new = if while program call return`

These are explained below in the `Commands` section

#### OperatorSymbols
`& | == != < + - * /`

These are explained below in the `Operators` section

### Literals

Literals let you directly insert values into your code

#### Int Literals:
any symbol composed only of numbers. E.G. this line should be lexed as 5 IntLiterals:

`1 143 13426534637 2 4`

#### String Literals:
everything between an opening `'` and a closing `'`. E.G.

`''` empty string

`'Hello World!'` spaces are allowed

#### Boolean Literals:
`?true ?false`

### User Defined Symbols

These symbols are created using certain **Commands** and can be used by other **Commands**

#### Variable Names
Any string of characters starting with `:` E.G.

`:var :num :thing1 :thing2`

#### Program Names
Any string of characters starting with `$` E.G.

`$fib $power $max`

## HeapAccessSymbols

MorkleRork comes with an infinite heap of Values (the same way turing machines do, your interpreter/compiler may provide a limited heap. This can be visualized as an array of boxes that can be any type.

You can access a box using `[` and `]` with a SingleExpression between that evaluates to an Int value. E.G.

`[2]` access the box numbered 2

`[:index]` access the box numbered by the int in the variable `:index`

`[[2]]` access the box numbered by the int in the box numbered 2

HeapAccess can read, or write the value of the box depending on the context

Note: boxes are numbered from 0 up, inclusive

Note: Heap access can neven have a space between `[` and `]`, since they need to be lexed as a single **Symbol**

## Expressions

Expressions are combinations of literals, Variables, and HeapAccess using operators.

Sometimes **Commands** call for a 'SingleExpression' which means 'Expressions with only one symbol in them', I.E. one literal, one Variable, or one HeapAccess

E.G.

`1 + 2 + 3` adding an arbitrary group of numbers

`1 + 2 * 3` morkleRork follows BODMAS, so here 2 * 3 will be evaluated first

`'hello ' + 'world'` you can add strings, for concatenation

`'Your name is ' + :name` adding a Variable to a string

## Commands

### log

```morklerork
log <Expression>
```

The log command will evaluate the expression then output it into stdout

Examples:
```morklerork
log 5
log 'hello\n'
log 'Your Variable Value Is: ' + :var + '\n'
```


### new

```morklerork
new <VariableName> <Expression>
```

The new command will create a new variable in the current scope by the given name, with the value resulting from the expression.

If you reference a variable before it is created with `new` you will get an error at runtime

Examples:

```morklerork
new :num 5
new :string "hello\n"
new :isEnabled ?false
```

Since `new` will define a variable for in the current scope, and a **Block** always results in a new scope, you can 'mask' VariableNames in 'higher scope' if you need to. These Variables will be dis-guarded when the **Block** completes

Example
```morklerork
new :index 5
if :index < 10
    new :index 'Hello '
    log :index
log :index
```
This code will log `Hello 5`

### = (assign)

```morklerork
= <Variable|HeapAccess> <Expression>
```

The = command will set the target Variable or HeapCell to the value resulting from the expression.

Examples:

```morklerork
= :num 5
= :string 'hello\n'
= :isEnabled ?false
= [0] [5]
```

### if

```morklerork
if <Expression with Boolean value>
```

The if command will evaluate the expression, and if it gets the ?true value, it will execute the following block.

if the expressions is ?false the block, and all sub blocks will be skipped

Examples:


```morklerork
if 5 == :num
  log 'Hello! :num is 5\n'
  log 'I\'m only executed if the value of :num is 5\n'
log 'I\'m always executed\n'
```

### while

```morklerork
while <Expression with Boolean value>
```

The while command will evaluate the expression, and if it gets the ?true value, it will execute the following block.

Once the block is complete, the while command will re-evaluate the expression, and continue to execute the block as long as the expression resolves to ?true

if the expression resolves to ?false, the block will be skipped, even if its ?false the first time

Note: The `scope` created by the `while` command will be disgarded after each loop, and re-created if the **Block** is re-executed. This means you can use the `new` command freely in while **Blocks**

Examples:

```morklerork
new :num 0
while :num < 5
  log 'I will print 5 times\n'
  = :num :num + 1
```

### program

```morklerork
program <ProgramName> (<VariableName>)...
```

The program command takes a ProgramName, and any number of VariableNames (including none, for a program with no inputs)

It defines a new program, and immediately skips the entire block.

The block will be executed if the program is 'called' using the `call` command

The VariableNames are names to be given to variables passed to the program during the `call` command

These VariableNames are defined in the program scope only, and do not need `new` commands to use

Example:

Define a program for working out one number to the power of annother
```morklerork
program $pow :a :b
    new :ret 1
    while 0 < :b
        = :ret :ret * :a
        = :b :b - 1
    return :ret
```

### return

```morklerork
return (<Expression>)
```

When called within a program, it immediately exits the program block (even if called within a deeper block such as in an if statement)

if the expression is provided it will pass the value back to the `call` command, which will place its value in the return variable if one is provided

Define a program returning the minimum of two values
```morklerork
program $min :a :b
    if :a < :b
        return :a
    return :b
```


### call

```morklerork
call (<VariableName>) <ProgramName> (<SingleExpression>)...
```

Calls a program, passing any number of values to that program

If the first VariableName is included, it will receive any value returned from the program using the `return` command

Note: Call creates a completely new scope stack, so the program invoked cannot access any VariableNames defined in any upper scopes

Print the value of 2 to the power of 4 using a program
```morklerork
program $pow :a :b
    new :ret 1
    while 0 < :b
        = :ret :ret * :a
        = :b :b - 1
    return :ret

new :ret 0
call :ret $pow 2 4
log '' + :ret + '\n'
```

## Operators

MorkleRork only has 9 operators

Operators always operate on two values, operators are only valid for a subset of types

Operators of the same type execute left to right.

Operators obey the following precedence (high operators execute first)

```morklerork
0: "&"
1: "|"
2: "=="
3: "!="
4: "<"
5: "+"
6: "-"
7: "*"
8: "/"
```


### &

`<Bool> & <Bool>`

evaluates to the 'logical and' of the two values

### |

`<Bool> | <Bool>`

evaluates to the 'logical or' of the two values

### ==

`<Int> == <Int>`
`<String> == <String>`

evaluates to `?true` if the values are the same, or `?false` if they are not

### !=
`<Int> != <Int>`
`<String> != <String>`

evaluates to the opposite that `==` evaluates to

### <
`<Int> < <Int>`

evaluates to true if the left hand side is numerically lower than the right hand side

Note: There is no `>` in morklerork, as this can always be achieved by simply flipping the operands around

Note: there is no `<=` as this can be archived with either `:a < :b + 1` or `:a < :b | :a == :b`

### +
`<Int> + <Int>`

evaluates to the integer sum of the integers

`<String> + <String>`

evaluates to the concatenation of the two strings

`<String> + <Int>`

evaluates to the concatenation of the int (as a string of digits in base 10) to the string

Note: there is no `<Int> + <String>`, since you can always `"" + 1 + " bottle of beer on the wall"` to first get the 1 in a string 

### - * /
`<Int> - <Int>`
`<Int> * <Int>`
`<Int> / <Int>`

evaluates to the integer value of the mathematical operation

Note: `/` always rounds down

## System functions and Standard Library

System functions will all be prefixed with a `#`. They are implemented by the intepreter to do things like interfacing with the operating system

MorkleRork has a few more tricks up its sleeve that are coming soon, such as:
* Read a string, character by character into the heap
* Read characters one by one out of the heap into a string
* Read another morklerork file into the current file (to allow some level of code management)
* Receive stdin as a string
* Receive the cli args as a string
* Read a file as a string

Standard Library functions are implemented _in_ MorkleRork, and so are interpreter independent.

* A set of functions to initialize sections of the heap into a general memory allocation region supporting probably:
  * `call $heap$init <Int SingleExpression | heap start address> <Int SingleExpression | heap size>`
  * `call :address $heap$new <Int Expression | number of cells>`
  * `call $heap$free <Int Expression | address>`

* Maybe a set of functions to implement common collections, such as Stacks, Queues, Maps, etc:
  * `call :address $stack$create`
  * `call $stack$push :address :value`
  * `call :value $stack$pop :address`