create table if not exists `Stems` (
    `uuid`    blob     not null primary key,
    `bloom`   blob,
    `name`    text     not null,
    `created` datetime not null,

    foreign key (`bloom`) references `Orders`(`uuid`)
);

create view if not exists `Stems_View` as
    select
        s.`uuid`,
        s.`bloom`,
        s.`name`,
        coalesce(o.`version`, 0) as 'version',
        s.`created`,
        coalesce(o.`created`, s.`created`) as 'updated'
    from
        `Stems` s
        left join `Orders` o on s.`bloom` = o.`uuid`;