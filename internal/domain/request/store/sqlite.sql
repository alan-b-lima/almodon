create table if not exists `Requests` (
	`uuid`    blob     not null primary key,
	`number`  integer  not null,
	`author`  text     not null,
	`title`   text     not null,
	`memo`    text     not null,
	`status`  integer  not null,
	`created` datetime not null,
	`updated` datetime not null,

	foreign key (`author`) references `Users`(`name`)
);

create index if not exists `Requests_number` on `Requests`(`number`);

create view if not exists `Requests_View` as
	select
		r.`uuid`,
		r.`number`,
		r.`author`,
		u.`name`,
		u.`siape`,
		r.`title`,
		r.`memo`,
		r.`status`,
		r.`created`,
		r.`updated`
	from
		`Requests` r
		join `Users` u on r.`author` = u.`uuid`;
