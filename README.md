# kv
A collection of a few Go packages for Key Value storage like caches.

## Packages
* `kv1` - A Key Value sharded cache with time expiration using the swiss map.
* `kv1s` - A Key Value sharded cache using swiss map without any auto eviction.
* `kv2` - A Key Value sharded cache with a max size and automated eviction using the swiss map.
* `stack` - A last in, first out stack implementation without concurrency support. Used in the ``kv2`` package.
