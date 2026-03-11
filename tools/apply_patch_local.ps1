param(
    [Parameter(Position = 0)]
    [string]$Patch
)

$ErrorActionPreference = 'Stop'

$workspaceRoot = Split-Path -Parent $PSScriptRoot
$localCodexDir = Join-Path $workspaceRoot '.codex-temp'
$localCodexExe = Join-Path $localCodexDir 'codex.exe'
$localAppAsar = Join-Path $localCodexDir 'app.asar'
$sourceDir = 'C:\Program Files\WindowsApps\OpenAI.Codex_26.306.996.0_x64__2p2nqsd0c76g0\app\resources'

if (-not (Test-Path $localCodexDir)) {
    New-Item -ItemType Directory -Force -Path $localCodexDir | Out-Null
}

if (-not (Test-Path $localCodexExe)) {
    Copy-Item (Join-Path $sourceDir 'codex.exe') $localCodexExe -Force
}

if (-not (Test-Path $localAppAsar)) {
    Copy-Item (Join-Path $sourceDir 'app.asar') $localAppAsar -Force
}

if ([string]::IsNullOrEmpty($Patch)) {
    $Patch = [Console]::In.ReadToEnd()
}

if ([string]::IsNullOrWhiteSpace($Patch)) {
    throw 'Patch content is required.'
}

if (-not $Patch.TrimEnd("`r", "`n").EndsWith('*** End Patch')) {
    throw "Invalid patch: The last line of the patch must be '*** End Patch'"
}

& $localCodexExe --codex-run-as-apply-patch $Patch
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
