# Architecture

HTTP transport calls domain validation, repository interfaces and the deterministic solver. Solver never reads wall clock except for an explicit budget; tie-break order is priority then container ID then slot coordinates. Dispatch adapters are isolated behind an interface so simulation and production cannot share mutable state. Plan writes use a revision field as an optimistic lock and transitions are checked before persistence.
