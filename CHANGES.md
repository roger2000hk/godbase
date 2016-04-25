# changes in reverse chronological order
### last updated 2016-04-24

### 2016-04-24 added iters & basic range ops
Added an Iter interface; and Any.Cut()/Find() ops. Added start/end params and adapted Any.Insert/Delete to use Iter.

### 2016-04-25 consolidated hash map implementations
Consolidated hash implementations by extracting slot logic into a Slots interface implemented as AnySlots, ESkipSlots, HashSlots, MapSlots & SkipSlots.

### 2016-04-26 added maps.Wrap & maps.Suffix
Added wrap struct for easy api extension and implemented a basic suffix map using the new functionality.