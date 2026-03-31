@echo off
chcp 65001 >nul
echo Поиск pdflatex...
echo.

REM Попробуем найти pdflatex в стандартных местах
set PDFLATEX_PATH=

if exist "C:\Program Files\MiKTeX\miktex\bin\x64\pdflatex.exe" (
    set "PDFLATEX_PATH=C:\Program Files\MiKTeX\miktex\bin\x64\pdflatex.exe"
    set "BIBTEX_PATH=C:\Program Files\MiKTeX\miktex\bin\x64\bibtex.exe"
    goto :found
)

if exist "%LOCALAPPDATA%\Programs\MiKTeX\miktex\bin\x64\pdflatex.exe" (
    set "PDFLATEX_PATH=%LOCALAPPDATA%\Programs\MiKTeX\miktex\bin\x64\pdflatex.exe"
    goto :found
)

if exist "C:\Program Files (x86)\MiKTeX\miktex\bin\pdflatex.exe" (
    set "PDFLATEX_PATH=C:\Program Files (x86)\MiKTeX\miktex\bin\pdflatex.exe"
    goto :found
)

echo pdflatex не найден в стандартных местах!
echo.
echo Пожалуйста, укажите полный путь к pdflatex.exe
echo Например: C:\Program Files\MiKTeX\miktex\bin\x64\pdflatex.exe
set /p PDFLATEX_PATH="Введите путь к pdflatex.exe: "
if "%PDFLATEX_PATH%"=="" (
    echo Ошибка: путь не указан
    pause
    exit /b 1
)

:found
echo Найден pdflatex: %PDFLATEX_PATH%
echo.
echo ========================================
echo Компиляция LaTeX документа
echo ========================================
echo.

cd /d "%~dp0"

echo [1/4] Первый проход pdflatex...
"%PDFLATEX_PATH%" -interaction=nonstopmode pract-example.tex
REM Продолжаем даже при ошибках (могут быть предупреждения о шрифтах)

echo.
echo [2/4] Запуск bibtex для обработки библиографии...
if exist "%BIBTEX_PATH%" (
    "%BIBTEX_PATH%" pract-example
) else (
    echo ПРЕДУПРЕЖДЕНИЕ: bibtex не найден, пропускаем шаг
)

echo.
echo [3/4] Второй проход pdflatex...
"%PDFLATEX_PATH%" -interaction=nonstopmode pract-example.tex
REM Продолжаем даже при ошибках

echo.
echo [4/4] Третий проход pdflatex (финальный)...
"%PDFLATEX_PATH%" -interaction=nonstopmode pract-example.tex
REM Продолжаем даже при ошибках

echo.
echo ========================================
echo Компиляция завершена успешно!
echo PDF файл: pract-example.pdf
echo ========================================
echo.
pause
