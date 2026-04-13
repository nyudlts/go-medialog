param(
    [string]$Os,
    [string]$Od,
    [string]$Path
)


$o = "{0}" -f $Os
$p = "{0}" -f $Path
$d = "{0}" -f $Od

$Env:GOOS = $o; go build -o $d $p