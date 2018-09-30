#
# input only one paramter for node number. 
#
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

echo wait for pull it-chain docker image
#
# docker pull itchain latest version
#
docker pull teamit/itchain:latest

#
# create node in docker container.
# each node has port to api-gateway : 4000X
#		       grpc-gateway : 5000X
# X = node number
#

echo wait for create node container...

docker network create --subnet=172.18.0.0/16 it-chain > /dev/null 2>&1
sleep 2

docker run -d -p 40001:4000 -p 50001:5000 -v $abspath/config_docker_leader.yaml:/go/src/github.com/it-chain/engine/conf/config_docker.yaml \
		--privileged --net it-chain --ip 172.18.1.0 teamit/itchain:latest

#
# if node num is 1, not try to connect to other node. just finish
#

if [ $1 -eq 1 ]; then
	sleep 20
	echo finish create node!
	echo node1 - api-gateway : 40001, grpc-gateway : 50001
	exit
fi



for i in $(eval echo {2..$1})
do
	cp $abspath/config_docker_member.yaml config_docker_member_$i.yaml
	sed -i -e "s/0.0.0.0/172.18.$i.0/" config_docker_member_$i.yaml
	docker run -d -p 4000$i:4000 -p 5000$i:5000 \
		   -v $abspath/config_docker_member_$i.yaml:/go/src/github.com/it-chain/engine/conf/config_docker.yaml \
		   --privileged --net it-chain --ip 172.18.$i.0 teamit/itchain:latest
done

#
# wait for 20 sec for waiting start node.
# if something err with using this script,
# try to increase this time first
#

echo wait for start all node... 
echo it will take about 60 sec ...

for i in $(eval echo {1..$1})
do
	echo try to check node$i...
	status_code=$(curl --write-out %{http_code} --silent --output /dev/null http://127.0.0.1:4000$i/peers)
	while [ "$status_code" != 200 ]
	do
		sleep 1
		status_code=$(curl --write-out %{http_code} --silent --output /dev/null http://127.0.0.1:4000$i/peers)
		sleep 5
	done
done

echo start to connect each node...

for i in $(eval echo {2..$1})
do
	curl --header "Content-Type: application/json" \
	  --request POST \
	  --data '{"type":"join","address":"172.18.1.0:5000"}' \
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

rm $abspath/config_docker_member_*
