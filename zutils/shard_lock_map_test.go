package zutils

import (
	"hash/fnv"
	"log/slog"
	"sort"
	"strconv"
	"testing"
)

type TestUser struct {
	name string
}

func TestCreatMap(t *testing.T) {
	slm := NewShardLockMaps()
	if slm.shards == nil {
		t.Error("shardLockMaps is null.")
	}

	if slm.Count() != 0 {
		t.Error("new shardLockMaps should be empty.")
	}
}

func TestSet(t *testing.T) {
	slm := NewShardLockMaps()

	slm.Set("user", "14March")
	mData := make(map[string]interface{})
	mData["aaa"] = 111
	mData["bbb"] = "222"
	slm.MSet(mData)
	slm.SetNX("user", "14March")
	bo := TestUser{"bo"}
	slm.SetNX("bo", bo)

	if slm.Count() != 4 {
		t.Error("shardLockMaps should contain exactly one elements.")
	}

}

func TestGet(t *testing.T) {
	slm := NewShardLockMaps()

	val, ok := slm.Get("user")

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	tony := TestUser{"tony"}
	slm.Set("tony", tony)

	tmp, ok := slm.Get("tony")
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	tony, ok = tmp.(TestUser)
	if !ok {
		t.Error("expecting an element, not null.")
	}

	if tony.name != "tony" {
		t.Error("item was modified.")
	}

}

func TestHas(t *testing.T) {
	slm := NewShardLockMaps()

	if slm.Has("Money") == true {
		t.Error("element shouldn't exists")
	}

	slm.Set("user", "14March")

	if slm.Has("user") == false {
		t.Error("element exists, user Has to return True.")
	}

}

func TestCount(t *testing.T) {
	slm := NewShardLockMaps()
	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	if slm.Count() != 100 {
		t.Error("Expecting 100 element within map.")
	}

}

func TestRemove(t *testing.T) {
	slm := NewShardLockMaps()

	david := TestUser{"David"}
	slm.Set("david", david)

	slm.Remove("david")

	if slm.Count() != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok := slm.Get("david")

	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}

	slm.Remove("user")

	isEmpty := slm.IsEmpty()
	if !isEmpty {
		t.Error("map should be empty.")
	}

}

func TestRemoveCb(t *testing.T) {
	slm := NewShardLockMaps()

	tony := TestUser{"tony"}
	slm.Set("tony", tony)
	david := TestUser{"david"}
	slm.Set("david", david)

	var (
		mapKey   string
		mapVal   interface{}
		wasFound bool
	)
	cb := func(key string, val interface{}, exists bool) bool {
		mapKey = key
		mapVal = val
		wasFound = exists

		if user, ok := val.(TestUser); ok {
			return user.name == "tony"
		}
		return false
	}

	result := slm.RemoveCb("tony", cb)
	if !result {
		t.Errorf("Result was not true")
	}

	if mapKey != "tony" {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != tony {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if slm.Has("tony") {
		t.Errorf("Key was not removed")
	}

	result = slm.RemoveCb("david", cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != "david" {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != david {
		t.Errorf("Wrong value was provided to the value")
	}

	if !wasFound {
		t.Errorf("Key was not found")
	}

	if !slm.Has("david") {
		t.Errorf("Key was removed")
	}

	result = slm.RemoveCb("danny", cb)
	if result {
		t.Errorf("Result was true")
	}

	if mapKey != "danny" {
		t.Error("Wrong key was provided to the callback")
	}

	if mapVal != nil {
		t.Errorf("Wrong value was provided to the value")
	}

	if wasFound {
		t.Errorf("Key was found")
	}

	if slm.Has("danny") {
		t.Errorf("Key was created")
	}
}

func TestPop(t *testing.T) {
	slm := NewShardLockMaps()

	_, exists := slm.Pop("user")
	if exists {
		t.Error("user should be not exists.")
	}

	slm.Set("user", "14March")

	val, exists := slm.Pop("user")
	if exists {
		t.Logf("user should be exists %v.", val)
	}

}

func TestBufferedIterator(t *testing.T) {
	slm := NewShardLockMaps()

	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	counter := 0
	for item := range slm.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}

}

func TestClear(t *testing.T) {
	slm := NewShardLockMaps()

	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	slm.Clear()

	if slm.Count() != 0 {
		t.Error("should have 0 elements.")
	}
}

func TestIterCb(t *testing.T) {
	slm := NewShardLockMaps()

	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	counter := 0
	slm.IterCb(func(key string, v interface{}) {
		_, ok := v.(TestUser)
		if !ok {
			t.Error("Expecting an user object")
		}

		counter++
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {
	slm := NewShardLockMaps()

	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	items := slm.Items()

	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}

}

func TestJsonMarshal(t *testing.T) {
	expected := "{\"a\":1,\"b\":2}"
	slm := NewShardLockMapsWithCount(2)
	slm.Set("a", 1)
	slm.Set("b", 2)
	j, err := slm.MarshalJSON()
	if err != nil {
		t.Error(err)
	}

	if string(j) != expected {
		t.Error("json", string(j), "differ from expected", expected)
		return
	}
}

func TestUnmarshalJSON(t *testing.T) {
	slm := NewShardLockMaps()
	bytes := []byte("{\"ccc\":1,\"ddd\":2}")

	err := slm.UnmarshalJSON(bytes)
	if err != nil {
		t.Error(err)
	}

	if slm.Count() != 2 {
		t.Error("Expecting count to be 2 once item was removed.")
	}

	for _, shard := range slm.shards {
		slog.Debug("shard.items", "items", shard.items)
	}

}

func TestKeys(t *testing.T) {
	slm := NewShardLockMaps()

	for i := 0; i < 100; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	keys := slm.Keys()
	if len(keys) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestFnv32(t *testing.T) {
	key := []byte("ABC")

	hasher := fnv.New32()
	_, err := hasher.Write(key)
	if err != nil {
		t.Errorf("%v", err)
	}
	if (&Fnv32Hash{}).Sum(string(key)) != hasher.Sum32() {
		t.Errorf("Bundled fnv32 produced %d, expected result from hash/fnv32 is %d", (&Fnv32Hash{}).Sum(string(key)), hasher.Sum32())
	}

}

func TestKeysWhenRemoving(t *testing.T) {
	slm := NewShardLockMaps()

	Total := 100
	for i := 0; i < Total; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}

	Num := 10
	for i := 0; i < Num; i++ {
		go func(c *ShardLockMaps, n int) {
			c.Remove(strconv.Itoa(n))
		}(&slm, i)
	}
	keys := slm.Keys()
	for _, k := range keys {
		if k == "" {
			t.Error("Empty keys returned")
		}
	}
}

func TestUnDrainedIterBuffered(t *testing.T) {
	slm := NewShardLockMaps()

	Total := 100
	for i := 0; i < Total; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}
	counter := 0

	// Iterate over elements.
	ch := slm.IterBuffered()
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
		if counter == 42 {
			break
		}
	}
	for i := Total; i < 2*Total; i++ {
		slm.Set(strconv.Itoa(i), TestUser{strconv.Itoa(i)})
	}
	for item := range ch {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have been right where we stopped")
	}

	counter = 0
	for item := range slm.IterBuffered() {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 200 {
		t.Error("We should have counted 200 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	slm := NewShardLockMaps()
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			slm.Set(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := slm.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			slm.Set(strconv.Itoa(i), i)

			// Retrieve item from map.
			val, _ := slm.Get(strconv.Itoa(i))

			// Write to channel inserted value.
			ch <- val.(int)
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])

	// Make sure map contains 1000 elements.
	if slm.Count() != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

// TestGetOrSet tests the GetOrSet method
func TestGetOrSet(t *testing.T) {
	slm := NewShardLockMaps()

	// Test setting a new value
	value, isSet := slm.GetOrSet("key1", "value1")
	if !isSet {
		t.Error("Expected isSet to be true for new key")
	}
	if value != "value1" {
		t.Error("Expected value to be 'value1'")
	}

	// Test getting existing value
	value, isSet = slm.GetOrSet("key1", "value2")
	if isSet {
		t.Error("Expected isSet to be false for existing key")
	}
	if value != "value1" {
		t.Error("Expected value to be 'value1'")
	}
}

// TestGetOrSetFunc tests the GetOrSetFunc method
func TestGetOrSetFunc(t *testing.T) {
	slm := NewShardLockMaps()

	// Test setting a new value with function
	value, isSet := slm.GetOrSetFunc("key1", func(key string) interface{} {
		return "generated_" + key
	})
	if !isSet {
		t.Error("Expected isSet to be true for new key")
	}
	if value != "generated_key1" {
		t.Error("Expected value to be 'generated_key1'")
	}

	// Test getting existing value
	value, isSet = slm.GetOrSetFunc("key1", func(key string) interface{} {
		return "should_not_be_called"
	})
	if isSet {
		t.Error("Expected isSet to be false for existing key")
	}
	if value != "generated_key1" {
		t.Error("Expected value to be 'generated_key1'")
	}
}

// TestGetOrSetFuncLock tests the GetOrSetFuncLock method
func TestGetOrSetFuncLock(t *testing.T) {
	slm := NewShardLockMaps()

	// Test setting a new value with function in lock
	value, isSet := slm.GetOrSetFuncLock("key1", func(key string) interface{} {
		return "locked_" + key
	})
	if !isSet {
		t.Error("Expected isSet to be true for new key")
	}
	if value != "locked_key1" {
		t.Error("Expected value to be 'locked_key1'")
	}

	// Test getting existing value
	value, isSet = slm.GetOrSetFuncLock("key1", func(key string) interface{} {
		return "should_not_be_called"
	})
	if isSet {
		t.Error("Expected isSet to be false for existing key")
	}
	if value != "locked_key1" {
		t.Error("Expected value to be 'locked_key1'")
	}
}

// TestMGet tests the MGet method
func TestMGet(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")
	slm.Set("key3", "value3")

	// Test getting multiple existing keys
	values := slm.MGet("key1", "key2", "key3")
	if len(values) != 3 {
		t.Error("Expected 3 values")
	}
	if values["key1"] != "value1" || values["key2"] != "value2" || values["key3"] != "value3" {
		t.Error("Expected correct values")
	}

	// Test getting mix of existing and non-existing keys
	values = slm.MGet("key1", "nonexistent", "key2")
	if len(values) != 2 {
		t.Error("Expected 2 values for mix of existing and non-existing keys")
	}
	if values["key1"] != "value1" || values["key2"] != "value2" {
		t.Error("Expected correct values for existing keys")
	}
}

// TestGetAll tests the GetAll method
func TestGetAll(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	all := slm.GetAll()
	if len(all) != 2 {
		t.Error("Expected 2 items in GetAll")
	}
	if all["key1"] != "value1" || all["key2"] != "value2" {
		t.Error("Expected correct values in GetAll")
	}
}

// TestLockFuncWithKey tests the LockFuncWithKey method
func TestLockFuncWithKey(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	// Test modifying a specific key's shard
	slm.LockFuncWithKey("key1", func(data map[string]interface{}) {
		if val, ok := data["key1"]; ok {
			data["key1"] = val.(string) + "_modified"
		}
	})

	value, ok := slm.Get("key1")
	if !ok || value != "value1_modified" {
		t.Error("Expected key1 to be modified")
	}
}

// TestRLockFuncWithKey tests the RLockFuncWithKey method
func TestRLockFuncWithKey(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	// Test reading from a specific key's shard
	var found bool
	slm.RLockFuncWithKey("key1", func(data map[string]interface{}) {
		if _, ok := data["key1"]; ok {
			found = true
		}
	})

	if !found {
		t.Error("Expected to find key1 in shard")
	}
}

// TestLockFunc tests the LockFunc method
func TestLockFunc(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	// Test modifying all shards
	slm.LockFunc(func(data map[string]interface{}) {
		for key, val := range data {
			if str, ok := val.(string); ok {
				data[key] = str + "_modified"
			}
		}
	})

	value1, _ := slm.Get("key1")
	value2, _ := slm.Get("key2")
	if value1 != "value1_modified" || value2 != "value2_modified" {
		t.Error("Expected all values to be modified")
	}
}

// TestRLockFunc tests the RLockFunc method
func TestRLockFunc(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	// Test reading from all shards
	var count int
	slm.RLockFunc(func(data map[string]interface{}) {
		count += len(data)
	})

	if count != 2 {
		t.Error("Expected to find 2 items across all shards")
	}
}

// TestClearWithFuncLock tests the ClearWithFuncLock method
func TestClearWithFuncLock(t *testing.T) {
	slm := NewShardLockMaps()
	slm.Set("key1", "value1")
	slm.Set("key2", "value2")

	var clearedKeys []string
	slm.ClearWithFuncLock(func(key string, val interface{}) {
		clearedKeys = append(clearedKeys, key)
	})

	if len(clearedKeys) != 2 {
		t.Error("Expected 2 keys to be cleared")
	}
	if slm.Count() != 0 {
		t.Error("Expected map to be empty after ClearWithFuncLock")
	}
}
