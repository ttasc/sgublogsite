drop table if exists users;
create table users (
    id int primary key auto_increment,
    email varchar(128) not null,
    fullname varchar(32) not null,
    password varchar(128) not null,
    session_token varchar(128),
    csrf_token varchar(128)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;
