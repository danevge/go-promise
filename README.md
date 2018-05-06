
Go promise like JS promise ([ECMA-262](http://www.ecma-international.org/ecma-262/6.0/index.html#sec-promise-abstract-operations))

## Installation
```
go get github.com/danevge/go-promise
```
## Quick start

```
value, err := 
     NewPromise(func(d interface{}) interface{} {return "parent"})).
     Then(func(d interface{}) interface{} { return d.(string) + " and first child" }).
     Then(func(d interface{}) interface{} { return d.(string) + " and second child" }).
     Get()
     
log.printf("value: %v", value)
log.printf("error: %v", err)
```     
Result:
```
value: parent and first child and second child
error: nil
```
## Features

### Create promise

Create new parent promise
```
NewPromise(func(d interface{}) interface{} { ... })
```

Create new child promise
```
.Then(func(d interface{}) interface{} { ... })
```

Add catch function
```
.Catch(func(err error) interface{} { ... })
```

If *.Catch* returned this and new error, promise will be considered fail. 
Error transfer in child's *.Catch* function and etc. 

However, if *.Catch* returned not error then promise will be success. 
New result transfer in child's.

Default catch function: 
```
.Catch(func(err error) interface{} { return err })
```

Create new child promise with catch

```
.ThenAndCatch(
   func(d interface{}) interface{} { ... },
   func(err error) interface{} { ... })
```

#### Return value 

Methods *NewPromise*, *.Then* and *.Catch* can return any value.

```
return 5
```
```
return "string"
```
```
return customObject
```
```
return nil
```
if returned new promise with new result then this result be use current promise
```
return NewPromise(...)
```
And error 
```
return fmt.Errorf("ups")
```

### Get result

Promise don't waiting method *.Get*. Promise start process while creating. 
However, child promise waiting result parent promise.

With default timeout (250ms):
```
value, err :=  NewPromise(...).Then(...).Get()
```
With custom timeout:
```
value, err :=  NewPromise(...).Then(...).GetWithTimeout(400 * time.Millisecond)
```

### Await async promises
```
All(
		func(d interface{}) interface{} { ... },
		func(d interface{}) interface{} { ... },
		func(d interface{}) interface{} { ... },
		...
	)
```
Result *.All* is new promise.

if one of inner promise return error then all promise will be considered fail.

Result work is result array.

### Race condition
```
Race (
		func(d interface{}) interface{} { ... },
		func(d interface{}) interface{} { ... },
		func(d interface{}) interface{} { ... },
		...
	)
```

Result work is result array
Result *.Race* is new promise.

Race promise return first success result.

### Simplify 

Function for promise can have simple construction. Examples:
```
		F(func(d interface{}) {...})
		F(func() interface{} {...} )
		F(func() {...})
		F("string")
		F(err error) 
		F(anyObject) 
```

