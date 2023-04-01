INSERT INTO client(
    name, 
    contact_number, 
    email_address, 
    mailing_address,
    postal_code,
    country,
    is_active
)
VALUES (
	'MockClient',
    '+6512345678',
    '36 College Ave E, North Tower, Singapore 138600',
    'mockclient@email.com',
    '138600',
    'SG',
    true
);

INSERT INTO clientapikey(
    client_id,
    api_key
)
VALUES (
    10000,
    'bd30cb14f01f2c0757a9dea306a2abcc14bc7e1aa0600290f180fdffef3d0fdb125174d0f14151b05d07c40282bff81b'
);