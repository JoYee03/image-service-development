# upload iCARES.jpg to image/Watermark/iCARES.jpg
$wmBase64 = Get-Content -Raw "E:\Internship_Test\iCARES\icares_base64.txt"
$body2 = @{
    content = "$wmBase64"
    type = "image/jpeg"
    filename = "image/iCARES.jpg"
} | ConvertTo-Json -Compress

Write-Host "`Uploading iCARES.jpg..."
$response2 = Invoke-RestMethod -Uri https://image-service-development-735683043266.asia-southeast1.run.app/testImageUpload -Method Post -Body $body2 -ContentType "application/json"
Write-Host "iCARES.jpg upload result:"
$response2

# apply the watermark to Blank.jpg using iCARES.jpg
$body3 = @{
    image_path = "image/202505/Blank.jpg"
    watermark_path = "image/iCARES.jpg"
} | ConvertTo-Json -Compress

Write-Host "`Applying watermark..."
$response3 = Invoke-RestMethod -Uri https://image-service-development-735683043266.asia-southeast1.run.app/testWatermarkImage -Method Post -Body $body3 -ContentType "application/json"
Write-Host "Watermark result:"
$response3
