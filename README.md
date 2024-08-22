# castopod-sub-token-retriever

Exposes a simple REST API endpoint that allows users in a Castopod database to reset their private podcast feed URL.

This leverages a database connection - Castopod recommends mariadb as its supported database type, but it can technically work with any compliant driver that implements Go's `database/sql` interfaces (you'll need to make some tweaks to the sql parameterization tokens such as `?` or `$1 $2` in the example server code)

This is compatible with **Castopod v1.12.3**. Compatibility with any other version is not guaranteed.

## Example Usage

See [`examples/server/README.md`](./examples/server/README.md).

## Warnings and limitations

First and foremost, you must manage the Redis cache separately from this application. Castopod stores data in Redis (if available), and updates to the database will not be reflected. One easy way to manage this is to purge the Redis cache frequently whenever this application runs:

```bash
export REDIS_PASSWORD=your_password_goes_here
redis-cli FLUSHALL
```

If you're running it in a container, do the following *immediately* after using this library:

```bash
podman exec -it castopod-redis redis-cli -a your_password_goes_here FLUSHALL
```

Additionally, **Please verify the output of this application before using it**. Stability is never guaranteed. Make backups.

## Structure

In order to keep dependencies to zero (see [`go.mod`](./go.mod) - it only uses the standard library), this Go module is structured as a library that can be imported by any application.

This has the benefit of not requiring any specific SQL driver - it accepts SQL rows themselves from the `database/sql` package. You can use any compliant driver, such as sqlite or postgres - although I haven't tried anything aside from mariadb, so tread carefully.

## Development notes

At this time, unit tests are not implemented for this module or example server. I structured it with a `pkg/` directory but may end up making further changes to this as time goes on, if needed.

## Compatibility

Here is a list of all tested Castopod versions that are known to work:

- v1.12.3

If any database migrations that modify the tables used by this library occur upstream, this application will likely break.
