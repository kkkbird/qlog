# filebeat-qlog

## format setting

* udp: JSON with prettycaller set to truncated

```yaml
  udp:
    enabled: true
    level: debug
    host: 192.168.1.112:6060
    formatter:
      name: json
      opts:
        prettycaller: truncated
```

* file: CLASSIC with prettycaller set to truncated / NOT tested

```yaml
  file:
    enabled: true
    path: ./logs/
    name: message.log
    level: trace
    formatter:
      name: classic
      opts:
        prettycaller: truncated
```
