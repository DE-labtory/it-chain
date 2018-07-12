@echo off
SET BAT_PATH=%~dp0
if exist %BAT_PATH%/main/GeneratorBuildFile.exe (
	echo GeneratorBuildFile already exist
	echo Trying to remove old generator file
	del %BAT_PATH%\main\GeneratorBuildFile.exe
)

if not exist %BAT_PATH%/main/conf_file_generator.go (
	echo conf_file_generator.go is not exist
	exit 1
)

SET CONFIG_NAME=config
IF NOT * == %1* (
    SET CONFIG_NAME=%1
)
echo Build generator files...
go build -o main/GeneratorBuildFile.exe main/conf_file_generator.go
echo trying to excute generator...
%BAT_PATH%/main/GeneratorBuildFile.exe -name %CONFIG_NAME%
echo trying to remove excuted generator...
del "%BAT_PATH%\main\GeneratorBuildFile.exe"
echo %CONFIG_NAME%.yaml generate finish!!!
pause