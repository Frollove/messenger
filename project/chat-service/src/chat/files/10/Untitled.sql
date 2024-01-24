CREATE TABLE "timetable" (
  "id" integer PRIMARY KEY,
  "route_id" integer,
  "date_departure" timestamp,
  "train_id" integer,
  "seats_free" integer,
  "seats_sold" integer
);

CREATE TABLE "worker" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "surname" varchar,
  "patronymic" varchar,
  "age" integer,
  "sex" varchar,
  "experience" integer,
  "pasport_id" integer,
  "team_id" integer
);

CREATE TABLE "worker_pasport" (
  "id" integer PRIMARY KEY,
  "number" integer,
  "recive_date" timestamp,
  "recieve_organization" varchar
);

CREATE TABLE "train" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "team_id" integer
);

CREATE TABLE "route" (
  "id" integer PRIMARY KEY,
  "arive_from" varchar,
  "arive_at" varchar
);

CREATE TABLE "works" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "duration" timestamp
);

CREATE TABLE "team" (
  "id" integer PRIMARY KEY,
  "name" varchar,
  "works_id" integer
);

CREATE TABLE "ticket" (
  "id" integer PRIMARY KEY,
  "seller_id" integer,
  "user_id" integer,
  "van" integer,
  "place" integer,
  "timetable_id" integer
);

CREATE TABLE "user" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "password" varchar,
  "name" varchar,
  "surname" varchar,
  "patronymic" varchar,
  "birthday" timestamp,
  "pasport_id" integer
);

CREATE TABLE "user_pasport" (
  "id" integer PRIMARY KEY,
  "number" integer,
  "recive_date" timestamp,
  "recieve_organization" varchar
);

ALTER TABLE "route" ADD FOREIGN KEY ("id") REFERENCES "timetable" ("route_id");

ALTER TABLE "train" ADD FOREIGN KEY ("id") REFERENCES "timetable" ("train_id");

ALTER TABLE "team" ADD FOREIGN KEY ("id") REFERENCES "train" ("team_id");

ALTER TABLE "works" ADD FOREIGN KEY ("id") REFERENCES "team" ("works_id");

ALTER TABLE "worker_pasport" ADD FOREIGN KEY ("id") REFERENCES "worker" ("pasport_id");

ALTER TABLE "team" ADD FOREIGN KEY ("id") REFERENCES "worker" ("team_id");

ALTER TABLE "worker" ADD FOREIGN KEY ("id") REFERENCES "ticket" ("seller_id");

ALTER TABLE "timetable" ADD FOREIGN KEY ("id") REFERENCES "ticket" ("timetable_id");

ALTER TABLE "ticket" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

ALTER TABLE "user_pasport" ADD FOREIGN KEY ("id") REFERENCES "user" ("pasport_id");
