# beating The Map
### a glorious quest to write a faster, more flexible map

## why?
Good question. This rabbit hole was a bit deeper than anticipated, but I still feel it was worth the effort given the results. To cut a long story short, I was itching for multi-capable sorted maps to implement in-memory indexing in godbase. I quickly found out that my lofty goal of adding multi-capability and sorting, while matching the performance of native maps; was far from a walk in the park. Which is probably part of the reason I couldn't find anyone else trying.

## how?
I ended up with two designs based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes; the other with embedded nodes, constant number of levels, using one node per value. These do pretty good by themselves, especially considering that they add sorting to the mix; both currently hovering around 2-5 times slower than a native map. The embedded flavor usually wins the allocation race by a slim margin but pays the price of having a fixed number of levels for tiny / huge datasets. Still, somewhere along a line; not separately allocating nodes affects overall performance positively.

Once sorted maps were working properly, I had the crazy idea to put a hash on top just to see what happens. It turns out that when you don't need sorting, dividing the dataset into a tuned number of ordered sets helps puts you consistently ahead of native maps in the synthetic performance game for millions of items. That's far from the end of the story though. I still haven't had much time to ponder the consequences of having access to both aspects simultaneously, but I have a hunch it will bend the rules to my advantage in a number of tricky scenarios.

Any hashing strategy will work, dividing streams into a bucket per timeframe is an example that comes to mind. The Sorted interface is implemented per bucket for hash maps. Additionally; the hash adapter supports any kind of Map for slot chains, which opens the door for multi level hashing where each chain is another hash that further divides the dataset along an orthogonal axis.

## status
Only insert/delete/iterators implemented so far, set operations and polish are still in the oven.

## examples
test_map.go runs a basic test loop for each available map type, that should be enough to get started.