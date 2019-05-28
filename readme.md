# SA:MP Forum Discord Bot

[![All Contributors](https://img.shields.io/badge/all_contributors-14-orange.svg?style=flat-square)](#contributors)

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

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore -->
<table><tr><td align="center"><a href="https://wopss.net"><img src="https://avatars3.githubusercontent.com/u/3403191?v=4" width="100px;" alt="Octavian Dima"/><br /><sub><b>Octavian Dima</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=WopsS" title="Code">💻</a> <a href="#ideas-WopsS" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Southclaws/cj/issues?q=author%3AWopsS" title="Bug reports">🐛</a></td><td align="center"><a href="https://github.com/J0shES"><img src="https://avatars0.githubusercontent.com/u/18373054?v=4" width="100px;" alt="J0shES"/><br /><sub><b>J0shES</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=J0shES" title="Code">💻</a> <a href="#ideas-J0shES" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/Southclaws/cj/issues?q=author%3AJ0shES" title="Bug reports">🐛</a> <a href="#maintenance-J0shES" title="Maintenance">🚧</a></td><td align="center"><a href="https://github.com/Dayvison"><img src="https://avatars0.githubusercontent.com/u/10089094?v=4" width="100px;" alt="Dayvison"/><br /><sub><b>Dayvison</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=Dayvison" title="Code">💻</a></td><td align="center"><a href="https://adriangraber.com"><img src="https://avatars1.githubusercontent.com/u/18301034?v=4" width="100px;" alt="Adrian Graber"/><br /><sub><b>Adrian Graber</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=AGraber" title="Code">💻</a></td><td align="center"><a href="https://github.com/Sreyas-Sreelal"><img src="https://avatars3.githubusercontent.com/u/17766494?v=4" width="100px;" alt="__SyS__"/><br /><sub><b>__SyS__</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=Sreyas-Sreelal" title="Code">💻</a></td><td align="center"><a href="https://gigabitz.pw"><img src="https://avatars3.githubusercontent.com/u/15860096?v=4" width="100px;" alt="Robster"/><br /><sub><b>Robster</b></sub></a><br /><a href="#content-Gigabitzz" title="Content">🖋</a></td><td align="center"><a href="https://twitter.com/dakyskye"><img src="https://avatars1.githubusercontent.com/u/32128756?v=4" width="100px;" alt="Lasha Kanteladze"/><br /><sub><b>Lasha Kanteladze</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=dakyskye" title="Code">💻</a> <a href="https://github.com/Southclaws/cj/commits?author=dakyskye" title="Tests">⚠️</a> <a href="#ideas-dakyskye" title="Ideas, Planning, & Feedback">🤔</a></td></tr><tr><td align="center"><a href="https://kristo.xyz"><img src="https://avatars3.githubusercontent.com/u/7974602?v=4" width="100px;" alt="Kristo Isberg"/><br /><sub><b>Kristo Isberg</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=kristoisberg" title="Code">💻</a> <a href="https://github.com/Southclaws/cj/commits?author=kristoisberg" title="Tests">⚠️</a> <a href="#ideas-kristoisberg" title="Ideas, Planning, & Feedback">🤔</a></td><td align="center"><a href="https://marcelschr.me"><img src="https://avatars3.githubusercontent.com/u/19377618?v=4" width="100px;" alt="Marcel Schramm"/><br /><sub><b>Marcel Schramm</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=Bios-Marcel" title="Code">💻</a></td><td align="center"><a href="https://github.com/thecodeah"><img src="https://avatars0.githubusercontent.com/u/21268739?v=4" width="100px;" alt="Codeah"/><br /><sub><b>Codeah</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=thecodeah" title="Code">💻</a></td><td align="center"><a href="https://github.com/GiampaoloFalqui"><img src="https://avatars3.githubusercontent.com/u/4460702?v=4" width="100px;" alt="Giampaolo Falqui"/><br /><sub><b>Giampaolo Falqui</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=GiampaoloFalqui" title="Documentation">📖</a></td><td align="center"><a href="https://github.com/Sunehildeep"><img src="https://avatars1.githubusercontent.com/u/23412507?v=4" width="100px;" alt="Sunehildeep"/><br /><sub><b>Sunehildeep</b></sub></a><br /><a href="https://github.com/Southclaws/cj/commits?author=Sunehildeep" title="Code">💻</a></td><td align="center"><a href="https://redcountyrp.com"><img src="https://avatars0.githubusercontent.com/u/5786576?v=4" width="100px;" alt="TommyB"/><br /><sub><b>TommyB</b></sub></a><br /><a href="#content-TommyB123" title="Content">🖋</a></td><td align="center"><a href="https://github.com/Hual"><img src="https://avatars0.githubusercontent.com/u/1867646?v=4" width="100px;" alt="Nikola Yanakiev"/><br /><sub><b>Nikola Yanakiev</b></sub></a><br /><a href="#content-Hual" title="Content">🖋</a></td></tr></table>

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification.
Contributions of any kind welcome!
