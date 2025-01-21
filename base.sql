create table users (
    id serial primary key,
    email varchar(255),
    password varchar(255),
    created_at timestamp default now(),
    updated_at timestamp default null
);

create table forgot_password (
    id serial primary key,
    token varchar(255),
    user_id int references users(id),
    created_at timestamp default now(),
    updated_at timestamp default null
)