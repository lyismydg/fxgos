CREATE TABLE user
(
    id int PRIMARY KEY AUTO_INCREMENT,
    code nvarchar(50) NOT NULL,
    name nvarchar(100) NOT NULL,
    password varchar(50) NOT NULL,
    create_time datetime NOT NULL
);
CREATE UNIQUE INDEX user_code_uindex ON user (code);

CREATE TABLE trace_log
(
    id int PRIMARY KEY AUTO_INCREMENT,
    log_time datetime NOT NULL,
    user_id int NOT NULL,
    user nvarchar(50) NOT NULL,
    request_url nvarchar(100),
    code nvarchar(100) NOT NULL,
    type nvarchar(50) NOT NULL,
    message nvarchar(500),
    info text
);