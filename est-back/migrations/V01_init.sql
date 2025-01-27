create or replace function set_updated()
    returns trigger as
$$
begin
    new.updated = now();
    return new;
end;
$$ language plpgsql;

create table if not exists users
(
    id            uuid         not null primary key,
    created       timestamp    not null default now(),
    updated       timestamp    not null default now(),
    username      varchar(64)  not null unique,
    password_hash varchar(512) not null,
    email         varchar(128) not null unique,
    avatar        varchar(64),
    last_login    timestamp    not null default now()
);

create index if not exists idx_username on users (username);

create trigger update_users_timestamp
    before update
    on users
    for each row
execute function set_updated();

create table if not exists board
(
    id               uuid        not null primary key,
    created          timestamp   not null default now(),
    updated          timestamp   not null default now(),
    name             varchar(64) not null,
    description      varchar(256),
    owner_id         uuid        not null,
    link_shared_mode varchar(16) not null,
    foreign key (owner_id) references users (id)
);

create trigger update_board_timestamp
    before update
    on board
    for each row
execute function set_updated();

create table if not exists board_sharing
(
    id           uuid      not null primary key,
    created      timestamp not null default now(),
    updated      timestamp not null default now(),
    user_id      uuid      not null,
    board_id     uuid      not null,
    sharing_mode varchar(16),
    foreign key (user_id) references users (id),
    foreign key (board_id) references board (id)
);

create trigger update_board_sharing_timestamp
    before update
    on board_sharing
    for each row
execute function set_updated();


create table if not exists figure
(
    id          uuid      not null primary key,
    created     timestamp not null default now(),
    updated     timestamp not null default now(),
    board_id    uuid      not null,
    figure_data bytea,
    foreign key (board_id) references board (id)
);

create trigger update_figure_timestamp
    before update
    on figure
    for each row
execute function set_updated();

create table if not exists recent_board
(
    id        uuid      not null primary key,
    created   timestamp not null default now(),
    updated   timestamp not null default now(),
    board_id  uuid      not null,
    user_id   uuid      not null,
    last_used timestamp not null,
    foreign key (board_id) references board (id),
    foreign key (user_id) references users (id),
    unique (board_id, user_id)
);

create trigger update_recent_board_timestamp
    before update
    on recent_board
    for each row
execute function set_updated();
