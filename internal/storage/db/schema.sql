create function trigger_set_timestamp() returns trigger
	language plpgsql
as $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$;

alter function trigger_set_timestamp() owner to gzlgiwbcviwknc;

create trigger set_timestamp
    before update
    on categories
    for each row
    execute procedure trigger_set_timestamp();

create trigger set_timestamp
    before update
    on notifications
    for each row
    execute procedure trigger_set_timestamp();

create trigger set_timestamp
    before update
    on observation_state
    for each row
    execute procedure trigger_set_timestamp();

create trigger set_timestamp
    before update
    on observations
    for each row
    execute procedure trigger_set_timestamp();


create trigger set_timestamp
    before update
    on species
    for each row
    execute procedure trigger_set_timestamp();

create trigger set_timestamp
    before update
    on states
    for each row
    execute procedure trigger_set_timestamp();

create trigger set_timestamp
    before update
    on volunteers
    for each row
    execute procedure trigger_set_timestamp();

create table categories
(
    id serial not null
        constraint categories_pk
        primary key,
    name varchar(255) not null,
    created_at timestamp(6) not null,
    updated_at timestamp(6) not null
);

alter table categories owner to mipqivqxykcawh;

create table species
(
    id serial not null
        constraint species_pk
        primary key,
    name varchar(255) not null,
    common_name varchar(255) not null,
    created_at timestamp(6) not null,
    updated_at timestamp(6) not null
);

alter table species owner to mipqivqxykcawh;

create table observations
(
    id bigint not null
        constraint observations_pk
        primary key,
    link text not null,
    "species_ID" bigint not null
        constraint observations_fk0
        references species,
    created_at timestamp not null,
    updated_at timestamp not null
);

alter table observations owner to mipqivqxykcawh;

create table states
(
    id serial not null
        constraint states_pk
        primary key,
    category_id serial not null
        constraint states_fk0
        references categories,
    name varchar(255) not null,
    created_at timestamp(6) not null,
    updated_at timestamp(6) not null
);

alter table states owner to mipqivqxykcawh;

create table volunteers
(
    id serial not null
        constraint volunteers_pk
        primary key,
    fname varchar(255) not null,
    lname varchar(255) not null,
    role_id integer not null,
    created_at timestamp(6) not null,
    updated_at timestamp(6) not null
);

alter table volunteers owner to mipqivqxykcawh;

create table observation_state
(
    id serial not null
        constraint observation_state_pk
        primary key,
    observation_id serial not null
        constraint observation_state_fk0
        references observations,
    volunteer_id bigint not null
        constraint observation_state_fk1
        references volunteers,
    state_id bigint not null
        constraint observation_state_fk2
        references states,
    comment text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

alter table observation_state owner to mipqivqxykcawh;

create table notifications
(
    id serial not null
        constraint notifications_pk
        primary key,
    observation_id serial not null
        constraint notifications_fk0
        references observations,
    created_at timestamp not null,
    updated_at timestamp not null,
    comment text not null,
    link text not null
);

alter table notifications owner to mipqivqxykcawh;

