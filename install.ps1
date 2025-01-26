# Installation directory
$installDir = "$env:ProgramFiles\dockerizer"
$binary = "dockerizer-windows-amd64.exe"

# Create installation directory if it doesn't exist
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Download latest release
Write-Host "Downloading latest release..."
$latestRelease = (Invoke-RestMethod "https://api.github.com/repos/ravanbabayev/dockerizer-cli/releases/latest").tag_name
$downloadUrl = "https://github.com/ravanbabayev/dockerizer-cli/releases/download/$latestRelease/$binary.zip"

# Download and extract
$zipPath = "$env:TEMP\dockerizer.zip"
Invoke-WebRequest -Uri $downloadUrl -OutFile $zipPath
Expand-Archive -Path $zipPath -DestinationPath $installDir -Force
Remove-Item $zipPath

# Rename binary
Move-Item -Path "$installDir\$binary" -Destination "$installDir\dockerizer.exe" -Force

# Add to PATH if not already there
$currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($currentPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$currentPath;$installDir",
        "Machine"
    )
}

Write-Host "Installation complete! Please restart your terminal to use 'dockerizer' command." 