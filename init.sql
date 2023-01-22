create table expression
(
    id serial not null
        constraint expression_pkey
            primary key,
    description string
);