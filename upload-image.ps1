# upload-image into storage
$base64Image = Get-Content -Raw -Path "blank_base64.txt"

$body = @{
    content  = "$base64Image"
    type     = "image/jpg"
    filename = "image/2025/Blank.jpg"
} | ConvertTo-Json -Depth 3 -Compress

Invoke-WebRequest -Uri "https://image-service-development-735683043266.asia-southeast1.run.app/testImageUpload" `
  -Method POST `
  -Body $body `
  -ContentType "application/json"