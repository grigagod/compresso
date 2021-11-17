#!bin/bash

# update
sudo yum update

# install docker
sudo yum install docker

# start docker
sudo service docker start

# create docker group
sudo groupadd docker

# add user to docker group to allow docker commands without sudo
sudo usermod -a -G docker $USER

# add docker to autostart and reboot to see if it's working
sudo chkconfig docker on
sudo reboot
