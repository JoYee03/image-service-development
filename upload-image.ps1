$blankBase64 = Get-Content -Raw "E:\Internship_Test\iCARES\blank_base64.txt"
$body = @{
    content = "$blankBase64"
    type = "image/jpeg"
    filename = "image/202505/Blank.jpg"
} | ConvertTo-Json -Compress

Write-Host "Uploading Blank.jpg..."
$response = Invoke-RestMethod -Uri https://image-service-development-735683043266.asia-southeast1.run.app/testImageUpload -Method Post -Body $body -ContentType "application/json"
Write-Host "Upload result:"
$response
