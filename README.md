# Wait Until

Command line utility which holds up execution of a bash pipeline until a given command returns a desired exit code.

Useful for preventing race conditions in build pipelines.

## Examples

Wait until the database container is accepting connections before importing the SQL dump.

```bash
docker-compose up -d
wait-until --timeout=60s -- mysqladmin ping -h 127.0.0.1
mysql < db.sql
```

Wait for redis to accept connections before starting the app.

```bash
wait-until --retries=5 --sleep=5s -- redis-cli ping
./app start
```

## Usage

```bash
usage: wait-until [<flags>] <command>

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose          Enabled verbose output.
  -t, --timeout=TIMEOUT  Timeout before aborting pipeline. Omit for no limit.
  -r, --retries=-1       Number of attempts before aborting pipeline. -1 for no limit.
  -s, --sleep=1s         Sleep time between each execution.
  -e, --exit-code=0      Desired exit code before allowing pipeline to proceed.

Args:
  <command>  Command to repeatedly execute until exit code met, timeout exceeded, or retry limit exceeded.
```

## Notes

It is recommended to use the `+pipefail` bash directive when using this tool to ensure failures in this command terminate the pipeline.  