# beats-forwarder
A simple forwarder to send beats everywhere through udp, tcp, syslog and third-party




## How-to use it

### Installation

### Installation
Notice that we've released a docker version: get the docker [beats-forwarder]() image
Download the last version of the beats-forwarder ([all versions available]()):

```sh
# linux 64
mkdir -p beats-forwarder/etc && cd beats-forwarder
curl -OL https://github.com/logmatic/beats-forwarder/releases/download/v0.1-alpha/beats-forwarder
cd etc && curl -OL https://raw.githubusercontent.com/logmatic/beats-forwarder/dev/etc/config.yml
cd ..
chmod +x beats-forwarder
```

### Configuration
Beats-forwarder use a Yaml file as configuration. By default, beats-forwarder listens on
all interfaces of the server and use 5044 as port. Beats are forwarder to the local syslog.
Beats-forwarder allows you to ship them to:
* A syslog server (local or remote)
* A tcp/udp endpoint (secure or not)
* Logmatic.io

The default configuration can be found here: [config.yaml](logmatic/beats-forwarder/blob/dev/etc/config.yml)
Create a new file, `beats-fwdr.yaml` and set these attributes at least (it's a
recommendation not a mandatory):

```yaml
####
#### beats-forwarder default configuration
####

input:
  # The port to listen on
  port: 5044

output:

  # The wanted output (syslog|udp_tcp|logmatic), by default syslog
  type: syslog

  # Syslog specific settings
  syslog:
    # Tag or application reported for each log
    tag: beats-fowarder-demo

  # Logmatic specific settings
  logmatic:

    # The Logmatic API Key for authentification
    key: "<YOUR_API_KEY>"

```

If you want to send beats directly to Logmatic.io, just set `output.type` to `logmatic` and
add copy/paste your Logmatic APY Key to `output.type.logmatic.key`.

### Run
Here we are!
```
./beats-forwarder -c beats-fwdr.yaml
```

Now, just configure your already existing beats to send them to the forwarder.

### Beats configuration
All you need to do is to add and configure the logstash output for each beat.s
Edit `beatsname-config.yml` and add at the end the following code:

```yaml
output:

  logstash:
    # Set the beats-forwarder address
    hosts: [ "localhost:5044"]

  # The rest of the output configuration goes here ...

```

Restart the beat, and check the incomming beats. If you have followed this tutorial,
beats are sent to the local syslog.

```sh
# this can be different depending on your OS
journalctl -f
```
And the magic goes on:
```
...
Oct 06 15:33:56 jarvis beats-by-gpolaert[6417]: {"@metadata":{"beat":"topbeat","type":"filesystem"},"@timestamp":"2016-10-06T13:33:55.108Z","beat":{"hostname":"jarvis","name":"jarvis"},"count":1,"fs":{"avail":0,"device_name":"cgroup","files":0,"free":0,"free_files":0,"mount_point":"/sys/fs/cgroup/freezer","total":0,"used":0,"used_p":0},"type":"filesystem"}
Oct 06 15:33:56 jarvis beats-by-gpolaert[6417]: {"@metadata":{"beat":"topbeat","type":"filesystem"},"@timestamp":"2016-10-06T13:33:55.108Z","beat":{"hostname":"jarvis","name":"jarvis"},"count":1,"fs":{"avail":0,"device_name":"cgroup","files":0,"free":0,"free_files":0,"mount_point":"/sys/fs/cgroup/pids","total":0,"used":0,"used_p":0},"type":"filesystem"}
Oct 06 15:33:56 jarvis beats-by-gpolaert[6417]: {"@metadata":{"beat":"topbeat","type":"filesystem"},"@timestamp":"2016-10-06T13:33:55.108Z","beat":{"hostname":"jarvis","name":"jarvis"},"count":1,"fs":{"avail":0,"device_name":"mqueue","files":0,"free":0,"free_files":0,"mount_point":"/dev/mqueue","total":0,"used":0,"used_p":0},"type":"filesystem"}
Oct 06 15:33:56 jarvis beats-by-gpolaert[6417]: {"@metadata":{"beat":"topbeat","type":"filesystem"},"@timestamp":"2016-10-06T13:33:55.108Z","beat":{"hostname":"jarvis","name":"jarvis"},"count":1,"fs":{"avail":0,"device_name":"configfs","files":0,"free":0,"free_files":0,"mount_point":"/sys/kernel/config","total":0,"used":0,"used_p":0},"type":"filesystem"}
```

