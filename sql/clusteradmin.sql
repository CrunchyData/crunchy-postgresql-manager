drop database if exists clusteradmin;

create database clusteradmin;

\c clusteradmin;


create table settings (
	name varchar(30) primary key,
	value varchar(50) not null,
	description varchar(80) not null,
	updatedt timestamp not null
);

create table project (
	id serial primary key,
	name varchar(30) unique not null,
	description varchar(230),
	updatedt timestamp not null
);

insert into project (name, description, updatedt) values ('default', 'default project', now());

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
	name varchar(60) unique not null,
	clusterid int,
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


insert into settings (name, value, description, updatedt) values ('DOCKER-BRIDGES', '172.17.42.1:172.18.42.1', 'docker bridges to allow in the pg_hba.conf', now());
insert into settings (name, value, description, updatedt) values ('PG-DATA-PATH', '/var/cpm/data/pgsql', 'file path root of PG data files', now());
insert into settings (name, value, description, updatedt) values ('S-DOCKER-PROFILE-CPU', '256', 'small Docker profile CPU shares', now());
insert into settings (name, value, description, updatedt) values ('S-DOCKER-PROFILE-MEM', '512m', 'small Docker profile Memory limit', now());
insert into settings (name, value, description, updatedt) values ('M-DOCKER-PROFILE-CPU', '512', 'medium Docker profile CPU shares', now());
insert into settings (name, value, description, updatedt) values ('M-DOCKER-PROFILE-MEM', '1g', 'medium Docker profile Memory limit', now());
insert into settings (name, value, description, updatedt) values ('L-DOCKER-PROFILE-CPU', '0', 'large Docker profile CPU shares', now());
insert into settings (name, value, description, updatedt) values ('L-DOCKER-PROFILE-MEM', '0', 'large Docker profile Memory limit', now());
insert into settings (name, value, description, updatedt) values ('DOCKER-REGISTRY', 'registry:5000', 'not used currently', now());
insert into settings (name, value, description, updatedt) values ('PG-PORT', '5432', 'Postgresql port to use', now());
insert into settings (name, value, description, updatedt) values ('DOMAIN-NAME', 'crunchy.lab', 'domain name as configured in CPM', now());
insert into settings (name, value, description, updatedt) values ('ADMIN-URL', 'http://cpm:13001', 'CPM admin url', now());

insert into settings (name, value, description, updatedt) values ('CP-SM-COUNT', '1', 'small cluster profile standby count', now());
insert into settings (name, value, description, updatedt) values ('CP-SM-M-PROFILE', 'small', 'small cluster profile master Docker profile', now());
insert into settings (name, value, description, updatedt) values ('CP-SM-S-PROFILE', 'small', 'small cluster profile standby Docker profile', now());

insert into settings (name, value, description, updatedt) values ('CP-MED-COUNT', '1', 'medium cluster profile standby count', now());
insert into settings (name, value, description, updatedt) values ('CP-MED-M-PROFILE', 'small', 'medium cluster profile master Docker profile', now());
insert into settings (name, value, description, updatedt) values ('CP-MED-S-PROFILE', 'small', 'medium cluster profile standby Docker profile', now());

insert into settings (name, value, description, updatedt) values ('CP-LG-COUNT', '1', 'large cluster profile standby count', now());
insert into settings (name, value, description, updatedt) values ('CP-LG-M-PROFILE', 'small', 'large cluster profile master Docker profile', now());
insert into settings (name, value, description, updatedt) values ('CP-LG-S-PROFILE', 'small', 'large cluster profile standby Docker profile', now());
insert into settings (name, value, description, updatedt) values ('CP-SM-M-SERVER', 'low', 'small cluster profile master server size', now());
insert into settings (name, value, description, updatedt) values ('CP-SM-S-SERVER', 'low', 'small cluster profile standby server size', now());
insert into settings (name, value, description, updatedt) values ('CP-MED-M-SERVER', 'low', 'medium cluster profile master server size', now());
insert into settings (name, value, description, updatedt) values ('CP-MED-S-SERVER', 'low', 'medium cluster profile standby server size', now());
insert into settings (name, value, description, updatedt) values ('CP-LG-M-SERVER', 'low', 'large cluster profile master server size', now());
insert into settings (name, value, description, updatedt) values ('CP-LG-S-SERVER', 'low', 'large cluster profile standby server size', now());
insert into settings (name, value, description, updatedt) values ('CP-SM-ALGO', 'round-robin', 'small cluster placement algorithm', now());
insert into settings (name, value, description, updatedt) values ('CP-MED-ALGO', 'round-robin', 'medium cluster placement algorithm', now());
insert into settings (name, value, description, updatedt) values ('CP-LG-ALGO', 'round-robin', 'large cluster placement algorithm', now());
insert into settings (name, value, description, updatedt) values ('SLEEP-PROV', '2s', 'time to sleep during provisioning check', now());


insert into settings (name, value, description, updatedt) values ('TUNE-LG-MWM', '2GB', 'tuning parameter - maintenance_work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-CCT', '0.9', 'tuning parameter - checkpoint_completion_target', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-ECS', '24GB', 'tuning parameter - effective_cache_size', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-WM', '160MB', 'tuning parameter - work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-WB', '16MB', 'tuning parameter - wal_buffers', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-CS', '32', 'tuning parameter - checkpoint_segments', now());
insert into settings (name, value, description, updatedt) values ('TUNE-LG-SB', '8GB', 'tuning parameter - shared_buffers', now());

insert into settings (name, value, description, updatedt) values ('TUNE-MED-MWM', '1GB', 'tuning parameter - maintenance_work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-CCT', '0.9', 'tuning parameter - checkpoint_completion_target', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-ECS', '12GB', 'tuning parameter - effective_cache_size', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-WM', '80MB', 'tuning parameter - work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-WB', '16MB', 'tuning parameter - wal_buffers', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-CS', '32', 'tuning parameter - checkpoint_segments', now());
insert into settings (name, value, description, updatedt) values ('TUNE-MED-SB', '4GB', 'tuning parameter - shared_buffers', now());

insert into settings (name, value, description, updatedt) values ('TUNE-SM-MWM', '512MB', 'tuning parameter - maintenance_work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-CCT', '0.9', 'tuning parameter - checkpoint_completion_target', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-ECS', '6GB', 'tuning parameter - effective_cache_size', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-WM', '40MB', 'tuning parameter - work_mem', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-WB', '16MB', 'tuning parameter - wal_buffers', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-CS', '32', 'tuning parameter - checkpoint_segments', now());
insert into settings (name, value, description, updatedt) values ('TUNE-SM-SB', '2GB', 'tuning parameter - shared_buffers', now());

create table taskprofile (
	id serial primary key,
	name varchar(30) unique not null
);
insert into taskprofile (name) values ('pg_basebackup');
insert into taskprofile (name) values ('pg_dumpall');
insert into taskprofile (name) values ('vacuum-analyze');
insert into taskprofile (name) values ('backrest-backup');


create table taskschedule (
	id serial primary key,
	containername varchar(60) references container (name) on delete cascade not null,
	profilename varchar(30) references taskprofile (name) not null,
	name varchar(30) not null,
	enabled varchar(3) not null,
	minutes varchar(80) not null,
	hours varchar(80) not null,
	dayofmonth varchar(80) not null,
	month varchar(80) not null,
	dayofweek varchar(80) not null,
	restoreset varchar(80) not null,
	restoreremotepath varchar(80) not null,
	restoreremotehost varchar(80) not null,
	restoreremoteuser varchar(30) not null,
	restoredbuser varchar(80) not null,
	restoredbpass varchar(80) not null,
	serverip varchar(80) not null,
	updatedt timestamp not null,
	constraint valid_enabled check (
		enabled in ('YES', 'NO')
	)
);

create table taskstatus (
	id serial primary key,
	containername varchar(60) not null,
	profilename varchar(30) not null,
	scheduleid int references taskschedule (id) on delete cascade not null ,
	starttime timestamp not null,
	taskname varchar(30) not null,
	path varchar(80) not null,
	elapsedtime varchar(30) not null,
	tasksize varchar(30) not null,
	status varchar(200) not null,
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
	containername varchar(60) references container (name) on delete cascade not null,
	usename varchar(30) not null,
	passwd varchar(30) not null,
	updatedt timestamp not null,
	unique (containername, usename)
);

create table proxy (
	id serial primary key,
	containerid int references container (id) on delete cascade not null,
	projectid int references project (id) on delete cascade not null,
	port varchar(30) not null,
	host varchar(30) not null,
	databasename varchar(30) not null,
	usename varchar(30) not null,
	passwd varchar(30) not null,
	updatedt timestamp not null
);

insert into settings (name, value, description, updatedt) values ('POSTGRESPSW', '', 'postgres superuser password', now());
insert into settings (name, value, description, updatedt) values ('CPMTESTPSW', 'cpmtest', 'CPM test user password', now());
insert into settings (name, value, description, updatedt) values ('PGPOOLPSW', 'pgpool', 'pgpool user account password', now());


create table healthcheck (
	id serial primary key,
	projectname varchar(20) not null,
	projectid int references project (id) on delete cascade,
	containername varchar(60) not null,
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


