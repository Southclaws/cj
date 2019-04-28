# SA:MP Forum Discord Bot

[![Travis](https://img.shields.io/travis/Southclaws/cj.svg)](https://travis-ci.org/Southclaws/cj)

![CJ](cj.png)

CJ verifies Burgershot (formerly SA:MP) forum accounts and performs other tasks if you ask nicely.

Not really in development and not accepting new features. I fix bugs from time to time but it serves it purpose as a
verification tool and basic forum interface.

## Development

This project is open to anyone who wants to contribute, large or small! Whether you noticed a typo or want to add a
whole new feature, go for it!

Large additions should be discussed in issues or on Discord first. If you're new to Golang, ask me on Discord for where
to start and you can use CJ as a starting point for a contribution.

### Testing/Workflow

To run the app, you need:

- A Discord server to test - you can't use the SA:MP Discord to do tests
- Go 1.11 - no guarantees on older versions

If you don't own/admin a Discord server, creating one is simple, you can do it from the same menu you join discord
servers from.

#### Running with a database

If you want to develop features that require persisting data, spin up a MongoDB database. If you have Docker installed,
this is as simple as running `make mongodb` which will start a MongoDB container with a user `root` that has no
password. If you don't have Docker, you'll need to
[install MongoDB onto your system.](https://docs.mongodb.com/manual/installation/).

#### Running without a database

If you don't need a database for your feature, just add `NO_DATABSE=true` to `.env`.

Finally, the application requires some configuration. Copy the `example.env` to `.env` and modify it to use your token
and various IDs. Depending on what you're working on, some values won't be necessary. For example, unless you're
actually working on the verification system, you don't need to set the verified role ID.

Now you can build and run the application with `make local`.
