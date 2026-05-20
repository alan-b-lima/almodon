create table if not exists `Materials` (
	`uuid`     blob     primary key,
	`name`     text     not null,
	`ecampus`  integer  not null,
	`catmat`   integer  not null,
	`siads`    integer  not null,
	`descript` text     not null,
	`unit`     text     not null,
	`min`      int      not null,
	`created`  datetime not null,
	`updated`  datetime not null
);
