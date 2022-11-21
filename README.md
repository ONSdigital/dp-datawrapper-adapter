
---

:warning: This repository will be archived in November 2022 as it is no longer in development. :warning:

---

# dp-datawrapper-adapter
Datawrapper adapter service

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default                    | Description
| ---------------------------- | -------------------------- | -----------
| BIND_ADDR                    | :28400                     | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                         | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                        | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                        | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| DATAWRAPPER_UI_URL           | https://app.datawrapper.de | Datawrapper UI URL
| DATAWRAPPER_API_URL          | https://api.datawrapper.de | Datawrapper API URL
| DATAWRAPPER_API_TOKEN        |                            | Datawrapper API Token of the admin user

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2022, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

