abspath=$PWD
if [ $1 == "--help" ]; then
	echo "use this shell script with output port num to assign node"
	echo "./single_node.sh [api-gateway port] [grpc-gateway port]"
	echo "example usage ) ./single_node.sh 40000 50000"
exit 0
fi

if [ $# -ne 2 ]; then
	echo "please input two output port number to create node or --help"
	exit 1
fi


echo "wait for pull it-chain docker image"
#
# docker pull itchain-solo latest version
#
#docker pull teamit/itchain:latest

echo "wait for remove legacy network and volume"
docker network rm it-chain-default-network # > /dev/null 2>&1
docker volume rm it-chain-default-volume #> /dev/null 2>&1
sleep 5

echo "wait for create node container..."
docker network create --subnet=172.88.1.0/24 it-chain-default-network #> /dev/null 2>&1
docker volume create it-chain-default-volume
sleep 2

docker run -d -p $1:4000 -p $2:5000 --net it-chain-default-network -v it-chain-default-volume:/go/src/github.com/it-chain/engine/.tmp \
-v /var/run/docker.sock:/var/run/docker.sock --ip 172.88.1.2 \
-v $abspath/docker_solo_config.yaml:/go/src/github.com/it-chain/engine/conf/config.yaml teamit/itchain:latest

echo "finish create node!"
echo "node1 - api-gateway : $1, grpc-gateway : $2"

