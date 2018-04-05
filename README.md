# Pubtail

Pubtail is tool for tailing Google Cloud Pubsub Message.

# Installation
```
go get github.com/sonnythehottest/pubtail
```

# Example
### Tail text message
```
smanurung@sonny-macbook pubtail (master)*$ ./pubtail -project=<project-id> -topics=playSonny1,playSonny2
INFO[0002] subscription already exists
INFO[0002] listening to topic playSonny1...
INFO[0002] subscription already exists
INFO[0002] listening to topic playSonny2...
INFO[0015] hello world                                   foo=bar topic=playSonny1
INFO[0058] halo broh                                     topic=playSonny2
```

### Tail avro message
```
smanurung@sonny-macbook pubtail (master)*$ ./pubtail -project=<projectid> -topics=<topic> -format=avro
INFO[0000] connecting to projectid <projectid>, using format avro
INFO[0006] listening to topic <topic>...
INFO[0018] map[adShown:[""] network_sessionId:2bc059e8d154188028bdb525a1b41869 filter_preorder:0 ip_address:x.x.x.x network_ipAddress:x.x.x.x user_email:<user_email> hash_id:YWRTaG93bj1bXSxmaWx0ZXI9bWFwW3BhZ2U6MSBkZXBhcnRtZW50SWQ6W10gaXNHbTowIGtyZWFzaUxva2FsOiB3aG9sZXNhbGU6MCBmcmVlcmV0dXJuOjAgY29uZGl0aW9uOjAgdmFyaWFudDpbXSBwcmljZU1pbjowIHByaWNlTWF4OjAgaXNPZmZpY2lhbDowIHByZW9yZGVyOjAgbnVtYmVyQ container=<container> event=<event> topic=<topic-name>
```

### TODO
- Read config from config files
- Set different colors for each topic
- Set format per topic basis