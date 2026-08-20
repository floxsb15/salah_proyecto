param(
  [int]$BackendPort = 5000,
  [int]$FrontendPort = 3000
)

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot
$backendDir = Join-Path $root "backend-salah-main"
$frontendDir = Join-Path $root "front-salah-main"

$env:SERVER_PORT = "$BackendPort"
$env:GOCACHE = Join-Path $root ".gocache"
$env:NUXT_BACKEND_URL = "http://127.0.0.1:$BackendPort"
$env:PORT = "$FrontendPort"

Write-Host "Levantando backend en $env:NUXT_BACKEND_URL ..."
$backend = Start-Process -FilePath "go" -ArgumentList "run ." -WorkingDirectory $backendDir -PassThru -WindowStyle Hidden

try {
  Write-Host "Levantando frontend en http://127.0.0.1:$FrontendPort ..."
  Push-Location $frontendDir
  npm.cmd run dev
}
finally {
  Pop-Location
  if ($backend -and -not $backend.HasExited) {
    Stop-Process -Id $backend.Id -Force
  }
}
