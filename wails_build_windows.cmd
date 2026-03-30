@echo off
setlocal EnableExtensions EnableDelayedExpansion

for %%I in ("%~dp0.") do set "PROJECT_ROOT=%%~fI"
set "MSYS64=%MSYS64_ROOT%"
if not defined MSYS64 set "MSYS64=C:\DevDisk\DevTools\msys64"
set "MINGW_BIN=%MSYS64%\mingw64\bin"
set "FREERDP_INSTALL=%PROJECT_ROOT%\install"
set "BUILD_BIN=%PROJECT_ROOT%\build\bin"

where wails >nul 2>nul
if errorlevel 1 (
    echo ERROR: wails not found in PATH
    exit /b 1
)
where go >nul 2>nul
if errorlevel 1 (
    echo ERROR: go not found in PATH
    exit /b 1
)
where node >nul 2>nul
if errorlevel 1 (
    echo ERROR: node not found in PATH
    exit /b 1
)

if not exist "%FREERDP_INSTALL%\bin\libfreerdp3.dll" (
    echo ERROR: missing %FREERDP_INSTALL%\bin\libfreerdp3.dll
    echo Please run: build_windows.cmd
    exit /b 1
)

set "PATH=%MINGW_BIN%;%FREERDP_INSTALL%\bin;%PATH%"
set "OPENSSL_CONF=%PROJECT_ROOT%\openssl.cnf"
set "OPENSSL_MODULES=%FREERDP_INSTALL%\bin\ossl-modules"

cd /d "%PROJECT_ROOT%"
echo Building Wails package ^(Windows^)...
wails build -clean %*
if errorlevel 1 (
    echo ERROR: wails build failed
    exit /b 1
)

if not exist "%BUILD_BIN%" mkdir "%BUILD_BIN%"
for %%F in ("%FREERDP_INSTALL%\bin\*.dll") do copy /Y "%%~fF" "%BUILD_BIN%\" >nul
if exist "%FREERDP_INSTALL%\bin\ossl-modules" (
    xcopy /E /I /Y "%FREERDP_INSTALL%\bin\ossl-modules" "%BUILD_BIN%\ossl-modules" >nul
)

echo Done: %BUILD_BIN%
