@echo off
REM Build the Claude Usage Tracker widget for Windows
REM Requires: Go 1.22+, CGo enabled (for WebView2 native mode)

echo Building Claude Usage Tracker...

REM Build with optimizations and strip debug info for smaller binary
set CGO_ENABLED=0
go build -ldflags="-s -w" -o claude-usage-tracker.exe ./cmd/widget

if %errorlevel% neq 0 (
    echo.
    echo Build failed! Trying with CGO for native widget support...
    set CGO_ENABLED=1
    go build -ldflags="-s -w -H windowsgui" -o claude-usage-tracker.exe ./cmd/widget
)

if %errorlevel% equ 0 (
    echo.
    echo Build successful: claude-usage-tracker.exe
    echo.
    echo Run with:
    echo   claude-usage-tracker.exe              (native widget, Windows only)
    echo   claude-usage-tracker.exe --browser    (opens in browser)
    echo   claude-usage-tracker.exe --port 8080  (custom port for browser mode)
) else (
    echo Build failed.
    exit /b 1
)
