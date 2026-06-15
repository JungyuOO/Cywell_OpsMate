param(
    [string]$Namespace = "opsmate",
    [string]$Name = "sample",
    [string]$CookieSecretName = "opsmate-admin-oauth-cookie",
    [string]$CookieSecretKey = "session_secret",
    [string]$AdminRouteHost = "",
    [string]$AdminGroup = "cyops-admins",
    [string]$ManifestPath = "config/samples/opsmate_v1alpha1_opsmateconfig.yaml",
    [switch]$SkipApply
)

$ErrorActionPreference = "Stop"

function Invoke-Oc {
    param([Parameter(ValueFromRemainingArguments = $true)][string[]]$Arguments)
    & oc @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "oc $($Arguments -join ' ') failed with exit code $LASTEXITCODE"
    }
}

function New-CookieSecretValue {
    $bytes = New-Object byte[] 32
    [System.Security.Cryptography.RandomNumberGenerator]::Fill($bytes)
    return [Convert]::ToBase64String($bytes)
}

function Assert-NoSecretMaterial {
    param([string]$Text)
    $patterns = @(
        "postgres://",
        "password=",
        "token=",
        "CYOPS_POSTGRES_DSN",
        "CYOPS_ADMIN_TOKEN",
        "CYOPS_EMBEDDING_TOKEN"
    )
    foreach ($pattern in $patterns) {
        if ($Text -match [regex]::Escape($pattern)) {
            throw "sensitive value pattern appeared in smoke output: $pattern"
        }
    }
}

Write-Host "Checking OpenShift login context"
Invoke-Oc whoami | Out-Null

Write-Host "Ensuring namespace exists: $Namespace"
Invoke-Oc get namespace $Namespace | Out-Null

Write-Host "Creating/updating OAuth cookie secret: $CookieSecretName"
$cookie = New-CookieSecretValue
$secretYaml = New-TemporaryFile
try {
    & oc -n $Namespace create secret generic $CookieSecretName "--from-literal=$CookieSecretKey=$cookie" --dry-run=client -o yaml | Set-Content -LiteralPath $secretYaml.FullName -Encoding utf8
    if ($LASTEXITCODE -ne 0) {
        throw "could not render cookie secret yaml"
    }
    Invoke-Oc apply -f $secretYaml.FullName
}
finally {
    Remove-Item -LiteralPath $secretYaml.FullName -Force -ErrorAction SilentlyContinue
}

if (-not $SkipApply) {
    Write-Host "Applying OpsMateConfig manifest: $ManifestPath"
    Invoke-Oc apply -n $Namespace -f $ManifestPath
}

if ($AdminRouteHost -ne "") {
    Write-Host "Patching admin Route host: $AdminRouteHost"
    Invoke-Oc -n $Namespace patch opsmateconfig $Name --type merge -p "{`"spec`":{`"console`":{`"adminRouteHost`":`"$AdminRouteHost`",`"adminGroups`":[`"$AdminGroup`"],`"adminAuthProxyEnabled`":true,`"adminAuthProxyCookieSecretRef`":`"$CookieSecretName`"}}}"
}

$routeName = "$Name-admin-authproxy"
$jobName = "$Name-pgvector-migration"

Write-Host "Waiting for admin auth proxy Route"
Invoke-Oc -n $Namespace wait "--for=jsonpath={.spec.to.name}=$routeName" "route/$routeName" --timeout=120s
$host = (& oc -n $Namespace get route $routeName -o jsonpath="{.spec.host}")
if ($LASTEXITCODE -ne 0 -or $host -eq "") {
    throw "could not resolve admin auth proxy Route host"
}
Write-Host "Admin Route: https://$host"

Write-Host "Checking admin Route redirects to OpenShift OAuth"
$response = & curl.exe -k -I -s "https://$host/api/ops/diagnostics"
if ($LASTEXITCODE -ne 0) {
    throw "curl against admin Route failed"
}
$responseText = $response -join "`n"
if ($responseText -notmatch "HTTP/.* 30[12378]" -and $responseText -notmatch "oauth") {
    throw "admin Route did not look like an OAuth-protected redirect"
}
Assert-NoSecretMaterial $responseText

Write-Host "Waiting for pgvector migration Job if it exists"
$jobExists = $true
& oc -n $Namespace get "job/$jobName" *> $null
if ($LASTEXITCODE -ne 0) {
    $jobExists = $false
}

if ($jobExists) {
    Invoke-Oc -n $Namespace wait "--for=condition=complete" "job/$jobName" --timeout=300s
    $logs = & oc -n $Namespace logs "job/$jobName"
    if ($LASTEXITCODE -ne 0) {
        throw "could not read migration Job logs"
    }
    $logText = $logs -join "`n"
    Assert-NoSecretMaterial $logText
    $pgVectorReady = (& oc -n $Namespace get opsmateconfig $Name -o jsonpath="{.status.pgVectorReady}")
    if ($LASTEXITCODE -ne 0 -or $pgVectorReady -ne "true") {
        throw "OpsMateConfig status.pgVectorReady is not true after migration Job completion"
    }
    Write-Host "Migration Job completed and status.pgVectorReady=true"
} else {
    Write-Host "Migration Job $jobName does not exist; skipping migration smoke"
}

Write-Host "v0.0.22 OpenShift smoke checks completed"
