Generics playground
---

This repository should not be used at all. It just contains my experiments with generics.

Notes:
---
"Do you want slow programmers, slow compilers and bloated binaries, or slow execution times?" (Russ Cox)

Let's say `result` is `nil`, then you will have to write : `is.Equal(result, []int(nil))`.

Sometimes, you need to return ugly things like `return *new(T), false`, just because you are passing strings to your generic functions.

You have to deref results when your generic function has to return pointers, like in the pattern `(*T, error)`.

There is no way to both support predeclared types AND support user defined types.

It would be great to have:

```go
type Swapper[T comparable] interface {
    Swap(i, j int)
}

type Ordered[T comparable] interface {
  constraints.Ordered | Swapper[T]
}
```

Implementation restriction: A union (with more than one term) cannot contain the predeclared identifier comparable or interfaces that specify methods, or embed comparable or interfaces that specify methods.

And the conclusions from [here](https://planetscale.com/blog/generics-can-make-your-go-code-slower) 
