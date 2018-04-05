# Pubtail

Pubtail is tool for tailing Google Cloud Pubsub Message.

# Installation
```
go get github.com/sonnythehottest/pubtail
```

# Example
```
smanurung@sonny-macbook pubtail (master)*$ ./pubtail -project=<project-id> -topics=playSonny1,playSonny2
INFO[0002] subscription already exists
INFO[0002] listening to topic playSonny1...
INFO[0002] subscription already exists
INFO[0002] listening to topic playSonny2...
INFO[0015] hello world                                   foo=bar topic=playSonny1
INFO[0058] halo broh                                     topic=playSonny2
```

### TODO
- Read config from config files
- Set different colors for each topic