$dirs = @(
    "cmd\desktop", "cmd\cli", "cmd\daemon",
    "internal\engine", "internal\models", "internal\server", "internal\client", 
    "internal\config", "internal\daemon", "internal\embedder", "internal\utils",
    "pkg\api", "pkg\protocol",
    "ui\backend", "ui\frontend\src\components", "ui\frontend\src\services", 
    "ui\frontend\src\stores", "ui\frontend\src\styles", "ui\frontend\public",
    "native\bitnet-engine\windows", "native\bitnet-engine\linux", "native\bitnet-engine\darwin",
    "native\bindings",
    "manifests\official", "manifests\community",
    "installer\windows\wix", "installer\windows\inno", "installer\windows\resources", "installer\windows\scripts",
    "installer\linux\debian", "installer\linux\rpm", "installer\linux\appimage", "installer\linux\snap",
    "installer\macos\dmg", "installer\macos\pkg\scripts",
    "scripts\build", "scripts\release", "scripts\dev",
    "tests\unit", "tests\integration", "tests\e2e", "tests\fixtures\models", "tests\fixtures\configs",
    "docs\architecture\diagrams", "docs\architecture\decisions", "docs\api", "docs\user", "docs\developer",
    "configs\examples",
    "assets\icons", "assets\images", "assets\branding",
    "build\bin", "build\dist", "build\temp",
    ".github\workflows", ".github\ISSUE_TEMPLATE"
)

foreach ($dir in $dirs) {
    New-Item -ItemType Directory -Force -Path $dir | Out-Null
    Write-Host "Created: $dir"
}

# Create empty .gitkeep files to ensure git tracks empty folders
Get-ChildItem -Recurse -Directory | ForEach-Object {
    if ((Get-ChildItem $_.FullName).Count -eq 0) {
        New-Item -ItemType File -Path "$($_.FullName)\.gitkeep" | Out-Null
    }
}