create table list
(
    id INTEGER
        constraint list_pk
            primary key autoincrement,
    name varchar(100) default NULL not null
);

