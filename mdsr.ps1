#!/usr/bin/env pwsh

#########################################################################
# Restart wrapper for markdown-serve (Windows/PowerShell): kills any
# already-running instance and starts a new one with the given args.
# Part of https://github.com/muquit/markdown-serve
# Jun-26-2026 
#########################################################################

param(
    [Parameter(ValueFromRemainingArguments = $true)]
    [string[]]$Args
)

$BinaryName = if ($env:MARKDOWN_SERVE_BIN) { $env:MARKDOWN_SERVE_BIN } else { "markdown-serve" }

$cmd = Get-Command $BinaryName -ErrorAction SilentlyContinue
if (-not $cmd) {
    Write-Error "'$BinaryName' not found in PATH (set `$env:MARKDOWN_SERVE_BIN to override)"
    exit 1
}
$resolved = $cmd.Source
$procName = [System.IO.Path]::GetFileNameWithoutExtension($resolved)

$running = Get-Process -Name $procName -ErrorAction SilentlyContinue

if ($running) {
    foreach ($p in $running) {
        $info = Get-CimInstance Win32_Process -Filter "ProcessId=$($p.Id)" -ErrorAction SilentlyContinue
        $cmdLine = if ($info) { $info.CommandLine } else { $p.ProcessName }
        Write-Host "Killing running markdown-serve (pid $($p.Id)): $cmdLine" -ForegroundColor Red
    }
    $running | Stop-Process -Force
    $running | Wait-Process -Timeout 5 -ErrorAction SilentlyContinue
}

Write-Host "Starting: $resolved $Args"
& $resolved @Args
exit $LASTEXITCODE
