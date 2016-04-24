# beating The Map
### a glorious quest to write a faster, more flexible map

## why?
Good question. This rabbit hole was a bit deeper than anticipated, but I still feel it was worth the effort given the results. To cut a long story short, I was itching for multi-capable sorted maps to implement in-memory indexing in godbase. I quickly found out that my lofty goal of adding multi-capability and sorting, while matching the performance of native maps; was far from a walk in the park. Which is probably part of the reason I couldn't find anyone else trying.

## how?
I ended up with two designs based on deterministic skip lists. One with a configurable number of levels, one node per level/value, and optionally slab-allocated nodes; the other with embedded nodes, constant number of levels, using one node per value. These do pretty good by themselves, especially considering that they add sorting to the mix; both currently hovering around 2-3 times slower than a native map. The embedded flavor usually wins the allocation race by a slim margin but pays the price of having a fixed number of levels for tiny / huge datasets. Still, somewhere along a line; not separately allocating nodes affects overall performance positively.

Once sorted maps were working properly, I had the crazy idea to put a hash on top just to see what happens. It turns out that dividing the dataset into a tuned number of ordered sets helps puts us consistently ahead of native maps in the synthetic performance game for millions of items. That's far from the end of the story though. I still haven't had enough time to ponder the consequences of having access to both hashed and ordered aspects of the data simultaneously, but I have a hunch it will bend the rules to my advantage in a number of tricky scenarios. Additionally; the hash adapter supports any kind of Map for slot chains, which opens the door for multi level hashing where each chain is another hash that further divides the dataset.

## status
Basic functionality and testing in place; bells, whistles & polish are still on the stove.

## api
More RISC/Lispy than your everyday set/map api. Turns out that providing an optimal api is half of implementing an optimal algorithm. And there's more low hanging fruit in the garden of set/map apis than most places. It's obvious to me that academic dogmatics and software (or life in general, for that matter) isn't really the match made in heaven it's being sold as.

### interfaces

### constructors