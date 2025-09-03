--
-- PostgreSQL database dump
--

\restrict 0vffcobEg1gJSBjwRCa3oodaoAZc7NuUu0CfyPBB8mIRbns0NoxshLbujwJrATC

-- Dumped from database version 16.10
-- Dumped by pg_dump version 16.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: drone_position; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.drone_position (
    id integer NOT NULL,
    lat double precision,
    long double precision,
    "timestamp" timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.drone_position OWNER TO app;

--
-- Name: dron_position_id_seq; Type: SEQUENCE; Schema: public; Owner: app
--

CREATE SEQUENCE public.dron_position_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.dron_position_id_seq OWNER TO app;

--
-- Name: dron_position_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: app
--

ALTER SEQUENCE public.dron_position_id_seq OWNED BY public.drone_position.id;


--
-- Name: healthprobe; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.healthprobe (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.healthprobe OWNER TO app;

--
-- Name: healthprobe_id_seq; Type: SEQUENCE; Schema: public; Owner: app
--

CREATE SEQUENCE public.healthprobe_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.healthprobe_id_seq OWNER TO app;

--
-- Name: healthprobe_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: app
--

ALTER SEQUENCE public.healthprobe_id_seq OWNED BY public.healthprobe.id;


--
-- Name: images; Type: TABLE; Schema: public; Owner: app
--

CREATE TABLE public.images (
    id integer NOT NULL,
    path character varying,
    captured_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.images OWNER TO app;

--
-- Name: images_id_seq; Type: SEQUENCE; Schema: public; Owner: app
--

CREATE SEQUENCE public.images_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.images_id_seq OWNER TO app;

--
-- Name: images_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: app
--

ALTER SEQUENCE public.images_id_seq OWNED BY public.images.id;


--
-- Name: point; Type: TABLE; Schema: public; Owner: app
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


ALTER TABLE public.point OWNER TO app;

--
-- Name: point_id_seq; Type: SEQUENCE; Schema: public; Owner: app
--

CREATE SEQUENCE public.point_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.point_id_seq OWNER TO app;

--
-- Name: point_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: app
--

ALTER SEQUENCE public.point_id_seq OWNED BY public.point.id;


--
-- Name: drone_position id; Type: DEFAULT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.drone_position ALTER COLUMN id SET DEFAULT nextval('public.dron_position_id_seq'::regclass);


--
-- Name: healthprobe id; Type: DEFAULT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.healthprobe ALTER COLUMN id SET DEFAULT nextval('public.healthprobe_id_seq'::regclass);


--
-- Name: images id; Type: DEFAULT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.images ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);


--
-- Name: point id; Type: DEFAULT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.point ALTER COLUMN id SET DEFAULT nextval('public.point_id_seq'::regclass);


--
-- Data for Name: drone_position; Type: TABLE DATA; Schema: public; Owner: app
--

COPY public.drone_position (id, lat, long, "timestamp") FROM stdin;
\.


--
-- Data for Name: healthprobe; Type: TABLE DATA; Schema: public; Owner: app
--

COPY public.healthprobe (id, created_at) FROM stdin;
\.


--
-- Data for Name: images; Type: TABLE DATA; Schema: public; Owner: app
--

COPY public.images (id, path, captured_at) FROM stdin;
1	test/test	2025-09-03 15:22:35.286199
\.


--
-- Data for Name: point; Type: TABLE DATA; Schema: public; Owner: app
--

COPY public.point (id, lat, long, status, image_id, started_at, finished_at) FROM stdin;
1	23	33	pending	1	2025-09-03 15:18:56.582664	\N
\.


--
-- Name: dron_position_id_seq; Type: SEQUENCE SET; Schema: public; Owner: app
--

SELECT pg_catalog.setval('public.dron_position_id_seq', 1, false);


--
-- Name: healthprobe_id_seq; Type: SEQUENCE SET; Schema: public; Owner: app
--

SELECT pg_catalog.setval('public.healthprobe_id_seq', 1, false);


--
-- Name: images_id_seq; Type: SEQUENCE SET; Schema: public; Owner: app
--

SELECT pg_catalog.setval('public.images_id_seq', 1, true);


--
-- Name: point_id_seq; Type: SEQUENCE SET; Schema: public; Owner: app
--

SELECT pg_catalog.setval('public.point_id_seq', 1, true);


--
-- Name: drone_position drone_position_pk; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.drone_position
    ADD CONSTRAINT drone_position_pk PRIMARY KEY (id);


--
-- Name: healthprobe healthprobe_pkey; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.healthprobe
    ADD CONSTRAINT healthprobe_pkey PRIMARY KEY (id);


--
-- Name: images images_pk; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_pk PRIMARY KEY (id);


--
-- Name: point point_pk; Type: CONSTRAINT; Schema: public; Owner: app
--

ALTER TABLE ONLY public.point
    ADD CONSTRAINT point_pk PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

\unrestrict 0vffcobEg1gJSBjwRCa3oodaoAZc7NuUu0CfyPBB8mIRbns0NoxshLbujwJrATC

