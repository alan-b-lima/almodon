create table if not exists `Orders` (
    `uuid`    blob     not null primary key,
    `version` int      not null,
    `stem`    blob     not null,
    `author`  blob     not null,
    `created` datetime not null,

    unique (`version`, `stem`),
    foreign key (`stem`) references `Stems`(`uuid`),
    foreign key (`author`) references `Users`(`uuid`)
);

create view if not exists `Orders_View` as
    select
        o.`uuid`,
        o.`version`,
        o.`stem`,
        s.`title`,
        o.`author`,
        u.`name`,
        u.`siape`,
        o.`created`
    from
        `Orders` o
        join `Stems` s on o.`stem` = s.`uuid`
        join `Users` u on o.`author` = u.`uuid`;

create view if not exists `Orders_Bloom_View` as
    select
        o.`uuid`,
        o.`version`,
        o.`stem`,
        s.`title`,
        o.`author`,
        u.`name`,
        u.`siape`,
        o.`created`
    from
        `Stems` s
        join `Orders` o on s.`bloom` = o.`uuid`
        join `Users` u on o.`author` = u.`uuid`;
