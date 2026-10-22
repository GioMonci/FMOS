# Contributing to FMOS

1. Create a branch from `main`.
2. Keep changes focused, readable, and easy to review.
3. Before opening a pull request, build the project and check it:

   ```powershell
   gofmt -l .
   go vet ./...
   go build -o build/fmos.exe ./cmd/fmos
   ```

   `gofmt -l .` should print nothing.
4. Explain what changed and how it was checked.

Windows repair behavior must be tested carefully on a supported Windows system.
Never add an unexplained or automatically destructive system command.
