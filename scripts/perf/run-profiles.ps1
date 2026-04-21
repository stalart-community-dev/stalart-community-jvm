param(
  [ValidateSet("synthetic", "real", "both")]
  [string]$Mode = "both",
  [string]$ConfigDir = "configs",
  [string]$OutDir = "artifacts/perf",
  [string]$RealCsv = "artifacts/perf/real-runs.csv"
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "../..")
Set-Location $repoRoot

$timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
$runDir = Join-Path $OutDir $timestamp
New-Item -ItemType Directory -Force -Path $runDir | Out-Null

Write-Host "Running profile benchmark harness..."
Write-Host "Mode: $Mode"
Write-Host "Output: $runDir"

python "scripts/perf/parse-results.py" `
  --mode $Mode `
  --config-dir $ConfigDir `
  --real-csv $RealCsv `
  --out-dir $runDir

python "scripts/perf/report.py" `
  --input-dir $runDir `
  --output-md (Join-Path $runDir "report.md")

Write-Host "Done. Report:"
Write-Host (Join-Path $runDir "report.md")
