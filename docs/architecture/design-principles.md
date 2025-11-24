# Design Principles

1. Minimal Surface: Few exported symbols per package.
2. No Magic: Behaviors are explicit (e.g. `Promise.Join()` may panic on error).
3. Ergonomic Defaults: `slice.Of`, `list.Of`, `promise.Of` mirror familiar patterns.
4. Fail Fast: Illegal states panic early (invalid iterator usage, missing Optional).
5. Clear States: `promise.State` and `result.IsSuccess()` make flow visible.
6. Generics First: Prefer parametric polymorphism over interface{} or reflection.
7. Composability: Functions return plain values or light wrappers reusable across packages.
