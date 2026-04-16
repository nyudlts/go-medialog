param(
    [string]$Os,
    [string]$Od,
    [string]$Path
)


$o = "{0}" -f $Os
$p = "{0}" -f $Path
$d = if ($o -eq 'windows') { "{0}.exe" -f $Od } else { "{0}" -f $Od }

$Env:GOOS = $o; go build -o $d $p