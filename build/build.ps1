param(
    [string]$Os
)


$o = "{0}" -f $Os

$Env:GOOS = $o; go build -o medialog ../medialog.go