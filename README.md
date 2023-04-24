# calproxy

Proxy for CALDAV clients to force delete events.

Usage:
```bash
Usage of ./calproxy:
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
  -secure
        use https listener
  -target string
        target URL (default "https://calendar.mail.ru")
```

Add to `deleted.txt` UIDs of events you want to delete.

## Adding UUId 

1. Delete event from calendar
2. Watch for `DELETE` request in log
3. Add UUId to `deleted.txt`, e.g. for `DELETE /calendars/A-B-C-D/ff-aa-ee.ics` add `ff-aa-ee`