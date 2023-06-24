-- Database: 

-- DROP DATABASE marketplace;

CREATE DATABASE marketplace
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'C'
    LC_CTYPE = 'C'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

COMMENT ON DATABASE marketplace
    IS 'Database to store all order,products,customer data';