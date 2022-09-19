// The default amount of time (in milliseconds) to allow a request to run before
// terminating it early under assumption that the server is down or otherwise.
const REQUEST_TIMEOUT = 5000;

// API endpoint base for all unauthorized requests.
const UNAUTHORIZED_API_BASE = "http://localhost:8081";

// API endpoint base for all authorized requests.
const AUTHORIZED_API_BASE = `${UNAUTHORIZED_API_BASE}/largecurrency`;

export { REQUEST_TIMEOUT, UNAUTHORIZED_API_BASE, AUTHORIZED_API_BASE };
