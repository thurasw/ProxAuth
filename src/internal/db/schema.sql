CREATE TABLE IF NOT EXISTS users (
	id INTEGER NOT NULL,
	username TEXT NOT NULL,
    password TEXT NOT NULL,
	PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS apps (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,
	color TEXT DEFAULT '#000' NOT NULL,
    LOGO BLOB,
	PRIMARY KEY (id)
);
