<#
Builds the files for an FMOS release into build\release\v<version>\:

  fmos-x64.exe     Windows on Intel or AMD
  fmos-arm64.exe   Windows on ARM
  fmos-macos       macOS on Apple Silicon
  SHA256SUMS.txt   the fingerprint (SHA-256) of each file above
  notes.md         release notes for you to fill in (not uploaded)

This is the Windows version of scripts/release.sh. It does the same steps and
the same checks. Every program file is built from committed code, with
-trimpath (so no folder paths from this computer end up inside) and with the
release version. It never uploads anything itself.

Usage, from the FMOS folder in PowerShell:

  powershell -ExecutionPolicy Bypass -File scripts\release.ps1 0.1.0

("-ExecutionPolicy Bypass" lets this one script run without changing any
Windows settings.)
#>

param([string]$Version)

$ErrorActionPreference = 'Stop'

$builds = @(
    @{ Goos = 'windows'; Goarch = 'amd64'; File = 'fmos-x64.exe' },
    @{ Goos = 'windows'; Goarch = 'arm64'; File = 'fmos-arm64.exe' },
    @{ Goos = 'darwin'; Goarch = 'arm64'; File = 'fmos-macos' }
)

function Invoke-Release {
    if (git status --porcelain) {
        throw 'There are uncommitted changes. Commit them first, so the release matches the code on GitHub.'
    }
    git rev-parse -q --verify "refs/tags/v$Version" | Out-Null
    if ($LASTEXITCODE -eq 0) {
        throw "The tag v$Version already exists. Pick a new version number."
    }

    $out = "build\release\v$Version"
    if (Test-Path $out) {
        throw "$out already exists. Delete that folder first if you want to build this version again."
    }
    New-Item -ItemType Directory -Path $out | Out-Null
    # .NET file functions need a full path, because they do not follow Set-Location.
    $outFull = (Resolve-Path $out).Path

    $env:CGO_ENABLED = '0'
    foreach ($build in $builds) {
        Write-Host "Building $($build.File)..."
        $env:GOOS = $build.Goos
        $env:GOARCH = $build.Goarch
        go build -trimpath -ldflags "-X main.version=$Version" -o "$out\$($build.File)" ./cmd/fmos
        if ($LASTEXITCODE -ne 0) {
            throw "Building $($build.File) failed. The error is shown above."
        }
    }

    # SHA256SUMS.txt: one "<hash>  <file name>" line per file, the same format
    # that "shasum -a 256" prints on a Mac. Written without a byte order mark.
    $lines = foreach ($build in $builds) {
        $hash = (Get-FileHash -Algorithm SHA256 -Path "$outFull\$($build.File)").Hash.ToLower()
        "$hash  $($build.File)"
    }
    $utf8 = New-Object System.Text.UTF8Encoding $false
    [System.IO.File]::WriteAllText("$outFull\SHA256SUMS.txt", (($lines -join "`n") + "`n"), $utf8)

    # Check 1: no folder paths from this computer inside the files. Go writes
    # paths with forward slashes, so look for both spellings.
    $profilePaths = @($env:USERPROFILE, ($env:USERPROFILE -replace '\\', '/')) | Where-Object { $_ }
    $latin1 = [System.Text.Encoding]::GetEncoding(28591)
    foreach ($build in $builds) {
        $content = $latin1.GetString([System.IO.File]::ReadAllBytes("$outFull\$($build.File)"))
        foreach ($path in $profilePaths) {
            if ($content.IndexOf($path, [System.StringComparison]::OrdinalIgnoreCase) -ge 0) {
                throw "Stopped: $($build.File) still contains paths from this computer ($path). It was not built with -trimpath."
            }
        }
    }

    # Check 2: the fingerprints in SHA256SUMS.txt match the files.
    foreach ($line in [System.IO.File]::ReadAllLines("$outFull\SHA256SUMS.txt")) {
        $expected, $file = $line -split '\s+', 2
        $actual = (Get-FileHash -Algorithm SHA256 -Path "$outFull\$file").Hash.ToLower()
        if ($actual -ne $expected) {
            throw "Stopped: the fingerprint of $file does not match SHA256SUMS.txt."
        }
    }

    # Check 3: the build for this computer reports the release version.
    $ownBuild = @{ AMD64 = 'fmos-x64.exe'; ARM64 = 'fmos-arm64.exe' }[$env:PROCESSOR_ARCHITECTURE]
    if ($ownBuild) {
        $reported = (& "$outFull\$ownBuild" --version | Out-String).Trim()
        if ($reported -ne "FMOS $Version") {
            throw "Stopped: $ownBuild reports `"$reported`" instead of `"FMOS $Version`"."
        }
    }

    $notes = "## What's new in FMOS $Version`n`n" +
        "- (Write a few short, plain-English lines. FMOS shows them on its Updates screen.)`n"
    [System.IO.File]::WriteAllText("$outFull\notes.md", $notes, $utf8)

    $commit = git rev-parse HEAD
    Write-Host ''
    Write-Host "Release files for FMOS $Version are ready in $($out):"
    foreach ($line in $lines) {
        $hash, $file = $line -split '\s+', 2
        Write-Host ('  {0,-15}  fingerprint starts with {1}' -f $file, $hash.Substring(0, 8))
    }
    Write-Host ''
    Write-Host "Checked: built with -trimpath, fingerprints match, version is $Version."
    Write-Host ''
    Write-Host 'Next:'
    Write-Host "  1. Write the release notes in $out\notes.md"
    Write-Host '  2. Make sure this commit is on GitHub: git push'
    Write-Host '  3. Publish, either with the GitHub CLI:'
    Write-Host ('       gh release create v{0} --target {1} --title "FMOS {0}" `' -f $Version, $commit)
    Write-Host ('         --notes-file {0}\notes.md `' -f $out)
    Write-Host ('         {0}\fmos-x64.exe {0}\fmos-arm64.exe {0}\fmos-macos {0}\SHA256SUMS.txt' -f $out)
    Write-Host "     or on github.com: Releases > Draft a new release > tag v$Version,"
    Write-Host '     paste notes.md, and attach the 3 programs and SHA256SUMS.txt.'
}

if ($Version -notmatch '^\d+\.\d+\.\d+$') {
    Write-Host 'Usage: scripts\release.ps1 <version>   for example: scripts\release.ps1 0.1.0' -ForegroundColor Red
    exit 1
}

$failed = $false
# The builds change these settings; put back whatever was there before.
$savedSettings = @{ CGO_ENABLED = $env:CGO_ENABLED; GOOS = $env:GOOS; GOARCH = $env:GOARCH }
Push-Location (Join-Path $PSScriptRoot '..')
try {
    Invoke-Release
} catch {
    Write-Host $_.Exception.Message -ForegroundColor Red
    $failed = $true
} finally {
    foreach ($name in $savedSettings.Keys) {
        [Environment]::SetEnvironmentVariable($name, $savedSettings[$name])
    }
    Pop-Location
}
if ($failed) {
    exit 1
}
