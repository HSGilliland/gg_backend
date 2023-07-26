CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--==================================================================================================================

--CREATE TABLES

--==================================================================================================================

CREATE TABLE persontype
(
    code varchar(6) PRIMARY KEY
    , name varchar(64) NOT NULL
    , description varchar(256)
    , createddatetime timestamp DEFAULT CURRENT_TIMESTAMP
    , isdeleted boolean DEFAULT false
);

CREATE TABLE person
(
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY
    , persontypecode varchar(6) REFERENCES persontype(code)
    , name varchar(64) NULL
    , surname varchar(64) NULL
    , email varchar(64) NULL
    , username varchar(64) NULL
    , password varchar(64) NULL
    , createddatetime timestamp DEFAULT CURRENT_TIMESTAMP
    , isdeleted boolean DEFAULT false
);

--==================================================================================================================

--INSERT DEFAULT DATA

--==================================================================================================================

INSERT INTO persontype(
	code, name, description, createddatetime, isdeleted)
	VALUES ('TCH', 'Teacher', 'Teacher that could upload content', CURRENT_TIMESTAMP, false);

INSERT INTO persontype(
	code, name, description, createddatetime, isdeleted)
	VALUES ('STD', 'Student', 'Student that could consume content.', CURRENT_TIMESTAMP, false);

INSERT INTO person(
	id, persontypecode, name, surname, username, password, createddatetime, isdeleted)
	VALUES ('584dac8e-e6f9-4146-a204-b8f2dd260a2f', 'STD', 'Johan', 'Stemmet', 'jstemmet', 'mnandi', CURRENT_TIMESTAMP, false);

INSERT INTO person(
	id, persontypecode, name, surname, username, password, createddatetime, isdeleted)
	VALUES ('46c7e785-9f51-4db0-92ab-514be41be0b4', 'TCH', 'Mnandi', 'Gilliland', 'mmgilliland', 'hermanus', CURRENT_TIMESTAMP, false);
    
--==================================================================================================================

--CREATE FUNCTIONS

--==================================================================================================================

--VALIDATE LOGIN

CREATE OR REPLACE FUNCTION login(
	var_username varchar(64),
	var_password varchar(64),
	OUT ret_person_exists boolean,
	OUT ret_person_type varchar(6),
    OUT ret_person_id uuid) 
	LANGUAGE 'plpgsql'
AS $BODY$
BEGIN
    
	select into ret_person_exists exists(select 1 from person where username=var_username AND password=var_password);
	
	IF ret_person_exists THEN
		select id, persontypecode
		into ret_person_id, ret_person_type
		from person 
		where username=var_username AND password=var_password;
	ELSE
		ret_person_type := 'NULL';
	END if;
END;
$BODY$;

--==================================================================================================================
--CHECK if username already exists
CREATE OR REPLACE FUNCTION public.registerverifyuser(
	var_username character varying,
	OUT ret_username_exists boolean)
    RETURNS boolean
    LANGUAGE 'plpgsql'

    COST 100
    VOLATILE 
AS $BODY$
BEGIN
	
	SELECT exists (SELECT 1 FROM person WHERE username = var_username) INTO ret_username_exists;
	
END;
$BODY$;

--==================================================================================================================

-- Register a new person
CREATE OR REPLACE FUNCTION registerperson(
	var_name varchar(64),
	var_surname varchar(64),
	var_username varchar(64),
	var_password varchar(64),
	var_person_type varchar(6),
	OUT res_created boolean,
	OUT ret_username varchar(64),
	OUT ret_person_id uuid)
    LANGUAGE 'plpgsql'
AS $BODY$
DECLARE
	personid uuid := uuid_generate_v4();
BEGIN
    	
	INSERT INTO person(
	id, persontypecode, name, surname, username, password, createddatetime, isdeleted)
	VALUES (personid, var_person_type, var_name, var_surname, var_username, var_password, CURRENT_TIMESTAMP, false);
	
	res_created := true;
	ret_username := var_username;
	ret_person_id := personid;
END;
$BODY$;

--==================================================================================================================