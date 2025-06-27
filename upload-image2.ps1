#upload watermark to storage
$base64Image = Get-Content -Raw -Path "icares_base64.txt"

$body = @{
    content  = "$base64Image"
    type     = "image/jpg"
    filename = "image/iCARES.jpg"
} | ConvertTo-Json -Depth 3 -Compress

Invoke-WebRequest -Uri "https://image-service-development-735683043266.asia-southeast1.run.app/testImageUpload" `
  -Method POST `
  -Body $body `
  -ContentType "application/json"