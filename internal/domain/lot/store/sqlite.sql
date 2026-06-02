create table if not exists `Lots` (
	`uuid`     blob not null primary key,
	`order`    blob,
	`supplier` text not null,
	`author`   blob not null,
	`arrival`  datetime not null,
	`note`     text not null,
	`created`  datetime not null,
	`updated`  datetime not null,

	foreign key (`order`) references `Orders`(`uuid`),
	foreign key (`author`) references `User`(`uuid`)
);

create trigger if not exists `Lots_update`
	before update on `Lots` for each row
	when
		old.`order` is not null
	begin
		select raise(abort, 'immutable');
	end;

create trigger if not exists `Lots_delete`
	before update on `Lots` for each row
	when
		old.`order` is not null
	begin
		select raise(abort, 'immutable');
	end;
