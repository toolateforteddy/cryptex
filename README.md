# Cryptex
Status: Work In Progress

# Goals
Storing and sharing secrets in a git project is hard. The goal of this project is to allow you to store your secrets securely inside of git, encrypted with a key. Secrets are stored in a toml file to allow for structured data.

# Usage Pattern
In the root, run `cryptex init` to create the empty secret storage. Then run `cryptex edit` to open the file and add your secrets.

When you later want to use your secrets, you can flatten the tree with `cryptex export` which will print out your secrets in the KEY="VALUE" form that bash expects.

# TODO
- Build `cryptex init`
- - Add relevant files to .gitignore if in a git repo.
- - If password is not provided create one.
- Build `cryptex keygen`
- Expand Key Loading/Management
- - Load from Environment if set.
- - Save to stored_key if loaded from environment.
- - Expose flag to pass in key as arg.
- - Explore alternative crypto pkgs and algorithms.
- - Should key be restructured to be 32 random bytes, then base64 encoded so that it is easier to work with? (Hint, Yes)
- More Tests
- More Documentation
- Example Secret Files
- Documentation for running in Docker
- - Maybe also a dockerfile?
- - Maybe also a build script?


