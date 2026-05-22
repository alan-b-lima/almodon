create table if not exists `Lots` (
	`uuid`     blob not null primary key,
	`order`    blob,
	`supplier` text not null,
	`arrival`  datatime not null,
	`note`     text not null,
	`created`  datetime not null,
	`updated`  datetime not null,

	foreign key (`order`) references `Orders`(`uuid`)
);
