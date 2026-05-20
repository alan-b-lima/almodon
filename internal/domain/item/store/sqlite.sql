create table if not exists `Items` (
	`uuid`      blob     not null primary key,
	`material`  blob     not null,
	`lot`       blob     not null,
	`available` real     not null,
	`unit_cost` int      not null,
	`expires`   datetime not null,
	`created`   datetime not null,
	`updated`   datetime not null,

	foreign key (`material`) references `Material`(`uuid`),
	foreign key (`lot`) references `Lot`(`uuid`)
);

create view if not exists `Items_View` as
	select
		i.`uuid`,
		i.`material`,
		m.`name`,
		m.`ecampus`,
		m.`catmat`,
		m.`siads`,
		m.`unit`,
		i.`lot`,
		i.`available`,
		i.`unit_cost`,
		i.`expires`,
		i.`created`,
		i.`updated`
	from
		`Items` i
		join `Materials` m on i.`material` = m.`uuid`;
