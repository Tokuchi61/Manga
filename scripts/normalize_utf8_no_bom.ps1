param(
    [string[]]$Paths = @('README.md', 'docs/*.md')
)

$ErrorActionPreference = 'Stop'
$utf8Strict = New-Object System.Text.UTF8Encoding($false, $true)
$utf8NoBom = New-Object System.Text.UTF8Encoding($false)
$latin = [System.Text.Encoding]::GetEncoding(1252)
$turkish = [System.Text.Encoding]::GetEncoding(1254)

function Get-MojibakeScore([string]$text) {
    if ([string]::IsNullOrEmpty($text)) {
        return 0
    }
    return [regex]::Matches($text, '[ÃÂÅÄ??]').Count
}

function Repair-Mojibake([string]$text) {
    if ([string]::IsNullOrEmpty($text)) {
        return $text
    }

    for ($i = 0; $i -lt 3; $i++) {
        $score = Get-MojibakeScore $text
        if ($score -eq 0) {
            break
        }

        $bestText = $text
        $bestScore = $score

        foreach ($encoding in @($latin, $turkish)) {
            try {
                $candidate = $utf8Strict.GetString($encoding.GetBytes($text))
                $candidateScore = Get-MojibakeScore $candidate
                if ($candidateScore -lt $bestScore) {
                    $bestText = $candidate
                    $bestScore = $candidateScore
                }
            }
            catch {
            }
        }

        if ($bestScore -ge $score) {
            break
        }

        $text = $bestText
    }

    return $text
}

$files = @()
foreach ($pattern in $Paths) {
    $files += Get-ChildItem -Path $pattern -File -ErrorAction SilentlyContinue
}
$files = $files | Sort-Object -Property FullName -Unique

foreach ($file in $files) {
    $bytes = [System.IO.File]::ReadAllBytes($file.FullName)

    try {
        $text = $utf8Strict.GetString($bytes)
    }
    catch {
        $text = $turkish.GetString($bytes)
    }

    $text = Repair-Mojibake $text
    [System.IO.File]::WriteAllText($file.FullName, $text, $utf8NoBom)
}

Write-Output "UTF-8 no-BOM normalize tamamlandi: $($files.Count) dosya"
