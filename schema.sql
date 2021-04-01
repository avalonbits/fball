CREATE TABLE IF NOT EXISTS RequestCache(
    Endpoint TEXT NOT NULL,
    Params TEXT NOT NULL,
    Timestamp INTEGER NOT NULL,
    Response BLOB NOT NULL,

    PRIMARY KEY (Endpoint, Params, Timestamp)
);
CREATE INDEX EndpointIDX ON RequestCache(Endpoint);
CREATE INDEX OrderedEndpointIDX ON RequestCache(Endpoint, Timestamp);
