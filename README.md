TODO:
- Add book search / including test handler and service
- Proper readme
- Remove pq.StringArray from users model
- Docker ping
- Move migrate and seed to own cmd
- Run clean for sum and mod


# Book API

## Setup

Running with Docker will launch the app and database containers. The required database will be created and populated with test data, including an admin with the username `admin` and the password `Adminl#1`. These admin credentials can be passed as the Basic Auth Header to the `/users/token` endpoint in order to retrieve the login token. The login token should then be used as the Bearer Token for all other endpoints.