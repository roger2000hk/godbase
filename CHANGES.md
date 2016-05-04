# changes in reverse chronological order
#### last updated 2016-05-04

### 2016-05-04 Changed fix to use int64 instead of big.Int

### 2016-05-04 Simplified hash map implementations
Cut versions using interface references for slots.

### 2016-05-04 Added ok result to Set.First/Last()
Added result to signal if key was found in Set.First/Last().

### 2016-05-03 Switched recs.Basic to custom type with sets.Sort for keys

### 2016-05-02 Added sets package

### 2016-05-01 Converted Suffix to proper implementation and added Trace wrap

### 2016-04-30 Moved Any-interfaces to godbase package
This means that maps.Any is named godbase.Map, same goes for cols.Any, recs.Any, idxs.Any & tbls.Any.
Also moved Key & Iter to godbase.

### 2016-04-28 added maps.Any.Clear/While
Added methods for clearing map & for looping with callback.

### 2016-04-27 added maps.Any.Get/Set
Added simplified methods to get / set value for key.

### 2016-04-27 changed maps.Any.Delete/Find/Insert iter semantics
All methods now return iter to current if found, not prev.

### 2016-04-27 added maps.Any.First
Added methods to get iter to first elem and implemented for sorted maps.

### 2016-04-27 added maps.Any.New
Added method to return a new map of the same type.

### 2016-04-26 added maps.Wrap & maps.Suffix
Added wrap struct for easy api extension and implemented a basic suffix map using the new functionality.

### 2016-04-25 consolidated hash map implementations
Consolidated hash implementations by extracting slot logic into a Slots interface implemented as AnySlots, ESkipSlots, HashSlots, MapSlots & SkipSlots.

### 2016-04-24 added iters & basic range ops
Added an Iter interface; and Any.Cut()/Find() ops. Added start/end params and adapted Any.Insert/Delete to use Iter.

