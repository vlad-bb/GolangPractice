package lru

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLruCache_PutAndGet(t *testing.T) {
	cache := NewLruCache(2)

	// Перевіряємо додавання і отримання значень
	cache.Put("key1", "value1")
	value, ok := cache.Get("key1")
	assert.True(t, ok, "expected key1 to be found")
	assert.Equal(t, "value1", value, "expected 'value1'")

	cache.Put("key2", "value2")
	value, ok = cache.Get("key2")
	assert.True(t, ok, "expected key2 to be found")
	assert.Equal(t, "value2", value, "expected 'value2'")

	// Перевіряємо, що значення зберігаються правильно
	value, ok = cache.Get("key1")
	assert.True(t, ok, "expected key1 to be found after adding 'key2'")
	assert.Equal(t, "value1", value, "expected 'value1'")
}

func TestLruCache_ReplaceLeastRecentlyUsed(t *testing.T) {
	cache := NewLruCache(2)

	// Додаємо два елементи
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")

	// Досягаємо ліміту кешу
	cache.Put("key3", "value3")

	// Перевіряємо, що перший доданий елемент було видалено
	_, ok := cache.Get("key1")
	assert.False(t, ok, "expected 'key1' to be evicted")

	// Перевіряємо, що інші елементи ще є в кеші
	value, ok := cache.Get("key2")
	assert.True(t, ok, "expected key2 to be found")
	assert.Equal(t, "value2", value, "expected 'value2'")

	value, ok = cache.Get("key3")
	assert.True(t, ok, "expected key3 to be found")
	assert.Equal(t, "value3", value, "expected 'value3'")
}

func TestLruCache_UpdateKey(t *testing.T) {
	cache := NewLruCache(2)

	// Додаємо елемент і оновлюємо його значення
	cache.Put("key1", "value1")
	cache.Put("key1", "updated_value1")

	// Перевіряємо, що значення оновлено
	value, ok := cache.Get("key1")
	assert.True(t, ok, "expected key1 to be found")
	assert.Equal(t, "updated_value1", value, "expected 'updated_value1'")
}

func TestLruCache_LRUBehavior(t *testing.T) {
	cache := NewLruCache(2)

	// Додаємо три елементи, перевіряємо видалення найстарішого
	cache.Put("key1", "value1")
	cache.Put("key2", "value2")
	cache.Get("key1")           // key1 тепер найновіший
	cache.Put("key3", "value3") // key2 має бути видалений

	// Перевіряємо, що key2 було видалено
	_, ok := cache.Get("key2")
	assert.False(t, ok, "expected 'key2' to be evicted")

	// Перевіряємо, що key1 та key3 залишились
	value, ok := cache.Get("key1")
	assert.True(t, ok, "expected key1 to be found")
	assert.Equal(t, "value1", value, "expected 'value1'")

	value, ok = cache.Get("key3")
	assert.True(t, ok, "expected key3 to be found")
	assert.Equal(t, "value3", value, "expected 'value3'")
}
