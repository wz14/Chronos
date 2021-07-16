# Code in AAC

## Install docker in ubuntu20

```shell script
sudo apt-get update
sudo apt-get install docker.io
# add docker group to avoid use `sudo` before docker
sudo groupadd docker
sudo usermod -aG docker $USER
# restart shell or terminal
```

## Remote docker access

```shell script
sudo service docker stop
```

Create `/etc/docker/daemon.json` and write.
```json
{
  "hosts" : ["unix:///var/run/docker.sock", "tcp://0.0.0.0:2375"]
}
```

Create `/etc/systemd/system/docker.service.d/override.conf` and write
```text
##Add this to the file for the docker daemon to use different ExecStart parameters (more things can be added here)
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd
```

```shell script
sudo systemctl daemon-reload
sudo systemctl restart docker.service
```

## Install docker compose 

```shell script
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

## Download source or Get source

note: Use jetbrain need to open module in go setting.

## Open new Branch and Merge

