# session-auth

This was originally intended to be an authentication microservice for the Sessions application that used a NO-SQL DB for storing auth credentials with a direct ID acting as a tunnel to the main DB. However, the approach ended up being overly complex, making the setup more complicated than necessary. With Sessions now operating under a monolithic API structure, a separate authentication service is no longer needed.
