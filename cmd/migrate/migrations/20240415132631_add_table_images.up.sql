create table if not exists images(
    id serial primary key,
    user_id uuid references auth.users,
    batch_id uuid not null,
    status int not null default 1,
    prompt text not null,
    image_location text,
    created_at timestamp not null default now(),
    deleted boolean not null default 'false',
    deleted_at timestamp
);