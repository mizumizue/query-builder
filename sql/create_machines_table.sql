create table machines
(
	machine_number integer auto_increment,
	machine_name varchar(20) not null,
	buy_date date not null,
	price float not null,
	owner varchar(20) not null,
	constraint machines_pk
		primary key (machine_number)
);
