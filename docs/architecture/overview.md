# Architecture Overview

Ego groups small, orthogonal generic helpers:
- Value containers: `box`, `optional`, `result`
- Async orchestration: `promise`
- Collections & iteration: `collection`, `list`, `slice`, `iterator`, `stream`
- Env & config: `dotenv`, `env`
- HTTP utilities: `httpx`
- Misc helpers: `pair`, `pointer`, `constraint`

Core pattern: thin abstractions over plain Go types; no hidden goroutines except in `promise`. Generics keep type safety without reflection.
