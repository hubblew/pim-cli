## Project Code style

* Prefer functions to methods when the receiver is not needed.
* When defining structs that implement interfaces, add `var _ Interface = (*Struct)(nil)` to ensure at compile time that
  the struct implements the interface.