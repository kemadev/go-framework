# `cmd`

## Files in this directory

- Are placed in subdirectories named after the behavior they implement, most likely there is only one `main` subdirectory that is the sole entry point of the application, but there can be specific subdirectories for very specific behaviors
- Are applications entry points (`main.go`) files
- Are meant to import code from `pkg` and `internal` only, are not implementing any business logic by just providing an entry point to the application
