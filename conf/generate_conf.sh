if [ -f main/GeneratorBuildFile ]; then
	echo "GeneratorBuildFile already exist"
	echo "Trying to remove old generator file"
	rm -rf main/GeneratorBuildFile
fi

if [ ! -f main/conf_file_generator.go ]; then
	echo "conf_file_generator.go is not exist"
	exit
fi
echo "Build generator file..."
go build -o main/GeneratorBuildFile main/conf_file_generator.go
echo "trying to excute generator..."
./main/GeneratorBuildFile
echo "trying to remove excuted generator..."
rm -rf main/GeneratorBuildFile
echo "conf file generate finish!!!"


