# fluent-bit-output-rsyslog
An output plugin of FluentBit to send log via rsyslog

## Usage
```
[SERVICE]
    Flush        1
    Daemon       Off
    Log_Level    info
    Log_File     /fluent-bit/log/fluent-bit.log
    Parsers_File parsers.conf
    Parsers_File parsers_java.conf
    Plugins_File plugins.conf

[INPUT]
    Name Forward
    Port 24224

[OUTPUT]
    Name stdout
    Match *

[OUTPUT]
    Name  rsyslog 
    Match *
    Network tcp 
    Addr rsyslog:514
    Tag k8s-monitor 
```
