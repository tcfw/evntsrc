CREATE TABLE event_store.events (
	id UUID PRIMARY KEY NOT NULL,
	stream INT NOT NULL,
	time TIMESTAMPTZ NOT NULL,
	type STRING NOT NULL,
	typeVersion STRING,
	ceVersion STRING,
	source STRING NOT NULL,
	subject STRING NOT NULL,
	acknowledged TIMESTAMP,
	metadata JSONB,
	contentType STRING,
	data BYTES,

	INDEX streamer (id, stream),
	INDEX ts (time)
)