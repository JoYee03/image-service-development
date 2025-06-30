param (
    [string]$ImagePath = "",          # Firebase path (e.g. "images/test.jpg")
    [string]$WatermarkPath = "",      # Firebase path (e.g. "images/iCARES.jpg")
    [string]$Endpoint = "https://image-service-development-647092317027.asia-southeast1.run.app/testWatermarkImage"
)

# ===== INPUT VALIDATION =====
if (-not $ImagePath -or -not $WatermarkPath) {
    Write-Host "Error: Both -ImagePath and -WatermarkPath are required!" -ForegroundColor Red
    exit 1
}

# ===== REQUEST CONFIGURATION =====
$body = @{
    image_path     = $ImagePath
    watermark_path = $WatermarkPath
} | ConvertTo-Json -Depth 3

Write-Host "Watermarking '$ImagePath' with '$WatermarkPath'..." -ForegroundColor Cyan

# ===== EXECUTE API CALL =====
try {
    $response = Invoke-RestMethod -Uri $Endpoint `
        -Method Post `
        -Body $body `
        -ContentType "application/json" `
        -ErrorAction Stop

    Write-Host " Success! Watermarked image saved to:" -ForegroundColor Green
    $response | Format-List | Out-String | Write-Host -ForegroundColor White

    if ($response.watermarked_path) {
        Start-Process $response.watermarked_path
    }
}
catch {
    Write-Host " HTTP Error ($($_.Exception.StatusCode)):" -ForegroundColor Red
    $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
    $errorDetails | Format-List | Out-String | Write-Host -ForegroundColor Yellow

    Write-Host "`n Full Error Details:" -ForegroundColor DarkYellow
    $_
}