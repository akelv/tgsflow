# Positive — Complex (While …, when …, the <System> shall …)
While the user is authenticated, when the token expires, the API shall return 401 and a renewal hint.

While a migration is in progress, when a write arrives, the Storage Service shall queue the write.

While maintenance mode is active, when a health check runs, the Probe shall return 503 with "maintenance".
