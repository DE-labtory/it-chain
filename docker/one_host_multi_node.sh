abspath=$PWD

if [ $# -ne 1 ]; then
	echo please input node number to init or --help
	echo "ex) ./create_network.sh 4"
	exit 1
fi

if [ $1 == "--help" ]; then
	echo use this shell script with node num to create
	echo "ex ) ./create_network.sh 4"
	echo "-> mean create 4 node and connect each to create network. it will make blocak chain network with 4 node."
	echo
	echo each api-gateway  port has 40001 ~ 4000X, X = node num. so node1 = 40001, node2 = 40002 ...
	echo each grpc-gateway port has 50001 ~ 5000X, X = node num. so node1 = 50001, node2 = 50002 ...
	exit 0
fi

echo "wait for pull it-chain docker image"
#
# docker pull itchain-solo latest version
#
docker pull teamit/itchain:latest

echo "wait for remove legacy network and volume"
docker network rm it-chain-default-network  > /dev/null 2>&1
docker volume rm it-chain-default-volume > /dev/null 2>&1
sleep 5

echo "wait for create node container..."
docker network create --subnet=172.88.0.0/16 it-chain-default-network > /dev/null 2>&1
docker volume create it-chain-default-volume > /dev/null 2>&1
sleep 2

docker run -d -p 40001:4000 -p 50001:5000 --net it-chain-default-network -v it-chain-default-volume:/go/src/github.com/DE-labtory/it-chain/.tmp \
-v /var/run/docker.sock:/var/run/docker.sock --ip 172.88.1.2 \
-v $abspath/docker_multi_config_leader.yaml:/go/src/github.com/DE-labtory/it-chain/conf/config.yaml teamit/itchain:latest > /dev/null 2>&1

if [ $1 -eq 1 ]; then
	sleep 10
	echo finish create node!
	echo node1 - api-gateway : 40001, grpc-gateway : 50001
	exit
fi

for i in $(eval echo {2..$1})
do
	cp $abspath/docker_multi_config_member.yaml docker_multi_config_member_$i.yaml
	sed -i -e "s/a.a.a.a/172.88.$i.2/" docker_multi_config_member_$i.yaml
	sed -i -e "s/b.b.b.b/172.88.$i.0/" docker_multi_config_member_$i.yaml
	docker run -d -p 0.0.0.0:4000$i:4000 -p 0.0.0.0:5000$i:5000 --net it-chain-default-network \
	-v it-chain-default-volume:/go/src/github.com/DE-labtory/it-chain/.tmp \
	-v /var/run/docker.sock:/var/run/docker.sock --ip 172.88.$i.2 \
	-v $abspath/docker_multi_config_member_$i.yaml:/go/src/github.com/DE-labtory/it-chain/conf/config.yaml teamit/itchain:latest \
	> /dev/null 2>&1
done

#
# wait for 20 sec for waiting start node.
# if something err with using this script,
# try to increase this time first
#

echo "wait for start all node... "
echo "it will take about 20 sec ..."
for i in $(eval echo {1..$1})
do
	echo "try to check node$i..."
	status_code=$(curl --write-out %{http_code} --silent --output /dev/null http://127.0.0.1:4000$i/peers)
	while [ "$status_code" != 200 ]
	do
		sleep 1
		status_code=$(curl --write-out %{http_code} --silent --output /dev/null http://127.0.0.1:4000$i/peers)
		sleep 5
	done
done

echo "start to connect each node..."

for i in $(eval echo {2..$1})
do
	curl --header "Content-Type: application/json" \
	  --request POST \
	  --data '{"type":"join","address":"172.88.1.2:5000"}' \
	  http://127.0.0.1:4000$i/peers
	sleep 5
done

echo finish to connect each node...
sleep 2

echo
echo

echo node port info------------
for i in $(eval echo {1..$1})
do
	echo node$i - api-gateway : 4000$i, grpc-gateway : 5000$i
done

rm $abspath/docker_multi_config_member_*


