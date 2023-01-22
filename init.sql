create table expression
(
    id serial not null
        constraint expression_pkey
            primary key,
    definition varchar(255) not null
);