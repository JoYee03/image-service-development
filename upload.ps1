 param (
    [string]$Base64Data = "",          # Direct base64 string input
    [string]$Base64Path = "",         # Path to base64 file (alternative to $Base64Data)
    [string]$MimeType = "image/jpeg", # Default MIME type
    [string]$Filename = $(if ($Base64Path) { [System.IO.Path]::GetFileName($Base64Path) } else { "uploaded_file.jpg" }),
    [string]$Endpoint = "https://image-service-development-647092317027.asia-southeast1.run.app/testImageUpload"
)

# ===== INPUT VALIDATION =====
if (-not $Base64Data -and -not $Base64Path) {
    Write-Host " Error: Either -Base64Data or -Base64Path must be provided!" -ForegroundColor Red
    exit 1
}

if ($Base64Path -and -not (Test-Path $Base64Path)) {
    Write-Host " Error: Base64 file not found at '$Base64Path'" -ForegroundColor Red
    exit 1
}

# ===== BASE64 HANDLING =====
if ($Base64Path) {
    $base64Image = (Get-Content -Raw -Path $Base64Path).Trim()
} else {
    $base64Image = $Base64Data.Trim()
}

if ($base64Image.Length -eq 0) {
    Write-Host " Error: Base64 data is empty!" -ForegroundColor Red
    exit 1
}

# ===== REQUEST CONFIGURATION =====
$body = @{
    content  = $base64Image
    type     = $MimeType
    filename = $Filename
} | ConvertTo-Json -Depth 3

Write-Host "Uploading '$Filename' ($MimeType) to $Endpoint..." -ForegroundColor Cyan

# ===== EXECUTE UPLOAD =====
try {
    $response = Invoke-RestMethod -Uri $Endpoint `
        -Method POST `
        -Body $body `
        -ContentType "application/json" `
        -ErrorAction Stop

    Write-Host " Success! File uploaded to:" -ForegroundColor Green
    $response | Format-List | Out-String | Write-Host -ForegroundColor White

    # Auto-open in browser if successful
    if ($response.public_url) {
        Start-Process $response.public_url
    }
}
catch {
    Write-Host " HTTP Error ($($_.Exception.StatusCode)):" -ForegroundColor Red
    $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
    $errorDetails | Format-List | Out-String | Write-Host -ForegroundColor Yellow
    
    Write-Host "`n Full Error Details:" -ForegroundColor DarkYellow
    $_
}

# ===== OUTPUT CLEANUP =====
if ($Base64Path) {
    Remove-Variable base64Image  # Clear sensitive data from memory
}