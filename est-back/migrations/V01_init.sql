create table if not exists users
(
    id            uuid         not null primary key,
    username      varchar(64)  not null unique,
    password_hash varchar(512) not null,
    email         varchar(128) not null unique,
    avatar        varchar(64)
);

create index if not exists idx_username on users (username);

create table if not exists board
(
    id               uuid        not null primary key,
    name             varchar(64) not null,
    description      varchar(256),
    owner_id         uuid        not null,
    link_shared_mode varchar(16) not null,
    foreign key (owner_id) references users (id)
);

create table if not exists board_sharing
(
    id           uuid not null primary key,
    user_id      uuid not null,
    board_id     uuid not null,
    sharing_mode varchar(16),
    foreign key (user_id) references users (id),
    foreign key (board_id) references board (id)
);

create table if not exists figure
(
    id          uuid not null primary key,
    board_id    uuid not null,
    figure_data bytea,
    foreign key (board_id) references board (id)
);

create table if not exists recent_board
(
    id        uuid      not null primary key,
    board_id  uuid      not null,
    user_id   uuid      not null,
    last_used timestamp not null,
    foreign key (board_id) references board (id),
    foreign key (user_id) references users (id),
    unique (board_id, user_id)
);

