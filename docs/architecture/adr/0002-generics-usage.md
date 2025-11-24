# ADR 0002: Generics Usage

- Status: accepted
- Date: 2025-11-24

## Context
Need type-safe helpers (collections, containers, async) without runtime reflection overhead.

## Decision
Use generics broadly where parametric polymorphism is natural: `Box[T]`, `Optional[T]`, `Result[T]`, `Promise[T]`, `ArrayList[T]`, `SliceIterator[T]`, `Stream[T]`. Constrain only where required (`constraint.Comparable`). Avoid exposing complex type parameters to callers unnecessarily.

## Alternatives
- Rely on interface{} + casting (unsafe, verbose).
- Heavy reflective operations (slower, brittle).

## Consequences
- Compile-time safety.
- Small binary impact.
- Clear intent through type signatures.
