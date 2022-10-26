# Docker data collector ![CI](https://github.com/crayonwow/docker_data_collector/actions/workflows/develop.yml/badge.svg)

A simple project for collection data from running docker containers on host.

### Origin
I have a digital ocean host that serves vpn's for me and my friends. My friends pay me every month. But sometimes they 
and me also forget about payments so me decided to  build this project as a helpful tool for notify myself that they need to pay :)

### How it works
It uses docker SDK to collect data from containers. To make it run in container we mount `/var/run/docker.sock` from a 
host machine to docker container that runs this application. This make able to fetch data from host docker form inside docker container.
I'm fully aware that not really safe and creates a lot of problems (read 
[article](http://jpetazzo.github.io/2015/09/03/do-not-use-docker-in-docker-for-ci/) for details).
But all what I do is fetch data about containers.

### TODO
* make it as a package
* cover with test