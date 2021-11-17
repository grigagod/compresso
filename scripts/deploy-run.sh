!bin/bash

# start watchtower
docker run --privileged -d --name wathtower \
    -v /var/run/docker.sock:/var/run/docker.sock containrrr/watchtower -i 30 --cleanup

echo -n "Please enter a service name: "
read svcname

# start service
docker run --name $svcname --env-file=.env \
    --log-driver=awslogs --log-opt awslogs-region=eu-central-1 --log-opt awslogs-group=compresso-$svcname \
    -p 80:80 -d compressorepo/$svcname
