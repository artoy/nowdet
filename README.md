# Nowdet

> [!WARNING]
> This repostitory is experimental and still under development.
> Now, you can check when **`time.Now()` is called in the same function which calls `spanner.Insert` or `spanner.Update`**.

Nowdet is a tool to detect `time.Now()` in your code is inserted to `allow_commit_timestamp` column in Cloud Spanner.

## Installation

```bash
go install github.com/nowdet/nowdet@latest
```

## Usage
Please specify a function name you want to check and a package name which contains it.

```bash
nowdet --pkg <package name>  --func <function name>
```
