create table if not exists `Stem` (
    `uuid`    blob     not null primary key,
    `bloom`   blob,
    `created` datetime not null,

    foreign key (`bloom`) references `Orders`(`uuid`)
);

create view if not exists `Stem_View` as
    select
        s.`uuid`,
        s.`bloom`,
        coalesce(o.`version`, 0) as 'version',
        s.`created`,
        coalesce(o.`created`, s.`created`) as 'updated'
    from
        `Stem` s
        left join `Orders` o on s.`bloom` = o.`uuid`;