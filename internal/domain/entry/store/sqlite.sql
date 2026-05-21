create table if not exists `Entries` (
    `order`  blob not null,
    `item`   blob not null,
    `amount` int  not null,

    primary key (`order`, `item`),
    foreign key (`order`) references `Orders`(`uuid`),
    foreign key (`item`) references `Items`(`uuid`)
);
