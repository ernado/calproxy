# calproxy

Proxy for CALDAV clients to force delete events.

Usage:
```bash
Usage of ./calproxy:
  -append
        automatically append deleted UIDs to file
  -cert string
        path to certificate (default "server.crt")
  -deleted string
        file with deleted UIDs (default "deleted.txt")
  -dump
        dump requests/responses
  -key string
        path to key (default "server.key")
  -listen string
        listen address (default "localhost:8080")
  -output string
        directory to store requests/responses (default "output")
  -rewrite
        rewrite events (default true)
  -secure
        use https listener
  -target string
        target URL (default "https://calendar.mail.ru")
```

Add to `deleted.txt` UIDs of events you want to delete.

This can be done automatically if `-append` flag is `true` which is on by default.

## Adding UUId 

1. Delete event from calendar
2. Watch for `DELETE` request in log
3. Add UUId to `deleted.txt`, e.g. for `DELETE /calendars/A-B-C-D/ff-aa-ee.ics` add `ff-aa-ee`

## Systemd example

Use `systemctl --user` to create service for current user.
Local proxy is the most secure way to use this proxy.

```bash
systemctl --user cat caldav.service 
# /home/ernado/.config/systemd/user/caldav.service
[Unit]
Description=CalDav Proxy
StartLimitIntervalSec=0

[Service]
ExecStart=/src/ernado/calproxy/calproxy --listen localhost:8312
Type=simple
RestartSec=3
Restart=always
WorkingDirectory=/src/ernado/calproxy

[Install]
WantedBy=default.target
```