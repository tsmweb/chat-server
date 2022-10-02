-- DROP ROLE salesapi;

CREATE ROLE salesapi WITH 
	NOSUPERUSER
	NOCREATEDB
	NOCREATEROLE
	INHERIT
	LOGIN
	NOREPLICATION
	NOBYPASSRLS
	CONNECTION LIMIT -1;

ALTER user salesapi ENCRYPTED PASSWORD 'password';

-- DROP SCHEMA chat_db;

CREATE SCHEMA chat_db AUTHORIZATION salesapi;

-- DROP TABLE chat_db."user";

CREATE TABLE chat_db."user" (
	id varchar(100) NOT NULL,
	"name" varchar(255) NOT NULL,
	lastname varchar(255) NULL,
	created_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp(0) NULL,
	CONSTRAINT client_pkey PRIMARY KEY (id)
);

-- DROP TABLE chat_db.login;

CREATE TABLE chat_db.login (
	user_id varchar(100) NOT NULL,
	"password" varchar(100) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL,
	CONSTRAINT user_pkey PRIMARY KEY (user_id)
);

-- chat_db.login foreign keys

ALTER TABLE chat_db.login ADD CONSTRAINT login_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db.contact;

CREATE TABLE chat_db.contact (
	user_id varchar(100) NOT NULL,
	contact_id varchar(100) NOT NULL,
	"name" varchar(255) NOT NULL,
	lastname varchar(255) NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL,
	CONSTRAINT contact_pkey PRIMARY KEY (user_id, contact_id)
);

-- chat_db.contact foreign keys

ALTER TABLE chat_db.contact ADD CONSTRAINT contact_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db.online_user;

CREATE TABLE chat_db.online_user (
	user_id varchar(100) NOT NULL,
	server_id varchar(255) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT online_user_pkey PRIMARY KEY (user_id)
);

-- chat_db.online_user foreign keys

ALTER TABLE chat_db.online_user ADD CONSTRAINT user_online_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db.blocked_user;

CREATE TABLE chat_db.blocked_user (
	user_id varchar(100) NOT NULL,
	blocked_user_id varchar(100) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT blocked_user_pkey PRIMARY KEY (user_id, blocked_user_id)
);

-- chat_db.blocked_user foreign keys

ALTER TABLE chat_db.blocked_user ADD CONSTRAINT contact_block_contact_id_fkey FOREIGN KEY (blocked_user_id) REFERENCES chat_db."user"(id);
ALTER TABLE chat_db.blocked_user ADD CONSTRAINT contact_block_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db."group";

CREATE TABLE chat_db."group" (
	id varchar(100) NOT NULL,
	owner_id varchar(100) NOT NULL,
	"name" varchar(255) NOT NULL,
	description varchar(255) NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL,
	updated_by varchar(100) NULL,
	CONSTRAINT group_pkey PRIMARY KEY (id)
);

-- chat_db."group" foreign keys

ALTER TABLE chat_db."group" ADD CONSTRAINT group_user_id_fkey FOREIGN KEY (owner_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db.group_member;

CREATE TABLE chat_db.group_member (
	group_id varchar(100) NOT NULL,
	user_id varchar(100) NOT NULL,
	created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL,
	updated_by varchar(100) NULL,
	"admin" bool NOT NULL,
	CONSTRAINT group_member_pkey PRIMARY KEY (group_id, user_id)
);

-- chat_db.group_member foreign keys

ALTER TABLE chat_db.group_member ADD CONSTRAINT group_member_group_id_fkey FOREIGN KEY (group_id) REFERENCES chat_db."group"(id);
ALTER TABLE chat_db.group_member ADD CONSTRAINT group_member_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);


-- DROP TABLE chat_db.group_member_notify;

CREATE TABLE chat_db.group_member_notify (
	group_id varchar(100) NOT NULL,
	user_id varchar(100) NOT NULL,
	created_at timestamp(0) NOT NULL,
	CONSTRAINT group_member_notify_pk PRIMARY KEY (group_id, user_id)
);

-- chat_db.group_member_notify foreign keys

ALTER TABLE chat_db.group_member_notify ADD CONSTRAINT group_member_notify_group_id_fkey FOREIGN KEY (group_id) REFERENCES chat_db."group"(id);
ALTER TABLE chat_db.group_member_notify ADD CONSTRAINT group_member_notify_user_id_fkey FOREIGN KEY (user_id) REFERENCES chat_db."user"(id);

-- DROP TABLE chat_db.offline_message;

CREATE TABLE chat_db.offline_message (
	msg_id varchar(100) NOT NULL,
	msg_status bpchar(1) NOT NULL,
	msg_from varchar(100) NOT NULL,
	msg_to varchar(100) NOT NULL,
	msg_group varchar(100) NULL,
	msg_date timestamp NOT NULL,
	msg_content_type varchar(10) NOT NULL,
	msg_content text NOT NULL,
	CONSTRAINT offline_message_pkey PRIMARY KEY (msg_id, msg_status)
);
CREATE INDEX offline_message_msg_to_idx ON chat_db.offline_message USING btree (msg_to, msg_status);

-- chat_db.group_member_notify foreign keys
