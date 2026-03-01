SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
CREATE SCHEMA osm2pgsql;
ALTER SCHEMA osm2pgsql OWNER TO osm2pgsql;
SET default_tablespace = '';
SET default_table_access_method = heap;
CREATE TABLE osm2pgsql.element (
    osm_type character(1) CONSTRAINT element_tmp_osm_type_not_null NOT NULL,
    osm_id bigint CONSTRAINT element_tmp_osm_id_not_null NOT NULL,
    name text,
    class text CONSTRAINT element_tmp_class_not_null NOT NULL,
    subclass text,
    tags jsonb CONSTRAINT element_tmp_tags_not_null NOT NULL,
    location public.geometry(Point,4326) CONSTRAINT element_tmp_location_not_null NOT NULL,
    lon double precision GENERATED ALWAYS AS (public.st_x(location)) STORED CONSTRAINT element_tmp_lon_not_null NOT NULL,
    lat double precision GENERATED ALWAYS AS (public.st_y(location)) STORED CONSTRAINT element_tmp_lat_not_null NOT NULL
);
ALTER TABLE osm2pgsql.element OWNER TO osm2pgsql;
CREATE TABLE osm2pgsql.osm2pgsql_properties (
    property text NOT NULL,
    value text NOT NULL
);
ALTER TABLE osm2pgsql.osm2pgsql_properties OWNER TO osm2pgsql;
ALTER TABLE ONLY osm2pgsql.osm2pgsql_properties
    ADD CONSTRAINT osm2pgsql_properties_pkey PRIMARY KEY (property);
CREATE INDEX element_location_idx ON osm2pgsql.element USING gist (location) WITH (fillfactor='100');
GRANT USAGE ON SCHEMA osm2pgsql TO socialmaps_api;
GRANT SELECT ON TABLE osm2pgsql.element TO socialmaps_api;
GRANT SELECT ON TABLE osm2pgsql.osm2pgsql_properties TO socialmaps_api;
ALTER DEFAULT PRIVILEGES FOR ROLE osm2pgsql IN SCHEMA osm2pgsql GRANT SELECT ON TABLES TO socialmaps_api;
