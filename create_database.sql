--
-- Drop tables
--
drop table if exists menu  ;
drop table if exists service;
drop table if exists host;
drop table if exists probe  ;

--
-- Create probe table
--
create table probe (
  key character(20) primary key not null unique,
  data jsonb
);
--
-- Define generic demo probe
--
insert into probe
values ('WPPAGE','{"name": "WordPress pages probe", "description": "Probe if WordPress site is up and running"}');
insert into probe
values ('PING', '{"name": "Basic ping probe", "description": "If IP address respond of ping", "command": "ping {ip}"}');

--
-- Create menu table
--
create table menu (
  key  character(20) primary key not null unique,
  priority integer,
  data jsonb
);

insert into menu
 values ('HOME', 1000, '{"Key": "1000HOME", "Link": "/", "Text": "Home", "Visible": true }');
insert into menu
 values ('LOGOUT', 9000, '{"Key": "9000LOGOUT", "Link": "/logout", "Text": "Logout", "Visible": false }');
insert into menu
 values ('PROBE', 2000, '{"Key": "2000PROBE", "Link": "/probes", "Text": "Probe", "Visible": true }');

--
-- Define host
--
create table host (
  key character(20) primary key not null unique,
  fqdn character(255),
  ipv4  character(15),
  data jsonb
);

insert into host
 values ('EVEUTERPE', 'ev-euterpe.ch', '', '{}');

--
-- Define services
--
create table service (
  key character(20) not null,
  host character(20) references host(key),
  probe character(20) references probe(key),
  interval integer default 15,
  data jsonb,
  primary key (key, host)
);
insert into service
 values ('EVEPING', 'EVEUTERPE', 'PING', 15, '{}');

