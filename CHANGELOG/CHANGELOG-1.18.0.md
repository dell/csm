<!--toc-->

# v1.18.0

## Changelog since v1.17.2

## Known Issues

## Changes by Kind

### Deprecation

### Features

- Reduced PowerStore REST API load by deriving `NodeGetVolumeStats` host details from the `host_volume_mapping` response, eliminating per-mapping `GetHost` calls in the CSI PowerStore driver. (ECSDF-50)
- Changed `csm-metrics-powerstore` default metrics polling interval from `20s` to `300s`. (ECSDF-50)
- Changed `csm-metrics-powerstore` default `POWERSTORE_MAX_CONCURRENT_QUERIES` value from `10` to `5`. (ECSDF-50)

### Bugs
