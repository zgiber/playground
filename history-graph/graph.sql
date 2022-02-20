drop table entity;
drop table version;

create extension if not exists "uuid-ossp";

-- every mutation increments the global version
create table if not exists "version" (
	"value" bigserial primary key,
	created_at timestamptz default now()
);

create table if not exists entity (
	id text default uuid_generate_v4(),
	"type" text, -- entity, identifier etc.
	"name" text, -- name (optional)
	properties jsonb, -- key:value
	created_with bigint not null references "version"("value"), --global version that created this entity
	superseded_with bigint references "version"("value") default null, --global version that modified this entity
	primary key (created_with, id)
);

create table if not exists entity_relation (
	-- id bigserial primary key,
	relation_type text,
	source_entity_id text not null,
	target_entity_id text not null,
	created_with bigint not null references "version"("value"),
	superseded_with bigint references "version"("value") default null
);

create table if not exists identifier (
	"namespace" text, -- e.g. SRay, ISIN etc.
	"value" text not null primary key,
	entity_id bigint not null,
	created_with bigint not null references "version"("value"),
	superseded_with bigint references "version"("value") default null
);

-- version must be incremented only if mutation is taking place:
-- - change of relation between entites
-- - change of any entity property
-- - ? change of an identifier (add new identifier, delete an identifier)

-- insert a few company entities
insert into entity (id,"type","name", created_with) values
('Company 1'),
('Company 2'),
('Company 3'),
('Company 4'),
('Company 5')
;

-- with

select * from graph;
