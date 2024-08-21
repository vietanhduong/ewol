# eWoL

Extended Wake-on-LAN. This repository is inspired by [wakeonlan](https://github.com/jpoliv/wakeonlan) but provides additional features like publishing the service as an HTTP server (`serve` mode) instead of running the straight `wakeonlan` command.

## Usage

`eWoL` provides two main features: `wake-on-lan` as the default command and the `remote` command.

### Wake-on-LAN

```console
$ ewol --help
Extended Wake-on-LAN is a tool to wake up devices on a local network.
You can also publish the service to the network and wake up the input device remotely via an API call.

Usage:
  ewol HARDWARE_ADDRESS [flags]
  ewol [command]

Examples:
# Wake-on-LAN directly
$ ewol 00:11:22:33:44:55

# Publish the service to the network and wake up the device remotely
# You can also provide the wake secret via the WAKE_SECRET environment variable
$ ewol 00:11:22:33:44:55 --serve --wake.secret mysecret

# Specify the IP address and port
$ ewol 00:11:22:33:44:55 --wake.ip 255.255.255.255 --wake.port 9

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  remote      Remote Wake on LAN for an eWoL server

Flags:
  -h, --help                            Help for ewol
      --log.format string               Log format. Available options: text, json (default "text")
      --log.level string                Log level (default "info")
      --serve                           Enable serve mode. This will create an HTTP server to listen for incoming requests
      --server.address string           Server listen address (default "0.0.0.0:8080")
      --server.drain-timeout duration   Server drain timeout (default 15s)
  -v, --version                         Print version and exit
  -i, --wake.ip string                  Destination IP address. Unless you have static ARP tables, you should use some kind of broadcast address (the broadcast address of the network where the computer resides or the limited broadcast address) (default "255.255.255.255")
  -p, --wake.port uint16                Destination port (default 9)
  -s, --wake.secret string              Secret key which will be used as a simple auth. Only works if you enable serve mode

Use "ewol [command] --help" for more information about a command.
```

If the `--serve` flag is specified, this will run as an HTTP server and serve incoming requests.

Currently, I expose only two APIs: `/wake` to wake up the input `HARDWARE_ADDRESS` and `/healthz` for health check purposes.

Please note that the `/wake` API only accepts the `HTTP POST` method.

### Remote Wake
```console
$ ewol remote --help
NOTE: This command is ONLY available for an eWoL server

Usage:
  ewol remote REMOTE_ADDRESS [flags]

Examples:
# Simple remote wake up
$ ewol remote http://localhost:8080

# Remote wake up with a secret key
$ ewol remote http://localhost:8080 --wake.secret mysecret

Flags:
  -h, --help                 Help for remote
      --log.format string    Log format. Available options: text, json (default "text")
      --log.level string     Log level (default "info")
  -v, --version              Print version and exit
  -s, --wake.secret string   Secret key to wake up the device
```

The `ewol remote` command is a wrap-up command to call the `/wake` API, which means to use this command, the input `REMOTE_ADDRESS` must be an `ewol` server (running with the `--serve` flag).

You can also achieve this by other methods, such as `curl`:

```console
$ curl -sSL http://localhost:8080/wake -XPOST -H "Authorization: mysecret"
```
