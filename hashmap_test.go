package hashmap_test

import (
	"testing"
	"github.com/DusanKasan/hashmap"
	"math/rand"
	"time"
)

func TestHashmap(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	//hash function that causes collisions
	hashFunc := func(i interface{}) int64 {
		v := i.(int64)
		if v != 0 && v % 5 == 0 {
			return v - 1
		}

		return v
	}

	//go multiple times over different sizes of input data
	for inputSize := 1; inputSize < 10; inputSize++ {
		for iteration := 1; iteration < 100; iteration++ {
			input := generateInputPool(inputSize)
			t.Logf("Running with input: %v", input)

			m := hashmap.New(hashFunc)

			for key, value := range (input) {
				t.Logf("Inserting key: %v", key)
				m.Insert(key, value)
			}

			for key, value := range (input) {
				v, found := m.Get(key)
				if !found {
					t.Errorf("Key not found: %v", key)
				} else {
					if v == nil {
						t.Errorf("Key %v has a nil value", key)
					} else if v.(int64) != value {
						t.Errorf("Key %v has wrong value. Expected %v, Got %v", key, value, v)
					}
				}
			}

			keys := getShuffledKeys(input)
			t.Logf("Shuffled keys: %v", keys)
			removedKeys := []int64{}
			for len(keys) > 0 {
				preservedKeys := []int64{}
				for i, k := range(keys) {

					if i == 0 {
						t.Logf("Removing key: %v", k)
						found := m.Remove(k)
						if !found {
							t.Errorf("Unable to find and remove key: %v", k)
						}
						t.Logf("Removed key: %v", k)

						removedKeys = append(removedKeys, k)
					} else {
						preservedKeys = append(preservedKeys, k)
					}
				}
				keys = preservedKeys

				for _, k := range(keys) {
					v, found := m.Get(k)
					if !found {
						t.Errorf("Key %v not found!", k)
					} else if v != input[k] {
						t.Errorf("Key %v has wrong value. Expected %v, Got %v", k, input[k], v)
					}
				}

				for _, k := range(removedKeys) {
					_, found := m.Get(k)
					if found {
						t.Errorf("Key %v found when it shouldn't have been!", k)
					}
				}
			}
		}
	}
}

//generate a map with randomized keys and values
func generateInputPool(size int) (map[int64]int64) {
	r := map[int64]int64{}
	values := rand.Perm(size * 4)

	for index, key := range (rand.Perm(size * 4)) {
		if index % 4 == 0 {
			r[int64(key)] = int64(values[index])
		}
	}

	return r
}

func getShuffledKeys(input map[int64]int64) []int64 {
	keys := []int64{}

	for key, _ := range(input) {
		keys = append(keys, key)
	}

	order := rand.Perm(len(keys))
	result := []int64{}
	for _, k := range(order) {
		result = append(result, keys[k])
	}

	return result
}
