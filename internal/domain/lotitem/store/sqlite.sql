create table if not exists `LotItems` (
	`uuid`      blob     not null primary key,
	`lot`       blob     not null,
	`material`  blob     not null,
	`amount`    int      not null,
	`unit_cost` int      not null,
	`expires`   datetime not null,
	`created`   datetime not null,
	`updated`   datetime not null,

	foreign key (`lot`) references `Lots`(`uuid`),
	foreign key (`material`) references `Materials`(`uuid`)
);

create index if not exists `LotItems_lot` on `LotItems`(`lot`);

create view if not exists `LotItems_View` as
	select
		i.`uuid`,
		i.`lot`,
		i.`material`,
		m.`name`,
		m.`unit`,
		i.`amount`,
		i.`unit_cost`,
		i.`expires`,
		i.`created`,
		i.`updated`
	from
		`Items` i
		join `Materials` m on i.`material` = m.`uuid`;

create trigger if not exists `LotItems_insert`
	before insert on `LotItems` for each row
	when exists (
		select
			1
		from
			`LotItems` i
			join `Lots` l on i.`lot` = l.`uuid`
		where
			i.`uuid` = new.`uuid`
			and l.`order` is not null
	)
	begin select raise(abort, 'immutable'); end;

create trigger if not exists `LotItems_update`
	before update on `LotItems` for each row
	when exists (
		select
			1
		from
			`LotItems` i
			join `Lots` l on i.`lot` = l.`uuid`
		where
			i.`uuid` = old.`uuid`
			and l.`order` is not null
	)
	begin select raise(abort, 'immutable'); end;

create trigger if not exists `LotItems_delete`
	before delete on `LotItems` for each row
	when exists (
		select
			1
		from
			`LotItems` i
			join `Lots` l on i.`lot` = l.`uuid`
		where
			i.`uuid` = old.`uuid`
			and l.`order` is not null
	)
	begin
		select raise(abort, 'immutable');
	end;
