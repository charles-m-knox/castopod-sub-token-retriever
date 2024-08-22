# Simple example

This example connects to a Castopod mariadb database, and performs the necessary queries to update your database with new tokens.

## Usage

Start by cloning the repository:

```bash
git clone https://git.cmcode.dev/cmcode/castopod-sub-token-retriever.git
cd examples/server
cp config.example.json config.json
```

Edit `config.json` according to your needs. The configuration is defined [here](../../pkg/cstr/lib.go) in the `Config` struct. Comments for each configurable property are defined.

Proceed to build this application and run it:

```bash
go get -v
go build -v
# operate without modifying the database with the -test flag.
# also, if you want to skip sending emails and instead print the emailed output
# to the console, set the value of "smtp":{"testing":true} in your config.json.
./server -f config.json -test -addr "127.0.0.1:19281"
```

When you're ready to run the real thing, you can remove the `-test` flag.

Once it's up and running, proceed to test with `curl`:

```bash
# step 1: this will send a verification code to the user's email address
curl -X POST -v -sSL  -k 'http://127.0.0.1:19281/podcast-access/token?email=user@example.com&handle=podcast-handle'

# step 2: the user will receive a link in their email that looks like this:
curl -v -sSL  -k 'http://127.0.0.1:19281/podcast-access/token?email=user@example.com&handle=podcast-handle&code=fCmWlgCDJqvwU6Zz'

# step 3: the user will receive a final email that contains a link with a fresh
# token:
# http://podcast.example.com/@podcast-handle/feed.xml?token=d8dHcb2c
```

That's it!

## Tips for connecting to a remote mariadb db

If your mariadb database is only accessible behind an ssh tunnel, you can use ssh forwarding to open up the Castopod connection, assuming one is on 3306:

```bash
ssh -L 3306:127.0.0.1:3306 user@server
```

## Alternative: podman/docker container image

If you prefer not to build from source, you can use the pre-built container image.

You must first create a `config.json` just like above, and ensure that the output directory is going to exist within the container.

```bash
podman run --rm -it \
    -v "$(pwd)/config.json:/config.json:ro" \
    git.cmcode.dev/cmcode/castopod-sub-token-retriever:server-mariadb
```

Note: If you're using an SSH port forwarding mechanism for the mariadb database connection, you may want to consider adding `--network host` to the above `podman run` command.

### Building the container image

If you're building from an Arch Linux host, you can use your host system's pacman mirrorlist for faster builds. If not, remove the `-v` flag from the `podman build` command below. It is recommended to run `export GOSUMDB=off`.

```bash
podman build \
    -v "/etc/pacman.d/mirrorlist:/etc/pacman.d/mirrorlist:ro" \
    --build-arg GOSUMDB="${GOSUMDB}" \
    --build-arg GOPROXY="${GOPROXY}" \
    -f containerfile \
    -t git.cmcode.dev/cmcode/castopod-sub-token-retriever:server-mariadb .
```
