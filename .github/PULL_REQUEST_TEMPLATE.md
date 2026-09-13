## Summary

<!-- What does this change and why? Link issues with "Fixes #123". -->

## Definition of Done

- [ ] Tests cover normal, boundary (0, unit edges, rounding promotion), negative and extreme (`math.MinInt64`, `math.MaxInt64`, NaN/Inf, invalid UTF-8) inputs
- [ ] Coverage stays at 100%
- [ ] Every exported identifier has a GoDoc comment with examples of real output
- [ ] Runnable `Example` function added or updated
- [ ] Benchmark added for new formatters (`-benchmem`, allocations noted)
- [ ] Fuzz target added or extended where input is arbitrary (strings, full numeric range)
- [ ] README updated if the public API changed
- [ ] `CHANGELOG.md` updated under `[Unreleased]`
- [ ] `make check` passes locally
- [ ] No new dependencies (`go.mod` has no `require`)
- [ ] Commit messages follow Conventional Commits

## Breaking changes

<!-- None, or describe the output/API change and migration. -->
