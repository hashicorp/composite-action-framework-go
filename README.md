# Composite Action Framework Go

A library of packages for writing _composite_ GitHub Actions in Go.

---

_This is intended for internal HashiCorp use only and is not generally supported for other users._

---

This is different from other Go GitHub Actions libraries which are focussed on
writing standalone actions that just run a go binary and then exit.

## Why Composite Actions?

Composite actions allow the composition of multiple other actions as
well as arbitrary build steps together in a single action. This means
that composite actions have the full actions ecosystem available.

Composite actions can also present some challenges:

- Need to be able to run discrete chunks of logic in separate steps.
- Need to be able to share configuration and calculated values between steps.
- When complex logic is involved, need to be able to test it and share it
  with other codebases.

## The Implied Strategy

This framework is intended to help you to write Go-based composite actions
by implementing a CLI to do the heavy lifting, and embedding calls to that
CLI in the action.yml.

The CLI package `./pkg/cli` is designed so that
each defined command is able to act as the entire CLI, or is able to be
embedded in another similar CLI and become a subcommand there.

This means that functionality written to support an action can be embedded
inside other tools, to be run locally, for instance.

# Release Process

1. Create a tag on the default branch of the commit to release.  Follow SemVer
   semantics when choosing versions.
1. Create a release from the tag, either with the GitHub web UI or the gh CLI.
