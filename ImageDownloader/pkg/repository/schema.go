package repository

var UrlSchema = `
CREATE TABLE IF NOT EXISTS images (
	url TEXT PRIMARY KEY,
	content BYTEA
)
`
