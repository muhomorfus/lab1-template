-- +goose Up
-- +goose StatementBegin
create table persons (
    id int primary key generated always as identity,
    name text not null,
    age int,
    address text,
    work text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table persons;
-- +goose StatementEnd
