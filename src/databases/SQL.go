package databases

const INSERT_PASSWORD_ENTRY = `
WITH client_info as (
	SELECT id 
	FROM client
	WHERE id = $1 AND is_active = true
), inserted as (
    INSERT INTO siteinfo (name, type, metadata, username) 
	VALUES ($2, $3, $4, $5) 
	RETURNING id
)
INSERT into passwordentry(reference_id, client_id, site_id, length)
VALUES (
	$6,
	(SELECT id FROM client_info),
	(SELECT id FROM inserted),
	$7
)
RETURNING id`

const RETRIEVE_PASSWORD_ENTRY = `
WITH client_info as (
	SELECT id 
	FROM public.client
	WHERE id = $1 AND is_active = true
)
SELECT p.length, s.metadata
FROM public.passwordentry p
JOIN client_info AS c
ON c.id = p.client_id
JOIN public.siteinfo AS s
ON s.id = p.site_id
WHERE p.id = $2
`

const RETRIEVE_API_KEY = `
WITH client_info as (
	SELECT id 
	FROM client
	WHERE id = $1 AND is_active = true
)
SELECT k.id, k.api_key
FROM clientapikey k
JOIN client_info AS c
ON c.id = k.id
`
