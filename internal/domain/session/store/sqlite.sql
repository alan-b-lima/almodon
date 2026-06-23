create table if not exists `Sessions` (
	`token`             blob     primary key,
	`user`              blob     not null,
	`hard_deadline`     datetime not null,
	`idle_deadline`     datetime not null,
	`password_verified` datetime not null,

	foreign key (`user`) references `Users`(`uuid`)
);

create unique index if not exists `Sessions_user` on `Sessions`(`user`);
