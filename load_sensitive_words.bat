@echo off
chcp 65001 >nul
echo 🚀 开始加载敏感词到数据库...

REM 检查Python是否安装
python --version >nul 2>&1
if errorlevel 1 (
    echo ❌ Python未安装或未添加到PATH环境变量
    pause
    exit /b 1
)

REM 检查是否安装了依赖
echo 📦 检查Python依赖...
pip show mysql-connector-python >nul 2>&1
if errorlevel 1 (
    echo 📥 安装Python依赖...
    pip install -r requirements.txt
    if errorlevel 1 (
        echo ❌ 依赖安装失败
        pause
        exit /b 1
    )
)

REM 运行Python脚本
echo 🐍 运行敏感词加载脚本...
python load_sensitive_words.py

echo.
echo 按任意键退出...
pause >nul 