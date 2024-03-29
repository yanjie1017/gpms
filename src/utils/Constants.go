package utils

const HEALTHCHECK_ENDPOINT = "/v1/health"
const PASSWORD_GENERATION_ENDPOINT = "/v1/password/generate"
const PASSWORD_ENTRY_CREATION_ENDPOINT = "/v1/entry/create"

const HTTP_ERROR_RESPONSE_PARSE_REQUEST = "Invalid request - Wrong field type"
const HTTP_ERROR_RESPONSE_MISSING_FIELD = "Missing field in request - "

const HTTP_ERROR_RESPONSE_GENERATE_RESPONSE = "Unable to generate response"

const HTTP_ERROR_RESPONSE_GENERATE_PASSWORD = "Unable to generate password"
const HTTP_ERROR_RESPONSE_PASSWORD_ENTRY_NOT_FOUND = "Password entry not found in the database"

const HTTP_ERROR_RESPONSE_PASSWORD_ENTRY_NOT_CREATED = "Unable to create password entry in database"
const HTTP_ERROR_RESPONSE_PASSWORD_LENGTH = "Invalid password length - Minumum length: 6"
