# 01: Reject broken Clusters

**Ticket:** "Users keep creating Clusters with typos: pool names like `ML_Pool`, `minimum` bigger than `maximum`, no version. Reject them with clear messages that name the exact field."

**Run:** `make test EX=01`

## Task
Implement `ValidateCluster` and `ValidateClusterSpec` in `validation.go`. Report **every** problem at once (users hate fixing one error per apply), each with its field path:

- `metadata.name` is a DNS-1123 label
- `spec.version` is set
- at least one worker pool
- pool names: DNS-1123 label, at most 15 characters, unique
- `minimum >= 0`, `maximum >= minimum`

## Hints
- `validation.IsDNS1123Label(s)` (`k8s.io/apimachinery/pkg/util/validation`) returns a list of messages; empty means valid.
- Build paths with `fldPath.Child("workers").Index(i).Child("name")`.
- Error constructors: `field.Required(path, msg)`, `field.Invalid(path, value, msg)`, `field.Duplicate(path, value)`.
- A `map[string]bool` remembers which names you've seen.

## What you learn
- Structs, slices, loops, maps.
- Returning a list of errors instead of the first one.
- This is exactly how the API server tells you `spec.workers[0].maximum: Invalid value: 1: must be greater than or equal to minimum`.

## In Gardener
`pkg/api/core/validation/shoot.go`: thousands of lines of the same pattern. Open it and search for `field.Invalid`.
