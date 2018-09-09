CREATE TABLE addresses (
  id          integer PRIMARY KEY,
  line1       varchar(40) NOT NULL,
  line2       varchar(40),
  city        varchar(40) NOT NULL,
  state       varchar(40) NOT NULL,
  zip         varchar(40) NOT NULL
);

CREATE TABLE guests (
  id          integer PRIMARY KEY,
  name        varchar(40) NOT NULL,
  email       varchar(40),
  address_id  integer REFERENCES addresses (id) NOT NULL
);

CREATE TABLE events (
  id          integer PRIMARY KEY,
  name        varchar(40) NOT NULL,
  address_id  integer REFERENCES addresses (id) NOT NULL
);

CREATE TABLE rsvps (
  id          integer PRIMARY KEY,
  head_count  integer NOT NULL,
  guest_id    integer REFERENCES guests (id) NOT NULL
);
