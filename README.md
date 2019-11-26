# withenv: runs processes with environment settings from a file

It's just like envdir, but for go, and it reads files.

## Goals

When I looked around, I couldn't find a tool for a small need: to run an application quickly with appropriate environment variables set, when no shell is available.

This particular use case helps running handlers in [sensu-go](https://sensu.io/), while secrets are not exposed to the API / web UI. Information about a file on a system is less threatening than seeing the actual secret set by an environment variable, which is not even redacted by the user interface.

## Legal

This project is licensed under [Blue Oak Model License v1.0.0](https://blueoakcouncil.org/license/1.0.0). It is not registered either at OSI or GNU, therefore GitHub is widely looking at the other direction. However, this is the license I'm most happy with: you can read and understand it with no legal degree, and there are no hidden or cryptic meanings in it.

The project is also governed with [Contributor Covenant](https://contributor-covenant.org/)'s [Code of Conduct](https://www.contributor-covenant.org/version/1/4/) in mind. I'm not copying it here, as a pledge for taking the verbatim version by the word, and we are not going to modify it in any way.

## Any issues?

Open a ticket, perhaps a pull request. We support [GitHub Flow](https://guides.github.com/introduction/flow/). You might want to [fork](https://guides.github.com/activities/forking/) this project first.
