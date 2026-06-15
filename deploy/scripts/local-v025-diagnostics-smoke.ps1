param(
    [string]$BaseUrl = "http://127.0.0.1:8080",
    [string]$AdminUser = "admin"
)

$ErrorActionPreference = "Stop"

function Invoke-Json {
    param([string]$Path)
    $headers = @{ "X-Forwarded-User" = $AdminUser }
    return Invoke-RestMethod -Uri "$BaseUrl$Path" -Headers $headers
}

$view = Invoke-WebRequest -Uri "$BaseUrl/console-plugin/diagnostics" -UseBasicParsing
if ($view.StatusCode -ne 200) {
    throw "diagnostics view returned $($view.StatusCode)"
}
if ($view.Content -notmatch "CYOps Diagnostics") {
    throw "diagnostics view did not include CYOps Diagnostics"
}

$script = Invoke-WebRequest -Uri "$BaseUrl/console-plugin/diagnostics.js" -UseBasicParsing
if ($script.Content -match "oauth") {
    throw "console diagnostics script contains oauth handling"
}

$diagnostics = Invoke-Json "/api/ops/diagnostics"
$schema = Invoke-Json "/api/ops/diagnostics/schema"

if ($diagnostics.ui.primaryEntry -ne "openshift-web-console") {
    throw "diagnostics primary entry is $($diagnostics.ui.primaryEntry)"
}
if (-not $schema.aggregateOnly) {
    throw "diagnostics schema is not aggregate-only"
}

Write-Host "v0.0.25 local diagnostics smoke passed"
