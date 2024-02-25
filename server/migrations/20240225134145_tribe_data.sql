-- +goose Up
-- +goose StatementBegin
insert into tribes values(
	0,
	unixepoch(),
	"Maximilian Kolbe",
	"Maxko",
	"13/12/07",
	null,
	"maxkolbe.de",
	null
);
insert into tribes values(
	1,
	unixepoch(),
	"Sankt Canisius",
	"Canisius",
	"13/12/13",
	null,
	null,
	null
);
insert into tribes values(
	2,
	unixepoch(),
	"Sankt Ansgar",
	"Ansgar",
	"13/12/08",
	null,
	null,
	null
);
insert into tribes values(
	3,
	unixepoch(),
	"Frieden Christi",
	"FC",
	"13/12/09",
	null,
	null,
	null
);
insert into tribes values(
	4,
	unixepoch(),
	"Pater Rupert Mayer",
	"PRM",
	"13/12/06",
	null,
	null,
	null
);
insert into tribes values(
	5,
	unixepoch(),
	"Maria Hilf",
	"Maria Hilf",
	"13/12/04",
	null,
	null,
	null
);
insert into tribes values(
	6,
	unixepoch(),
	"Sankt Anna",
	"St.Anna",
	"13/12/12",
	null,
	null,
	null
);
insert into tribes values(
	7,
	unixepoch(),
	"Heilig Engel",
	"Hl.Engel",
	"13/12/02",
	null,
	null,
	null
);
insert into tribes values(
	8,
	unixepoch(),
	"Swapingo",
	"Swapingo",
	"13/12/10",
	null,
	null,
	null
);
insert into tribes values(
	9,
	unixepoch(),
	"Sankt Severin",
	"Severin",
	"13/12/03",
	null,
	null,
	null
);
insert into tribes values(
	10,
	unixepoch(),
	"Heilig Kreuz",
	"Hl.Kreuz",
	"13/12/14",
	null,
	null,
	null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from tribes;
-- +goose StatementEnd
