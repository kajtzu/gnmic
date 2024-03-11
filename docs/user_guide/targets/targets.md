# Targets

Sometimes it is needed to perform an operation on multiple devices; be it getting the same leaf value from a given set of the network elements or setting a certain configuration element to some value.

For cases like that `gnmic` offers support for multiple targets operations which a user can configure both via CLI flags as well as with the [file-based configuration](../configuration_file.md).

### CLI configuration

Specifying multiple targets in the CLI is as easy as repeating the [`--address`](../../global_flags.md#address) flag.

```shell
❯ gnmic -a router1.lab.net:57400 \
        -a router2.lab.net:57400 \
        get --path /configure/system/name
```

### File-based configuration

With the file-based configuration a user has two options to specify multiple targets:

* using `address` option
* using `targets` option

#### address option

With `address` option the user must provide a list of addresses. In the YAML format that would look like that:

```yaml
address:
  - "router1.lab.net:57400"
  - "router2.lab.net:57400"
```

The limitation this approach has is that it is impossible to set different credentials for the targets, they will essentially share the credentials specified in a file or via flags.

#### target option

With the `targets` option it is possible to set target specific options (such as credentials, subscriptions, TLS config, outputs), and thus this option is recommended to use:

```yaml
targets:
  router1.lab.net:
    timeout: 2s
    username: r1
    password: gnmi_pass
  router2.lab.net:57000:
    username: r2
    password: gnmi_pass
    tls-key: /path/file1
    tls-cert: /path/file2
```

The target address is defined as the key under the `targets` section of the configuration file. The default port (57400) can be omitted as demonstrated with `router1.lab.net` target address. Have a look at the [file-based targets configuration](https://github.com/openconfig/gnmic/blob/main/config.yaml) example to get a glimpse of what it is capable of.

The target inherits the globally defined options if the matching options are not set on a target level. For example, if a target doesn't have a username defined, it will use the username value set on a global level.

#### secure/insecure connections

`gnmic` supports both secure and insecure gRPC connections to the target.

##### insecure connection

Using the `--insecure` flag it is possible to establish an insecure gRPC connection to the target.

```bash
gnmic -a router1:57400 \
      --insecure \
      get --path /configure/system/name
```

##### secure connection

- A one way secure connection without target certificate verification can be established using the `--skip-verify` flag.

```bash
gnmic -a router1:57400 \
      --skip-verify \
      get --path /configure/system/name
```

- Adding target certificate verification can be done using the `--tls-ca` flag.

```bash
gnmic -a router1:57400 \
      --tls-ca /path/to/ca/file \
      get --path /configure/system/name
```

- A two way secure connection can be established using the `--tls-cert` `--tls-key` flags.

```bash
gnmic -a router1:57400 \
      --tls-cert /path/to/certificate/file \
      --tls-key /path/to/certificate/file \
      get --path /configure/system/name
```

- It is also possible to control the negotiated TLS version using the `--tls-min-version`, `--tls-max-version` and `--tls-version` (preferred TLS version) flags.

##### Controlling the advertised cipher suites

It's possible to configure the advertised cipher suites gNMIc's gNMI client advertises to the target.
This can be done by setting the `tls-min-version` and `tls-max-version` or by explicitly listing cipher suites to be advertised.

By default the below list is advertised:

| Name                                           | Key Exchange | Auth      | Enc                  | MAC       |
|------------------------------------------------|--------------|-----------|----------------------|-----------|
| TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256        | ECDHE        | ECDSA     | AES_128_GCM          | SHA256    |
| TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256          | ECDHE        | RSA       | AES_128_GCM          | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384        | ECDHE        | ECDSA     | AES_256_GCM          | SHA384    |
| TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384          | ECDHE        | RSA       | AES_256_GCM          | SHA384    |
| TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256  | ECDHE        | ECDSA     | CHACHA20_POLY1305    | SHA256    |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256    | ECDHE        | RSA       | CHACHA20_POLY1305    | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA           | ECDHE        | ECDSA     | AES_128_CBC          | SHA       |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA             | ECDHE        | RSA       | AES_128_CBC          | SHA       |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA           | ECDHE        | ECDSA     | AES_256_CBC          | SHA       |
| TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA             | ECDHE        | RSA       | AES_256_CBC          | SHA       |
| TLS_RSA_WITH_AES_128_GCM_SHA256                | RSA          | RSA       | AES_128_GCM          | SHA256    |
| TLS_RSA_WITH_AES_256_GCM_SHA384                | RSA          | RSA       | AES_256_GCM          | SHA384    |
| TLS_RSA_WITH_AES_128_CBC_SHA                   | RSA          | RSA       | AES_128_CBC          | SHA       |
| TLS_RSA_WITH_AES_256_CBC_SHA                   | RSA          | RSA       | AES_256_CBC          | SHA       |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA            | ECDHE        | RSA       | 3DES_EDE_CBC         | SHA       |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA                  | RSA          | RSA       | 3DES_EDE_CBC         | SHA       |
| TLS_AES_128_GCM_SHA256                         | (TLS 1.3)    | (TLS 1.3) | AES_128_GCM          | SHA256    |
| TLS_AES_256_GCM_SHA384                         | (TLS 1.3)    | (TLS 1.3) | AES_256_GCM          | SHA384    |
| TLS_CHACHA20_POLY1305_SHA256                   | (TLS 1.3)    | (TLS 1.3) | CHACHA20_POLY1305    | SHA256    |

If the `tls-max-version` is set to "1.2", the TLS1.3 cipher suites will not be included:

| Name                                           | Key Exchange | Auth  | Enc               | MAC       |
|------------------------------------------------|--------------|-------|-------------------|-----------|
| TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256        | ECDHE        | ECDSA | AES_128_GCM       | SHA256    |
| TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256          | ECDHE        | RSA   | AES_128_GCM       | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384        | ECDHE        | ECDSA | AES_256_GCM       | SHA384    |
| TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384          | ECDHE        | RSA   | AES_256_GCM       | SHA384    |
| TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256  | ECDHE        | ECDSA | CHACHA20_POLY1305 | SHA256    |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256    | ECDHE        | RSA   | CHACHA20_POLY1305 | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA           | ECDHE        | ECDSA | AES_128_CBC       | SHA       |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA             | ECDHE        | RSA   | AES_128_CBC       | SHA       |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA           | ECDHE        | ECDSA | AES_256_CBC       | SHA       |
| TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA             | ECDHE        | RSA   | AES_256_CBC       | SHA       |
| TLS_RSA_WITH_AES_128_GCM_SHA256                | RSA          | RSA   | AES_128_GCM       | SHA256    |
| TLS_RSA_WITH_AES_256_GCM_SHA384                | RSA          | RSA   | AES_256_GCM       | SHA384    |
| TLS_RSA_WITH_AES_128_CBC_SHA                   | RSA          | RSA   | AES_128_CBC       | SHA       |
| TLS_RSA_WITH_AES_256_CBC_SHA                   | RSA          | RSA   | AES_256_CBC       | SHA       |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA            | ECDHE        | RSA   | 3DES_EDE_CBC      | SHA       |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA                  | RSA          | RSA   | 3DES_EDE_CBC      | SHA       |

If the `tls-max-version` and `tls-min-version` are set to "1.1", the below list of cipher suites is advertised:

| Name                                   | Key Exchange | Auth  | Enc          | MAC |
|----------------------------------------|--------------|-------|--------------|-----|
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA   | ECDHE        | ECDSA | AES_128_CBC  | SHA |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA     | ECDHE        | RSA   | AES_128_CBC  | SHA |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA   | ECDHE        | ECDSA | AES_256_CBC  | SHA |
| TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA     | ECDHE        | RSA   | AES_256_CBC  | SHA |
| TLS_RSA_WITH_AES_128_CBC_SHA           | RSA          | RSA   | AES_128_CBC  | SHA |
| TLS_RSA_WITH_AES_256_CBC_SHA           | RSA          | RSA   | AES_256_CBC  | SHA |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA    | ECDHE        | RSA   | 3DES_EDE_CBC | SHA |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA          | RSA          | RSA   | 3DES_EDE_CBC | SHA |

If you want to control which cipher suites are sent and in what order of preference, you can set the `cipher-suites` field under the target:

```yaml
targets:
  target1:
    # other fields
    cipher-suites:
      - TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      - TLS_AES_128_GCM_SHA256
```

The full list of supported cipher suites is:

| Name                                           | Key Exchange | Auth      | Enc                | MAC       |
|------------------------------------------------|--------------|-----------|--------------------|-----------|
| TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256        | ECDHE        | ECDSA     | AES_128_GCM        | SHA256    |
| TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256          | ECDHE        | RSA       | AES_128_GCM        | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384        | ECDHE        | ECDSA     | AES_256_GCM        | SHA384    |
| TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384          | ECDHE        | RSA       | AES_256_GCM        | SHA384    |
| TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256  | ECDHE        | ECDSA     | CHACHA20_POLY1305  | SHA256    |
| TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256    | ECDHE        | RSA       | CHACHA20_POLY1305  | SHA256    |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA           | ECDHE        | ECDSA     | AES_128_CBC        | SHA       |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA             | ECDHE        | RSA       | AES_128_CBC        | SHA       |
| TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA           | ECDHE        | ECDSA     | AES_256_CBC        | SHA       |
| TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA             | ECDHE        | RSA       | AES_256_CBC        | SHA       |
| TLS_RSA_WITH_AES_128_GCM_SHA256                | RSA          | RSA       | AES_128_GCM        | SHA256    |
| TLS_RSA_WITH_AES_256_GCM_SHA384                | RSA          | RSA       | AES_256_GCM        | SHA384    |
| TLS_RSA_WITH_AES_128_CBC_SHA                   | RSA          | RSA       | AES_128_CBC        | SHA       |
| TLS_RSA_WITH_AES_256_CBC_SHA                   | RSA          | RSA       | AES_256_CBC        | SHA       |
| TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA            | ECDHE        | RSA       | 3DES_EDE_CBC       | SHA       |
| TLS_RSA_WITH_3DES_EDE_CBC_SHA                  | RSA          | RSA       | 3DES_EDE_CBC       | SHA       |
| TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256        | ECDHE        | ECDSA     | AES_128_CBC        | SHA256    |
| TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256          | ECDHE        | RSA       | AES_128_CBC        | SHA256    |
| TLS_RSA_WITH_AES_128_CBC_SHA256                | RSA          | RSA       | AES_128_CBC        | SHA256    |
| TLS_ECDHE_ECDSA_WITH_RC4_128_SHA               | ECDHE        | ECDSA     | RC4_128            | SHA       |
| TLS_ECDHE_RSA_WITH_RC4_128_SHA                 | ECDHE        | RSA       | RC4_128            | SHA       |
| TLS_RSA_WITH_RC4_128_SHA                       | RSA          | RSA       | RC4_128            | SHA       |
| TLS_AES_128_GCM_SHA256                         | (TLS 1.3)    | (TLS 1.3) | AES_128_GCM        | SHA256    |
| TLS_AES_256_GCM_SHA384                         | (TLS 1.3)    | (TLS 1.3) | AES_256_GCM        | SHA384    |
| TLS_CHACHA20_POLY1305_SHA256                   | (TLS 1.3)    | (TLS 1.3) | CHACHA20_POLY1305  | SHA256    |

#### target configuration options

Target supported options:

```yaml
targets:
  # target name or an address (IP or DNS name).
  # if an address is set it can include a port number or not,
  # if a port is not included, the default gRPC port will be added.
  target_key:
    # target name, will default to the target_key if not specified
    name: target_key
    # target address, if missing the target_key is used as an address.
    # supports comma separated addresses.
    # if any of the addresses is missing a port, the default gRPC port will be added.
    # if multiple addresses are set, all of them will be tried simultaneously,
    # the first established gRPC connection will be used, the other attempts will be canceled.
    address:
    # target username
    username:
    # target password
    password:
    # authentication token, 
    # applied only in the case of a secure gRPC connection.
    token: 
    # target RPC timeout
    timeout:
    # establish an insecure connection
    insecure:
    # path to tls ca file
    tls-ca:
    # path to tls certificate
    tls-cert:
    # path to tls key
    tls-key:
    # max tls version to use during negotiation
    tls-max-version:
    # min tls version to use during negotiation
    tls-min-version:
    # preferred tls version to use during negotiation
    # this value overwrites both tls-min-version and 
    # tls-max-version
    tls-version:
    # enable logging of a pre-master TLS secret
    log-tls-secret:
    # do not verify the target certificate when using tls
    skip-verify:
    # server name used to verify the hostname on the returned 
    # certificates unless skip-verify is true.    
    tls-server-name:
    # list of subscription names to establish for this target.
    # if empty it defaults to all subscriptions defined under
    # the main level `subscriptions` field
    subscriptions:
    # string, case insensitive, defines the gNMI encoding to be used for 
    # the subscriptions to be established for this target.
    # This encoding value applies only if the subscription configuration does
    # NOT explicitly define an encoding.
    encoding:
    # list of output names to which the gnmi data will be written.
    # if empty if defaults to all outputs defined under
    # the main level `outputs` field
    outputs:
    # number of subscribe responses to keep in buffer before writing
    # the target outputs
    buffer-size:
    # target retry period
    retry:
    # list of tags, relevant when clustering is enabled.
    tags:
    # a mapping of static tags to add to all events from this target.
    # each key/value pair in this mapping will be added to metadata
    # on all events
    event-tags:
    # list of proto file names to decode protoBytes values
    proto-files:
    # list of directories to look for the proto files
    proto-dirs:
    # enable grpc gzip compression
    gzip: 
    # proxy type and address, only SOCKS5 is supported currently
    # example: socks5://<address>:<port>
    proxy:
    # list of custom TLS cipher suites to advertise to the target 
    # during the TLS handshake.
    cipher-suites:
    # a duration, sets the TCP keepalive time and keepalive interval.
    # The number of keepalive probes to send before sending a TCP RST
    # is not configurable, it inherits its value from the linux kernel
    # net.ipv4.tcp_keepalive_probes which usually has a default value of 9.
    # When empty or set to 0s, the Golang default (15s) applies.
    # Disabled if set to a negative value.
    tcp-keepalive: 0s
```

### Example

Whatever configuration option you choose, the multi-targeted operations will uniformly work across the commands that support them.

Consider the `get` command acting on two routers getting their names:

```shell
❯ gnmic -a router1.lab.net:57400 \
        -a router2.lab.net:57400 \
        get --path /configure/system/name

[router1.lab.net:57400] {
[router1.lab.net:57400]   "source": "router1.lab.net:57400",
[router1.lab.net:57400]   "timestamp": 1593009759618786781,
[router1.lab.net:57400]   "time": "2020-06-24T16:42:39.618786781+02:00",
[router1.lab.net:57400]   "updates": [
[router1.lab.net:57400]     {
[router1.lab.net:57400]       "Path": "configure/system/name",
[router1.lab.net:57400]       "values": {
[router1.lab.net:57400]         "configure/system/name": "gnmic_r1"
[router1.lab.net:57400]       }
[router1.lab.net:57400]     }
[router1.lab.net:57400]   ]
[router1.lab.net:57400] }

[router2.lab.net:57400] {
[router2.lab.net:57400]   "source": "router2.lab.net:57400",
[router2.lab.net:57400]   "timestamp": 1593009759748265232,
[router2.lab.net:57400]   "time": "2020-06-24T16:42:39.748265232+02:00",
[router2.lab.net:57400]   "updates": [
[router2.lab.net:57400]     {
[router2.lab.net:57400]       "Path": "configure/system/name",
[router2.lab.net:57400]       "values": {
[router2.lab.net:57400]         "configure/system/name": "gnmic_r2"
[router2.lab.net:57400]       }
[router2.lab.net:57400]     }
[router2.lab.net:57400]   ]
[router2.lab.net:57400] }
```

Notice how in the output the different gNMI targets are prefixed with the target address to make the output easy to read. If those prefixes are not needed, you can make them disappear with [`--no-prefix`](../../global_flags.md#no-prefix) global flag.
