# Hashmap

A [Red-Black Tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree) [Hash Map](https://en.wikipedia.org/wiki/Hash_table) implementation in Golang that uses user-supplied hashing algorithms.

## Usage

The hashmap supports the classic Get, Insert, Remove operations you'd expect.

### Inserting
```go
//hash function must take interface{} and return int64
hashFunc := func(i interface{}) int64 {
    return int64(i.(int))
}

m := hashmap.New(hashFunc)
//insertion of key, value. keep in mind the key will be used as input to your hashFunc
m.Insert(4, 0)
//you can store different types
m.Insert(19, "hello")
//panics, because the hashFunc doesn't support string keys
m.Insert("fail", "oh no")
```

### Selecting
```go
hashFunc := func(i interface{}) int64 {
    return int64(i.(int))
}

m := hashmap.New(hashFunc)
m.Insert(4, 0)
m.Insert(19, 10)

//returns value as interface{} and found flag (true if the key was found)
value, found := m.Get(19)
// found will be false
value, found = m.Get(123)
```

### Removing
```go
hashFunc := func(i interface{}) int64 {
    return int64(i.(int))
}

m := hashmap.New(hashFunc)
m.Insert(4, 0)
m.Insert(19, 10)

//returns found flag (true if the key was found)
found := m.Remove(19)
// found will be false
found = m.Remove(123)
```

## Type safety concerns

As this hash map supports keys and values of any type (by type hinting interface{}), there could be concerns of type safety and runtime problems. The suggested way to work around this is to wrap the hash map into type-specific proxy with methods such as `Get(key KeyType) (value ValueType, found bool)` and do the type assertions there.

Direct support for code generation by this package is still considered but not yet implemented.

##TODO

- implement as thread safe
- threadsafety: don't lock the whole tree but separate nodes?
- CI
- Performance optimizations
- Performance tests and docs