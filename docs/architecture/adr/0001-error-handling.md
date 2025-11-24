# ADR 0001: Error Handling Strategy

- Status: accepted
- Date: 2025-11-24

## Context

Need unified lightweight error semantics without forcing exceptions or heavy wrappers.

## Decision

Provide `optional.Optional` for presence, `result.Result` for (value,error) with invariant: empty value + nil error => internal sentinel error. `promise.Promise` propagates `error` and panics only on `Join()` if error present.

## Alternatives

- Use only `error` return tuples (verbose chaining).
- Introduce monadic error type everywhere (heavy adoption cost).

## Consequences

- Clear separation: presence (`Optional`), success/failure (`Result`), async (`Promise`).
- Occasional panic on misuse (intentional fail-fast).
