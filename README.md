# Streaming

### Create database

```
$ cd app/src
$ sqlite3 stream.db
$ CREATE TABLE video (
    id INTEGER constraint video_pk primary key autoincrement,
    name        VARCHAR(100) default NULL,
    internal_id VARCHAR(10)  default NULL,
    created_at  TIMESTAMP    default CURRENT_TIMESTAMP,
    list        VARCHAR(50)  default NULL
);
$ .quit
```