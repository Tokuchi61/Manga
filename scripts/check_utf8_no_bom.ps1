param(
    [string[]]$Paths = @('README.md', 'docs/*.md')
)

$ErrorActionPreference = 'Stop'
$utf8Strict = New-Object System.Text.UTF8Encoding($false, $true)
$failed = @()

$files = @()
foreach ($pattern in $Paths) {
    $files += Get-ChildItem -Path $pattern -File -ErrorAction SilentlyContinue
}
$files = $files | Sort-Object -Property FullName -Unique

foreach ($file in $files) {
    $bytes = [System.IO.File]::ReadAllBytes($file.FullName)

    if ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) {
        $failed += "$($file.FullName): UTF-8 BOM bulundu"
        continue
    }

    try {
        [void]$utf8Strict.GetString($bytes)
    }
    catch {
        $failed += "$($file.FullName): UTF-8 gecersiz bayt dizisi"
    }
}

if ($failed.Count -gt 0) {
    $failed | ForEach-Object { Write-Error $_ }
    exit 1
}

Write-Output "UTF-8 no-BOM kontrolu basarili: $($files.Count) dosya"
