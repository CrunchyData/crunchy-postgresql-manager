drop database if exists clusteradmin;

create database clusteradmin;

\c clusteradmin;


create table settings (
	name varchar(30) primary key,
	value varchar(50) not null,
	updatedt timestamp not null
);

create table project (
	id serial primary key,
	name varchar(30) unique not null,
	description varchar(230),
	updatedt timestamp not null
);

insert into project (name, description, updatedt) values ('default', 'default project', now());

create table server (
	id serial primary key,
	name varchar(20) unique not null,
	ipaddress varchar(20) unique not null,
	dockerbip varchar(20) unique not null,
	pgdatapath varchar(40) not null,
	serverclass varchar(20) not null,
	createdt timestamp not null,
	constraint valid_server_class check (
		serverclass in ('low', 'medium', 'high')
	)
);

create table cluster (
	id serial primary key,
	name varchar(20) unique not null,
	clustertype varchar(20) not null,
	status varchar(20) not null,
	createdt timestamp not null,
	projectid int references project (id) on delete cascade,
	constraint valid_cluster_types check (
		clustertype in ('bdr', 'asynchronous', 'synchronous')
	),
	constraint valid_status_types check (
		status in ('uninitialized', 'initialized')
	)
);

create table container (
	id serial primary key,
	name varchar(30) unique not null,
	clusterid int,
	serverid int references server (id) on delete cascade,
	projectid int references project (id) on delete cascade,
	role varchar(10) not null,
	image varchar(30) not null,
	createdt timestamp not null,
	constraint valid_roles check (
		role in ('standby', 'master', 'unassigned', 'standalone', 'pgpool')
	)
);

create table secuser (
	name varchar(20) not null primary key,
	password varchar(20) not null,
	updatedt timestamp not null);

create table secsession (
	token varchar(50) not null primary key,
	name varchar(20) references secuser (name) on delete cascade,
	updatedt timestamp not null);

create table secrole (
	name varchar(30) not null primary key,
	updatedt timestamp not null);


create table secuserrole (
	username varchar(20) references secuser (name) on delete cascade,
	role varchar(30) references secrole (name) on delete cascade,
	unique (username, role)
);

create table secperm (
	name varchar(20) not null primary key,
	description varchar(50) not null
);

create table secroleperm (
	role varchar(20) references secrole (name) on delete cascade,
	perm varchar(30) references secperm (name) on delete cascade,
	unique (role, perm)
);


insert into secuser values ('cpm', 'dd6ced', now());

insert into secrole values ('superuser', now());

insert into secuserrole values ('cpm', 'superuser');

insert into secperm values ('perm-server', 'maintain servers');
insert into secperm values ('perm-container', 'maintain containers');
insert into secperm values ('perm-cluster', 'maintain clusters');
insert into secperm values ('perm-setting', 'maintain settings');
insert into secperm values ('perm-backup', 'perform backups');
insert into secperm values ('perm-user', 'maintain users');

insert into secroleperm values ('superuser', 'perm-server');
insert into secroleperm values ('superuser', 'perm-container');
insert into secroleperm values ('superuser', 'perm-cluster');
insert into secroleperm values ('superuser', 'perm-setting');
insert into secroleperm values ('superuser', 'perm-backup');
insert into secroleperm values ('superuser', 'perm-user');


insert into settings (name, value, updatedt) values ('S-DOCKER-PROFILE-CPU', '256', now());
insert into settings (name, value, updatedt) values ('S-DOCKER-PROFILE-MEM', '512m', now());
insert into settings (name, value, updatedt) values ('M-DOCKER-PROFILE-CPU', '512', now());
insert into settings (name, value, updatedt) values ('M-DOCKER-PROFILE-MEM', '1g', now());
insert into settings (name, value, updatedt) values ('L-DOCKER-PROFILE-CPU', '0', now());
insert into settings (name, value, updatedt) values ('L-DOCKER-PROFILE-MEM', '0', now());
insert into settings (name, value, updatedt) values ('DOCKER-REGISTRY', 'registry:5000', now());
insert into settings (name, value, updatedt) values ('PG-PORT', '5432', now());
insert into settings (name, value, updatedt) values ('DOMAIN-NAME', 'crunchy.lab', now());
insert into settings (name, value, updatedt) values ('ADMIN-URL', 'http://cpm:13001', now());

insert into settings (name, value, updatedt) values ('CP-SM-COUNT', '1', now());
insert into settings (name, value, updatedt) values ('CP-SM-M-PROFILE', 'small', now());
insert into settings (name, value, updatedt) values ('CP-SM-S-PROFILE', 'small', now());

insert into settings (name, value, updatedt) values ('CP-MED-COUNT', '1', now());
insert into settings (name, value, updatedt) values ('CP-MED-M-PROFILE', 'small', now());
insert into settings (name, value, updatedt) values ('CP-MED-S-PROFILE', 'small', now());

insert into settings (name, value, updatedt) values ('CP-LG-COUNT', '1', now());
insert into settings (name, value, updatedt) values ('CP-LG-M-PROFILE', 'small', now());
insert into settings (name, value, updatedt) values ('CP-LG-S-PROFILE', 'small', now());
insert into settings (name, value, updatedt) values ('CP-SM-M-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-SM-S-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-MED-M-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-MED-S-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-LG-M-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-LG-S-SERVER', 'low', now());
insert into settings (name, value, updatedt) values ('CP-SM-ALGO', 'round-robin', now());
insert into settings (name, value, updatedt) values ('CP-MED-ALGO', 'round-robin', now());
insert into settings (name, value, updatedt) values ('CP-LG-ALGO', 'round-robin', now());

create table backupprofile (
	id serial primary key,
	name varchar(30) unique not null
);
insert into backupprofile (name) values ('pg_basebackup');
insert into backupprofile (name) values ('pg_dumpall');


create table backupschedule (
	id serial primary key,
	serverid int references server (id) on delete cascade not null,
	containername varchar(20) references container (name) on delete cascade not null,
	profilename varchar(30) references backupprofile (name) not null,
	name varchar(30) not null,
	enabled varchar(3) not null,
	minutes varchar(80) not null,
	hours varchar(80) not null,
	dayofmonth varchar(80) not null,
	month varchar(80) not null,
	dayofweek varchar(80) not null,
	updatedt timestamp not null,
	constraint valid_enabled check (
		enabled in ('YES', 'NO')
	)
);

create table backupstatus (
	id serial primary key,
	containername varchar(30) not null,
	profilename varchar(30) not null,
	scheduleid int references backupschedule (id) on delete cascade not null ,
	starttime timestamp not null,
	backupname varchar(30) not null,
	servername varchar(20) not null,
	serverip varchar(20) not null,
	path varchar(80) not null,
	elapsedtime varchar(30) not null,
	backupsize varchar(30) not null,
	status varchar(50) not null,
	updatedt timestamp not null,

	unique (containername, starttime)
);


drop table monmetric;
drop table monschedule;

create table monschedule (
	name varchar(30) not null,
	cronexp varchar(80) not null,
	unique (name)
);

insert into monschedule values ( 's1', '@every 0h5m0s');


create table monmetric (
	name varchar(30) unique not null,
	metrictype varchar(30) not null,
	schedule varchar(30) references monschedule (name),
	constraint valid_metrictype check (
		metrictype in ('server', 'database', 'healthck')
	)
);

insert into monmetric values ('cpu', 'server', 's1');
insert into monmetric values ('mem', 'server', 's1');
insert into monmetric values ('pg1', 'database', 's1');
insert into monmetric values ('pg2', 'database', 's1');
insert into monmetric values ('hc1', 'healthck', 's1');

create table containeruser (
	id serial primary key,
	containername varchar(20) references container (name) on delete cascade not null,
	usename varchar(30) not null,
	passwd varchar(30) not null,
	updatedt timestamp not null,
	unique (containername, usename)
);

insert into settings (name, value, updatedt) values ('POSTGRESPSW', '', now());
insert into settings (name, value, updatedt) values ('CPMTESTPSW', 'cpmtest', now());
insert into settings (name, value, updatedt) values ('PGPOOLPSW', 'pgpool', now());


create table healthcheck (
	id serial primary key,
	projectname varchar(20) not null,
	projectid int references project (id) on delete cascade,
	containername varchar(20) not null,
	containerid int references container (id) on delete cascade,
	containerrole varchar(10) not null,
	containerimage varchar(30) not null,
	status varchar(20) not null,
	updatedt timestamp not null
);

create table accessrule (
	id serial primary key,
	name varchar(30) unique not null,
	ruletype varchar(20),
	database varchar(60),
	ruleuser varchar(60),
	address varchar(60),
	method varchar(20),
	description varchar(230),
	createdt timestamp not null,
	updatedt timestamp not null
);

create table containeraccessrule (
	id serial primary key,
	containerid int references container (id) on delete cascade,
	accessruleid int references accessrule (id) on delete cascade,
	createdt timestamp not null,
	unique (containerid, accessruleid)
);

insert into  accessrule (
name , ruletype, database, ruleuser, address, method, description,
createdt, updatedt) values (
'samplerule' , 'host', 'all', 'all', '192.168.10.100/32', 'md5', 'sample rule',
now(), now());


