create database tfm;
create user tfm with password 'password';
grant connect on database tfm to tfm;
\connect tfm
drop schema tfm;
create schema tfm;
alter schema tfm owner to tfm;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

