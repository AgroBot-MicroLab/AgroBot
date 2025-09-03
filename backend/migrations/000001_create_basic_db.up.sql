-- Schema baseline for app

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET row_security = off;

SET default_tablespace = '';
SET default_table_access_method = heap;

--
-- Table: drone_position
--
CREATE TABLE public.drone_position (
    id integer NOT NULL,
    lat double precision,
    long double precision,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE SEQUENCE public.dron_position_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.dron_position_id_seq OWNED BY public.drone_position.id;

ALTER TABLE ONLY public.drone_position
    ALTER COLUMN id SET DEFAULT nextval('public.dron_position_id_seq'::regclass);

ALTER TABLE ONLY public.drone_position
    ADD CONSTRAINT drone_position_pk PRIMARY KEY (id);

--
-- Table: healthprobe
--
CREATE TABLE public.healthprobe (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE SEQUENCE public.healthprobe_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.healthprobe_id_seq OWNED BY public.healthprobe.id;

ALTER TABLE ONLY public.healthprobe
    ALTER COLUMN id SET DEFAULT nextval('public.healthprobe_id_seq'::regclass);

ALTER TABLE ONLY public.healthprobe
    ADD CONSTRAINT healthprobe_pkey PRIMARY KEY (id);

--
-- Table: images
--
CREATE TABLE public.images (
    id integer NOT NULL,
    path character varying,
    captured_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE SEQUENCE public.images_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.images_id_seq OWNED BY public.images.id;

ALTER TABLE ONLY public.images
    ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pk PRIMARY KEY (id);

--
-- Table: point
--
CREATE TABLE public.point (
    id integer NOT NULL,
    lat double precision,
    long double precision,
    status character varying,
    image_id integer,
    started_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    finished_at timestamp without time zone
);

CREATE SEQUENCE public.point_id_seq
    AS integer START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

ALTER SEQUENCE public.point_id_seq OWNED BY public.point.id;

ALTER TABLE ONLY public.point
    ALTER COLUMN id SET DEFAULT nextval('public.point_id_seq'::regclass);

ALTER TABLE ONLY public.point
    ADD CONSTRAINT point_pk PRIMARY KEY (id);

