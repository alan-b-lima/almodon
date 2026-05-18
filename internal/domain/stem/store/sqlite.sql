create table if not exists `Stems` (
    `uuid`    blob     not null primary key,
    `bloom`   blob,
    `title`   text     not null,
    `created` datetime not null,

    foreign key (`bloom`) references `Blooms`(`uuid`)
);

create view if not exists `Stems_View` as
    select
        s.`uuid`,
        s.`bloom`,
        s.`title`,
        o.`version`,
        s.`created`,
        o.`created` as 'updated'
    from
        `Stems` s
        join `Orders` o on s.`bloom` = o.`uuid`;