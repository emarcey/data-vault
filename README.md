# data_vault

## Table of Contents
<!-- MarkdownTOC -->

- [Storage](#storage)
- [Access](#access)
- [API](#api)
	- [Authentication](#authentication)
		- [Access Token](#access-token)
		- [Client Secret](#client-secret)
	- [Pagination](#pagination)
	- [Users](#users)
	- [User Groups](#user-groups)
	- [Secrets](#secrets)
	- [Secret Permissions](#secret-permissions)
- [Roadmap](#roadmap)
- [Components](#components)
- [Configuration](#configuration)
- [Development](#development)

<!-- /MarkdownTOC -->


A personal project for the storage of sensitive key-value pairs.

*Version:* v0.0.1

## Storage

The keys and values are stored in separated locations.

First, the key name and metadata is stored in a relational database (Postgres implementation provided). The value is encrypted using [AES-256](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard) and is paired with a randomly generated UUID and stored alongside metadata.

The Initialization Vector (IV) and the Encryption Key are stored in a separate datastore (MongoDB implementation provided), indexed by the UUID.

At decrypt time, the IV and Key are fetched from the separate datastore, then the value is decrypted in memory and returned to the caller.

## Access

Access is provisioned according to users, both standard and developer. For all interactions with the API, a user must first generate an access token which lasts 24 hours.

Developer users have the ability to:

1. Create a key-value pair
1. Fetch a key-value pair for keys they have created or to which they have been granted access
1. Grant/Revoke permissions on keys they have created to other users
1. List users/user groups

Admins have extended permissions. In addition to create/fetch, they have the ability to

1. Delete a key-value pair
1. Grant/Revoke permissions on any keys
1. Create/Delete users
1. Create/Delete user groups & add/remove users to/from groups


All API interactions with the key-value store (fetch, create, delete) are additionally logged in the secrets datastore (MongoDB implementation provided). The MongoDB implementation is structured for a time-series collection keyed on user ID.

## API

### Authentication

#### Access Token

On creation, a user will be granted a client ID/secret pair. This pair will be used to identify the user, and is necessary for generating an Access Token, which is used to access all other endpoints.

To create an access token, a user must call the following endpoint:

`GET: {base_url}/access_token`

With the headers:

```json
{
	"Client-Id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	"Client-Secret": "9e5d5da6-24ed-42a9-a105-c531bed8175d"
}
```

On success, this returns:

```json
{
    "id": "92b2198a-a0b6-4f37-ba53-955447604b31",
    "client_id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
    "invalid_at": "2022-03-24T11:18:58.911523-04:00",
    "is_latest": true
}
```

where `id` is the access token to be used for future requests.

In all subsequent requests, API endpoints should be queried with the header `Access-Token` set to the value of this id.

#### Client Secret

Should a user lose their secret, or if they believe it has been compromised, a user can rotate their secret, which is used to generate an access token.

This is done with the call:

`GET: {base_url}/rotate`

With the headers:

```json
{
	"Client-Id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	"Client-Secret": "9e5d5da6-24ed-42a9-a105-c531bed8175d"
}
```

On success, this returns:

```json
{
	"Client-Id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	"Client-Secret": "29a52d35-d8d5-4ead-ac4a-90dba908ecaa"
}
```

### Pagination

All `List` endpoints support limit/offset pagination.

Call the endpoints with the URL Query Parameters `pageSize` and `offset`, like:

`GET {base_url}/users?pageSize=10&offset=0`

If not set, page size will default to 10 and offset will default to 0.

### Users

**Note: All User Endpoints except List are currently Admin-Only**

1. List
	* Method: GET
	* URI: `/users`
	* Response: List of User objects
		```json
		[
			{
	        	"id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
		        "name": "admin",
		        "is_active": true,
		        "type": "admin"
	    	}
	    ]
		```
1. Get
	* Method: GET
	* URI: `/users/{userId}`
	* Response: Single User object
		```json
		{
	        "id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	        "name": "admin",
	        "is_active": true,
	        "type": "admin"
	    }
		```
1. Create
	* Method: POST
	* URI: `/users`
	* Request:
		```json
		{
	        "name": "new user1",
	        "type": "developer"
	    }
		```
	* Response: Single User object
		```json
		{
	        "id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	        "name": "new user1",
	        "is_active": true,
	        "type": "developer"
	    }
		```
1. Delete
	* Method: DELETE
	* URI: `/users/{userId}`
	* Response: None, if successful
	* Note: Delete is soft delete, so record will be inaccessible, but not deleted from the database entirely.

### User Groups

**Note: All User Group Endpoints except List are currently Admin-Only**

1. List
	* Method: GET
	* URI: `/user-groups`
	* Response: List of User Group objects
		```json
		[
			{
	        	"id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
		        "name": "admin"
	    	}
	    ]
		```
1. Get
	* Method: GET
	* URI: `/user-groups/{userGroupId}`
	* Response: Single User Group object
		```json
		{
	        "id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	        "name": "admin"
	    }
		```
1. Create
	* Method: POST
	* URI: `/user-groups`
	* Request:
		```json
		{
	        "name": "new user group1"
	    }
		```
	* Response: Single User Group object
		```json
		{
	        "id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
	        "name": "new user group1"
	    }
		```
1. Delete
	* Method: DELETE
	* URI: `/user-groups/{userGroupId}`
	* Response: None, if successful
	* Note: Delete is soft delete, so record will be inaccessible, but not deleted from the database entirely.
1. List Users in Group
	* Method: GET
	* URI: `user-groups/{userGroupId}/users`
	* Response: List of User objects
		```json
		[
			{
	        	"id": "03b6f72c-f3f4-43d9-a705-17b326924d74",
		        "name": "admin",
		        "is_active": true,
		        "type": "admin"
	    	}
	    ]
		```
1. Add Users to Group
	* Method: POST
	* URI: `user-groups/{userGroupId}/users`
	* Request:
		```json
		{
			"user_id": "03b6f72c-f3f4-43d9-a705-17b326924d74"
		}
		```
	* Response: None, if successful
1. Remove Users from Group
	* Method: DELETE
	* URI: `user-groups/{userGroupId}/users`
	* Request:
		```json
		{
			"user_id": "03b6f72c-f3f4-43d9-a705-17b326924d74"
		}
		```
	* Response: None, if successful
	* Note: Delete is soft delete, so record will be inaccessible, but not deleted from the database entirely.

### Secrets

**Note: Secret operations are performed against secret name rather than ID, as storing a separate secret ID in someone else's DB just seems like a waste of energy**

1. List
	* Method: GET
	* URI: `/secrets`
	* Response: list of secrets; value will not be set
		```json
		[
			{
			    "id": "c13dc88b-9563-43d8-bb70-81cb7f5af675",
			    "name": "my-key4",
			    "description": "something",
			    "created_by": "admin",
			    "updated_by": "admin"
			}
		]
		```
1. Get
	* Method: GET
	* URI: `/secrets/{secretName}`
	* Response: Decrypted secret
		```json
		{
		    "id": "c13dc88b-9563-43d8-bb70-81cb7f5af675",
		    "name": "my-key4",
		    "value": "doy2 ",
		    "description": "something",
		    "created_by": "admin",
		    "updated_by": "admin"
		}
		```
1. Create
	* Method: POST
	* URI: `/secrets`
	* Request:
		```json
		{
		    "name": "my-key4",
		    "value": "doy2 ",
		    "description": "something"
		}
		```
	* Response: Decrypted secret
		```json
		{
		    "id": "c13dc88b-9563-43d8-bb70-81cb7f5af675",
		    "name": "my-key4",
		    "value": "doy2 ",
		    "description": "something",
		    "created_by": "admin",
		    "updated_by": "admin"
		}
		```
1. Delete
	Method: DELETE
	* URI: `/secrets/{secretName}`
	* Response: None, if successful
	* Note: Delete is soft delete, so record will be inaccessible, but not deleted from the database entirely.
	* Note: endpoint is admin only


### Secret Permissions

Used to add read permissions for a user or group.

**Note: if both user_id and user_group_id are set in the request, will return an error**

1. Create
	* Method: POST
	* URI: `/secrets/{secretName}/permissions`
	* Request:
		```json
		{
			"user_id": "c13dc88b-9563-43d8-bb70-81cb7f5af675",
			"user_group_id": "c13dc88b-9563-43d8-bb70-81cb7f5af675"
		}
		```
	* Response: None, if successful
1. Delete
	* Method: DELETE
	* URI: `/secrets/{secretName}/permissions`
	* Request:
		```json
		{
			"user_id": "c13dc88b-9563-43d8-bb70-81cb7f5af675",
			"user_group_id": "c13dc88b-9563-43d8-bb70-81cb7f5af675"
		}
		```
	* Response: None, if successful


## Roadmap

* Improved permissioning
	* Wildcard-based access
* Extended support for interfaces
	* Datadog for tracing support
	* Other data stores for secret manager/db
* Better dev tools
* Automated tests

## Components

* Core database: stores user/permission data and the encrypted values of the secrets
	* Currently supported:
		* [Postgres](postgresql.org)
* Secrets database: stores access logs and encryption keys for secrets
	* Currently supported:
		* [MongoDB](mongodb.com)
* Logger: agent that provides logging for server
	* Currently supported (via [logrus](https://github.com/sirupsen/logrus):
		* basic text logger
		* JSON logger
* Tracer: agent that provides tracing for datastore/API access
	* Currently supported:
		* No Op. Literally does nothing
		* Local tracer (basically just a logger)
		* [Sentry.io](sentry.io)
* Service:
	* [Golang](golang.org) service
		* Endpoints supported with [go-kit](https://github.com/go-kit/kit) and [gorilla mux](https://github.com/gorilla/mux)


## Configuration

Server configuration is done using the `server_conf.yml` file.

This file allows the executor to configure the address, environment and other server configs (e.g. how long should access tokens last).

In addition, this is where connection settings for data stores and other dependencies are configured.

## Development

* Start up a postgres cluster
	* Run the contents `scripts/ddl.sql`
	* Manually create a new admin user with self-generated client ID/secret (Note: the secret in the db will be `sha256:{sha256 hash of client secret}`)
* Start up a MongoDB cluster
	* Create a database with collections for accessLogs and for secrets
* Copy `server_conf.example.yml` to `server_conf.yml`
	* Update `server_conf.yml` with postgres and MongoDB settings.
	* Adjust any other settings as needed
* Run `go mod vendor` to install vendor packages
* Run `go run main.go` to start the server