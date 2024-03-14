--
-- PostgreSQL database dump
--

-- Dumped from database version 14.11 (Homebrew)
-- Dumped by pg_dump version 14.11 (Homebrew)

-- Started on 2024-03-14 15:50:58 MDT

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

--
-- TOC entry 3650 (class 1262 OID 14088)
-- Name: postgres; Type: DATABASE; Schema: -; Owner: hess
--

CREATE DATABASE postgres WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE = 'C';


ALTER DATABASE postgres OWNER TO hess;

\connect postgres

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

--
-- TOC entry 3651 (class 0 OID 0)
-- Dependencies: 3650
-- Name: DATABASE postgres; Type: COMMENT; Schema: -; Owner: hess
--

COMMENT ON DATABASE postgres IS 'default administrative connection database';


--
-- TOC entry 3 (class 2615 OID 2200)
-- Name: public; Type: SCHEMA; Schema: -; Owner: hess
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO hess;

--
-- TOC entry 3652 (class 0 OID 0)
-- Dependencies: 3
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: hess
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- TOC entry 211 (class 1255 OID 16399)
-- Name: mynotify(); Type: FUNCTION; Schema: public; Owner: hess
--

CREATE FUNCTION public.mynotify() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
	row RECORD;
	output text;

begin 
	if (TG_OP = 'DELETE') then
		row = old;
	else
		row = new;
	end if;

	output = 'OPERATION= ' || TG_OP || ', ID= ' || row.id;

	perform pg_notify('ithappened', output);

	return null;
end; $$;


ALTER FUNCTION public.mynotify() OWNER TO hess;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 209 (class 1259 OID 16384)
-- Name: todos; Type: TABLE; Schema: public; Owner: hess
--

CREATE TABLE public.todos (
    item character varying,
    done boolean DEFAULT false,
    id integer NOT NULL
);


ALTER TABLE public.todos OWNER TO hess;

--
-- TOC entry 210 (class 1259 OID 16390)
-- Name: todos_id_seq; Type: SEQUENCE; Schema: public; Owner: hess
--

CREATE SEQUENCE public.todos_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.todos_id_seq OWNER TO hess;

--
-- TOC entry 3653 (class 0 OID 0)
-- Dependencies: 210
-- Name: todos_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: hess
--

ALTER SEQUENCE public.todos_id_seq OWNED BY public.todos.id;


--
-- TOC entry 3500 (class 2604 OID 16391)
-- Name: todos id; Type: DEFAULT; Schema: public; Owner: hess
--

ALTER TABLE ONLY public.todos ALTER COLUMN id SET DEFAULT nextval('public.todos_id_seq'::regclass);


--
-- TOC entry 3643 (class 0 OID 16384)
-- Dependencies: 209
-- Data for Name: todos; Type: TABLE DATA; Schema: public; Owner: hess
--

INSERT INTO public.todos VALUES ('Make POST and PUT endpoints', true, 1);
INSERT INTO public.todos VALUES ('try for update', false, 2);
INSERT INTO public.todos VALUES ('try for update', false, 3);
INSERT INTO public.todos VALUES ('insert test', false, 4);
INSERT INTO public.todos VALUES ('insert test 2', false, 5);


--
-- TOC entry 3654 (class 0 OID 0)
-- Dependencies: 210
-- Name: todos_id_seq; Type: SEQUENCE SET; Schema: public; Owner: hess
--

SELECT pg_catalog.setval('public.todos_id_seq', 5, true);


--
-- TOC entry 3502 (class 2606 OID 16398)
-- Name: todos todos_pk; Type: CONSTRAINT; Schema: public; Owner: hess
--

ALTER TABLE ONLY public.todos
    ADD CONSTRAINT todos_pk PRIMARY KEY (id);


--
-- TOC entry 3503 (class 2620 OID 16400)
-- Name: todos trigger_my_notify; Type: TRIGGER; Schema: public; Owner: hess
--

CREATE TRIGGER trigger_my_notify AFTER INSERT OR DELETE OR UPDATE ON public.todos FOR EACH ROW EXECUTE FUNCTION public.mynotify();


-- Completed on 2024-03-14 15:50:58 MDT

--
-- PostgreSQL database dump complete
--

