# changes in reverse chronological order
### last updated 2016-04-27

### 2016-04-24 added iters & basic range ops
Added an Iter interface; and Any.Cut()/Find() ops. Added start/end params and adapted Any.Insert/Delete to use Iter.

### 2016-04-25 consolidated hash map implementations
Consolidated hash implementations by extracting slot logic into a Slots interface implemented as AnySlots, ESkipSlots, HashSlots, MapSlots & SkipSlots.

### 2016-04-26 added maps.Wrap & maps.Suffix
Added wrap struct for easy api extension and implemented a basic suffix map using the new functionality.

### 2016-04-27 added maps.Any.Get/Set
Added simplified methods to get / set value for key.

### 2016-04-27 changed maps.Any.Delete/Find/Insert iter semantics
All methods now return iter to current if found, not prev.

### 2016-04-27 added maps.Any.First
Added methods to get iter to first elem and implemented for sorted maps.

### 2016-04-27 added maps.Any.New
Added method to return a new map of the same type.