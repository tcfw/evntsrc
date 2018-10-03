# Contributing guide
Pull requests are always welcome. By participating in this project, you agree to abide by the [code of conduct](https://github.com/tcfw/evntsrc/docs/CODE_OF_CONDUCT.md).

## Submitting changes

Please send a [Pull Request](https://github.com/tcfw/evntsrc/pull/new/master) with a clear list of what you've done. We can always use more test coverage. Please follow our coding conventions (below) and make sure all of your commits are atomic (one feature per commit).

Always write a clear log message for your commits. One-line messages are fine for small changes, but bigger changes should look like this:

```
$ git commit -m "A brief summary of the commit
> 
> A paragraph describing what changed and its impact."
```

## Coding conventions

Please follow the standards by running `golint`, `gofmt` and/or `go vet`