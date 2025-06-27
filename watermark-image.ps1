# apply watermark
$body = @{
    image_path     = "image/2025/Blank.jpg"
    watermark_path = "image/iCARES.jpg"
    filename       = "image/2025/watermarked-final.jpg"
} | ConvertTo-Json -Depth 3 -Compress

Invoke-RestMethod -Uri "https://image-service-development-735683043266.asia-southeast1.run.app/testWatermarkImage" `
  -Method POST `
  -Body $body `
  -ContentType "application/json"
