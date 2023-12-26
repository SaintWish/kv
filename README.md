# kv
A collection of a few Go packages for Key Value storage like caches.

## Packages
All of the packages of generic support for safety. Some of the packages use the [swiss map](https://github.com/dolthub/swiss) instead of the default Go map.

* `kv1` - A Key Value sharded cache with time expiration. Uses ``swiss`` map.
* `kv1s` - A Key Value sharded cache without any auto eviction. Uses ``swiss`` map.
* `kv2` - A Key Value sharded cache with a max size and automated eviction. Uses ``swiss`` map.
* `kvmap` - A Key Value sharded cache using vanilla Go map with no auto eviction.
* `ccmap` - A concurrent safe default Go map without sharding.
* `stack` - A last in, first out stack implementation without concurrency support. Used in the ``kv2`` package.

## Licensing
The [swiss map](https://github.com/dolthub/swiss) and this package are licensed with Apache-2.0