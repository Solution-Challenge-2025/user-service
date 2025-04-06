# Test script for user-service endpoints
$baseUrl = "http://localhost:8081/api/v1"

# Create a test file with content
Write-Host "Creating test file..."
$filePath = Join-Path $PSScriptRoot "test.txt"
$fileContent = @"
This is a test file
Created at: $(Get-Date)
Contains some sample content
Line 1
Line 2
Line 3
"@
$fileContent | Out-File -FilePath $filePath -Encoding UTF8
Write-Host "Test file created at: $filePath"

# Test file upload using multipart/form-data
Write-Host "`nTesting file upload..."
$fileBytes = [System.IO.File]::ReadAllBytes($filePath)
$fileEnc = [System.Text.Encoding]::GetEncoding('ISO-8859-1').GetString($fileBytes)
$boundary = [System.Guid]::NewGuid().ToString()
$LF = "`r`n"

$bodyLines = @(
    "--$boundary",
    "Content-Disposition: form-data; name=`"file`"; filename=`"test.txt`"",
    "Content-Type: text/plain",
    "",
    $fileEnc,
    "--$boundary--"
) -join $LF

try {
    $headers = @{
        "Content-Type" = "multipart/form-data; boundary=$boundary"
    }
    
    $fileUploadResponse = Invoke-WebRequest -Uri "$baseUrl/files/upload" `
        -Method Post `
        -Headers $headers `
        -Body $bodyLines
    
    $responseContent = $fileUploadResponse.Content | ConvertFrom-Json
    Write-Host "File upload response: $($responseContent | ConvertTo-Json)"
}
catch {
    Write-Host "Upload failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody"
    }
}

# Test file upload from URL
Write-Host "`nTesting file upload from URL..."
$urlUploadBody = @{
    url = "https://raw.githubusercontent.com/microsoft/PowerShell/master/README.md"
    name = "readme.md"
} | ConvertTo-Json

try {
    $urlUploadResponse = Invoke-RestMethod -Uri "$baseUrl/files/upload-url" `
        -Method POST `
        -Body $urlUploadBody `
        -ContentType "application/json"
    Write-Host "URL upload response: $($urlUploadResponse | ConvertTo-Json)"
}
catch {
    Write-Host "URL upload failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody"
    }
}

# Test list files
Write-Host "`nTesting list files..."
try {
    $listFilesResponse = Invoke-RestMethod -Uri "$baseUrl/files" -Method GET
    Write-Host "List files response: $($listFilesResponse | ConvertTo-Json)"
}
catch {
    Write-Host "List files failed: $($_.Exception.Message)"
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody"
    }
}

# Test download file (if we have a file ID from previous operations)
if ($fileUploadResponse.Content.id) {
    Write-Host "`nTesting file download..."
    try {
        $downloadPath = Join-Path $PSScriptRoot "downloaded_$($fileUploadResponse.Content.id).txt"
        $downloadResponse = Invoke-RestMethod -Uri "$baseUrl/files/$($fileUploadResponse.Content.id)/download" `
            -Method GET `
            -OutFile $downloadPath
        Write-Host "File downloaded successfully to: $downloadPath"
    }
    catch {
        Write-Host "Download failed: $($_.Exception.Message)"
    }
}

# Test hide file
if ($fileUploadResponse.Content.id) {
    Write-Host "`nTesting hide file..."
    try {
        $hideResponse = Invoke-RestMethod -Uri "$baseUrl/files/$($fileUploadResponse.Content.id)/hide" -Method PATCH
        Write-Host "Hide file response: $($hideResponse | ConvertTo-Json)"
    }
    catch {
        Write-Host "Hide file failed: $($_.Exception.Message)"
    }
}

# Test delete file
if ($fileUploadResponse.Content.id) {
    Write-Host "`nTesting delete file..."
    try {
        $deleteResponse = Invoke-RestMethod -Uri "$baseUrl/files/$($fileUploadResponse.Content.id)" -Method DELETE
        Write-Host "Delete file response: $($deleteResponse | ConvertTo-Json)"
    }
    catch {
        Write-Host "Delete file failed: $($_.Exception.Message)"
    }
}

# Clean up test file
if (Test-Path $filePath) {
    Remove-Item -Path $filePath -Force
    Write-Host "`nTest file cleaned up"
} 