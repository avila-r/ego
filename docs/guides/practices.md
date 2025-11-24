# Recommended Practices

- Prefer `Result` when both value & error matter; use `Optional` only for presence.
- Keep `Promise` chains shallow; compose via `Compose` for dependent async operations.
- Use `iterator.SliceIterator` only when stateful iteration is required; otherwise `slice.ForEach`.
- Export minimal surface from your own wrappers; re-export types sparingly.
