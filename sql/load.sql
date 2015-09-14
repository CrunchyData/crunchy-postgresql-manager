create table customer (
	id serial primary key,
	name varchar(60) not null,
	location varchar(220) not null
);


create table product (
	id serial primary key,
	productname varchar(30) not null,
	customerid int references customer (id) on delete cascade,
	productdesc varchar(100) not null
);

