
create extension pg_stat_statements;

create user cpmtest login superuser password 'cpmtest';

create user pgpool login superuser password 'pgpool';

create database cpmtest;

grant all privileges on database cpmtest to pgpool;
grant all privileges on database cpmtest to cpmtest;

\c cpmtest;

create extension adminpack;
create extension pg_stat_statements;


DROP TABLE IF EXISTS loadtest;
CREATE TABLE loadtest(
	id int,
	name varchar(200)
);

