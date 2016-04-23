# beating The Map
### a glorious quest to write a faster, more flexible map

## why?
Good question. This rabbit hole was a bit deeper than anticipated, but I still feel it was worth the effort, given the results. To cut a long story short, I was itching for multi-capable sorted maps to implement in-memory indexing in godbase. I quickly found out that my lofty goal of adding multi-capability and sorting, while matching the performance of native maps; was far from a walk in the park. Which is probably part of the reason I couldn't find anyone else trying.

## how?
I ended up with two designs based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes. And another with embedded nodes, constant number of levels, using one node per value. These do pretty good by themselves, especially considering that they add sorting to the mix; with both hovering around 2-5 times slower than a native map. The embedded flavor usually wins the allocation race by a slim margin but pays the price of having a fixed number of levels for tiny / huge datasets. On top of these, I implemented two kinds of hash maps; using native maps for slots and one of the ordered maps for chains. These, when properly tuned, are consistently faster than native maps for millions of items.

## status
Only insert/delete functionality implemented so far, range searching and set operations are still in the oven.

## examples
test_map.go runs a basic test loop for each available map type, that should be enough to get started.