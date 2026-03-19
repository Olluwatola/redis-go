package main

import (
	"bufio"
	"strings"
	"sync"
	"testing"
	"time"
)

// =============================================================================
// Store Tests
// =============================================================================

func TestStore_SetAndGet(t *testing.T) {
	store := NewStore()

	// Test basic set and get
	store.SetString("key1", "value1")
	value, exists, err := store.GetString("key1")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected key1 to exist")
	}
	if value != "value1" {
		t.Errorf("expected 'value1', got '%s'", value)
	}
}

func TestStore_GetNonExistent(t *testing.T) {
	store := NewStore()

	value, exists, err := store.GetString("nonexistent")

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected key to not exist")
	}
	if value != "" {
		t.Errorf("expected empty string, got '%s'", value)
	}
}

func TestStore_Overwrite(t *testing.T) {
	store := NewStore()

	store.SetString("key", "original")
	store.SetString("key", "updated")

	value, _, _ := store.GetString("key")
	if value != "updated" {
		t.Errorf("expected 'updated', got '%s'", value)
	}
}

func TestStore_Keys(t *testing.T) {
	store := NewStore()

	store.SetString("a", "1")
	store.SetString("b", "2")
	store.SetString("c", "3")

	keys := store.Keys()

	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	// Check all keys are present (order not guaranteed)
	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	for _, expected := range []string{"a", "b", "c"} {
		if !keyMap[expected] {
			t.Errorf("expected key '%s' to be present", expected)
		}
	}
}

func TestStore_Size(t *testing.T) {
	store := NewStore()

	if store.Size() != 0 {
		t.Errorf("expected size 0, got %d", store.Size())
	}

	store.SetString("key1", "value1")
	store.SetString("key2", "value2")

	if store.Size() != 2 {
		t.Errorf("expected size 2, got %d", store.Size())
	}
}

func TestStore_Concurrent(t *testing.T) {
	store := NewStore()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			store.SetString(key, "value")
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('a' + (i % 26)))
			store.GetString(key)
		}(i)
	}

	wg.Wait()
	// Test passes if no race conditions or deadlocks occur
}

// =============================================================================
// RESP Parsing Tests
// =============================================================================

func TestParseRESP_SimpleArray(t *testing.T) {
	// *2\r\n$4\r\nPING\r\n$4\r\ntest\r\n
	input := "*2\r\n$4\r\nPING\r\n$4\r\ntest\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	result, err := parseRESP(reader)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(result))
	}
	if result[0] != "PING" {
		t.Errorf("expected 'PING', got '%s'", result[0])
	}
	if result[1] != "test" {
		t.Errorf("expected 'test', got '%s'", result[1])
	}
}

func TestParseRESP_SingleElement(t *testing.T) {
	input := "*1\r\n$4\r\nPING\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	result, err := parseRESP(reader)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
	if result[0] != "PING" {
		t.Errorf("expected 'PING', got '%s'", result[0])
	}
}

func TestParseRESP_SetCommand(t *testing.T) {
	// SET key value
	input := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	result, err := parseRESP(reader)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}
	if result[0] != "SET" {
		t.Errorf("expected 'SET', got '%s'", result[0])
	}
	if result[1] != "key" {
		t.Errorf("expected 'key', got '%s'", result[1])
	}
	if result[2] != "value" {
		t.Errorf("expected 'value', got '%s'", result[2])
	}
}

func TestParseRESP_InvalidType(t *testing.T) {
	input := "+PING\r\n" // Simple string, not array
	reader := bufio.NewReader(strings.NewReader(input))

	_, err := parseRESP(reader)

	if err == nil {
		t.Error("expected error for invalid type")
	}
}

// =============================================================================
// Command Handler Tests
// =============================================================================

func TestHandleCommand_Ping(t *testing.T) {
	result := handleCommand([]string{"PING"})

	if result != "+PONG\r\n" {
		t.Errorf("expected '+PONG\\r\\n', got %q", result)
	}
}

func TestHandleCommand_PingLowercase(t *testing.T) {
	result := handleCommand([]string{"ping"})

	if result != "+PONG\r\n" {
		t.Errorf("expected '+PONG\\r\\n', got %q", result)
	}
}

func TestHandleCommand_Echo(t *testing.T) {
	result := handleCommand([]string{"ECHO", "hello"})

	expected := "$5\r\nhello\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestHandleCommand_EchoEmpty(t *testing.T) {
	result := handleCommand([]string{"ECHO"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_Set(t *testing.T) {
	// Reset store for test isolation
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "mykey", "myvalue"})

	if result != "+OK\r\n" {
		t.Errorf("expected '+OK\\r\\n', got %q", result)
	}

	// Verify value was stored
	value, exists, _ := globalStore.GetString("mykey")
	if !exists || value != "myvalue" {
		t.Errorf("expected 'myvalue', got '%s' (exists: %v)", value, exists)
	}
}

func TestHandleCommand_SetMissingArgs(t *testing.T) {
	result := handleCommand([]string{"SET", "key"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_SetWithEX(t *testing.T) {
	globalStore = NewStore()

	// SET key value EX 10 (expires in 10 seconds)
	result := handleCommand([]string{"SET", "exkey", "exvalue", "EX", "10"})

	if result != "+OK\r\n" {
		t.Errorf("expected '+OK\\r\\n', got %q", result)
	}

	// Key should exist immediately
	value, exists, _ := globalStore.GetString("exkey")
	if !exists || value != "exvalue" {
		t.Errorf("expected 'exvalue', got '%s' (exists: %v)", value, exists)
	}
}

func TestHandleCommand_SetWithPX(t *testing.T) {
	globalStore = NewStore()

	// SET key value PX 5000 (expires in 5000 milliseconds)
	result := handleCommand([]string{"SET", "pxkey", "pxvalue", "PX", "5000"})

	if result != "+OK\r\n" {
		t.Errorf("expected '+OK\\r\\n', got %q", result)
	}

	// Key should exist immediately
	value, exists, _ := globalStore.GetString("pxkey")
	if !exists || value != "pxvalue" {
		t.Errorf("expected 'pxvalue', got '%s' (exists: %v)", value, exists)
	}
}

func TestHandleCommand_SetWithPX_Expiration(t *testing.T) {
	globalStore = NewStore()

	// SET key value PX 50 (expires in 50 milliseconds)
	result := handleCommand([]string{"SET", "shortkey", "shortvalue", "PX", "50"})

	if result != "+OK\r\n" {
		t.Errorf("expected '+OK\\r\\n', got %q", result)
	}

	// Key should exist immediately
	value, exists, _ := globalStore.GetString("shortkey")
	if !exists || value != "shortvalue" {
		t.Errorf("expected key to exist immediately")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Key should now be expired
	_, exists, _ = globalStore.GetString("shortkey")
	if exists {
		t.Error("expected key to be expired")
	}
}

func TestHandleCommand_SetEX_MissingValue(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "EX"})

	if !strings.HasPrefix(result, "-ERR syntax error") {
		t.Errorf("expected syntax error, got %q", result)
	}
}

func TestHandleCommand_SetPX_MissingValue(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "PX"})

	if !strings.HasPrefix(result, "-ERR syntax error") {
		t.Errorf("expected syntax error, got %q", result)
	}
}

func TestHandleCommand_SetEX_InvalidInteger(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "EX", "notanumber"})

	if !strings.HasPrefix(result, "-ERR value is not an integer") {
		t.Errorf("expected integer error, got %q", result)
	}
}

func TestHandleCommand_SetEX_NegativeValue(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "EX", "-5"})

	if !strings.HasPrefix(result, "-ERR value is not an integer") {
		t.Errorf("expected integer error, got %q", result)
	}
}

func TestHandleCommand_SetEX_ZeroValue(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "EX", "0"})

	if !strings.HasPrefix(result, "-ERR value is not an integer") {
		t.Errorf("expected integer error, got %q", result)
	}
}

func TestHandleCommand_Set_UnsupportedOption(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"SET", "key", "value", "INVALID"})

	if !strings.HasPrefix(result, "-ERR unsupported option") {
		t.Errorf("expected unsupported option error, got %q", result)
	}
}

func TestHandleCommand_SetEX_LowercaseFlag(t *testing.T) {
	globalStore = NewStore()

	// Lowercase 'ex' should work
	result := handleCommand([]string{"SET", "lckey", "lcvalue", "ex", "10"})

	if result != "+OK\r\n" {
		t.Errorf("expected '+OK\\r\\n', got %q", result)
	}

	value, exists, _ := globalStore.GetString("lckey")
	if !exists || value != "lcvalue" {
		t.Errorf("expected 'lcvalue', got '%s' (exists: %v)", value, exists)
	}
}

// =============================================================================
// Store Expiration Tests
// =============================================================================

func TestStore_SetWithExpire(t *testing.T) {
	store := NewStore()

	expireAt := time.Now().Add(1 * time.Hour)
	store.SetStringWithExpire("key", "value", expireAt)

	value, exists, _ := store.GetString("key")
	if !exists || value != "value" {
		t.Errorf("expected 'value', got '%s' (exists: %v)", value, exists)
	}
}

func TestStore_SetWithExpire_Expired(t *testing.T) {
	store := NewStore()

	// Set with expiration in the past
	expireAt := time.Now().Add(-1 * time.Second)
	store.SetStringWithExpire("key", "value", expireAt)

	// Should not exist (lazy expiration)
	_, exists, _ := store.GetString("key")
	if exists {
		t.Error("expected key to be expired")
	}
}

func TestStore_SetClearsExpiration(t *testing.T) {
	store := NewStore()

	// Set with short expiration
	expireAt := time.Now().Add(50 * time.Millisecond)
	store.SetStringWithExpire("key", "value1", expireAt)

	// Overwrite without expiration
	store.SetString("key", "value2")

	// Wait past original expiration
	time.Sleep(100 * time.Millisecond)

	// Key should still exist (expiration was cleared)
	value, exists, _ := store.GetString("key")
	if !exists || value != "value2" {
		t.Errorf("expected 'value2', got '%s' (exists: %v)", value, exists)
	}
}

func TestHandleCommand_Get(t *testing.T) {
	// Reset store and set a value
	globalStore = NewStore()
	globalStore.SetString("testkey", "testvalue")

	result := handleCommand([]string{"GET", "testkey"})

	expected := "$9\r\ntestvalue\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestHandleCommand_GetNonExistent(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"GET", "nonexistent"})

	expected := "$-1\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestHandleCommand_GetMissingArgs(t *testing.T) {
	result := handleCommand([]string{"GET"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_UnknownCommand(t *testing.T) {
	result := handleCommand([]string{"UNKNOWN"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_EmptyCommand(t *testing.T) {
	result := handleCommand([]string{})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

// =============================================================================
// Encoding Tests
// =============================================================================

func TestEncodeBulkString(t *testing.T) {
	result := encodeBulkString("hello")

	expected := "$5\r\nhello\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeBulkString_Empty(t *testing.T) {
	result := encodeBulkString("")

	expected := "$0\r\n\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeSimpleString(t *testing.T) {
	result := encodeSimpleString("OK")

	expected := "+OK\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestEncodeError(t *testing.T) {
	result := encodeError("ERR unknown command")

	expected := "-ERR unknown command\r\n"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

// =============================================================================
// List Tests
// =============================================================================

func TestLPush(t *testing.T) {
	store := NewStore()

	// Test creating new list
	list, err := store.GetOrCreateList("mylist")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Push one value
	length := list.LPush("hello")
	if length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}

	if list.Head.Value != "hello" {
		t.Errorf("Expected 'hello', got %s", list.Head.Value)
	}

	// Push another value
	length = list.LPush("world")
	if length != 2 {
		t.Errorf("Expected length 2, got %d", length)
	}

	// Check order (LPUSH adds to head)
	if list.Head.Value != "world" {
		t.Errorf("Expected 'world' at head, got %s", list.Head.Value)
	}

	if list.Head.Next.Value != "hello" {
		t.Errorf("Expected 'hello' second, got %s", list.Head.Next.Value)
	}
}

func TestRPush(t *testing.T) {
	store := NewStore()

	list, err := store.GetOrCreateList("mylist")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Push one value
	length := list.RPush("hello")
	if length != 1 {
		t.Errorf("Expected length 1, got %d", length)
	}

	if list.Tail.Value != "hello" {
		t.Errorf("Expected 'hello', got %s", list.Tail.Value)
	}

	// Push another value
	length = list.RPush("world")
	if length != 2 {
		t.Errorf("Expected length 2, got %d", length)
	}

	// Check order (RPUSH adds to tail)
	if list.Head.Value != "hello" {
		t.Errorf("Expected 'hello' at head, got %s", list.Head.Value)
	}

	if list.Tail.Value != "world" {
		t.Errorf("Expected 'world' at tail, got %s", list.Tail.Value)
	}
}

func TestTypeCommand(t *testing.T) {
	store := NewStore()

	// String type
	store.SetString("strkey", "value")
	if store.Type("strkey") != "string" {
		t.Error("Expected type 'string'")
	}

	// List type
	list, _ := store.GetOrCreateList("listkey")
	list.LPush("item")
	if store.Type("listkey") != "list" {
		t.Error("Expected type 'list'")
	}

	// Non-existent key
	if store.Type("nonexistent") != "none" {
		t.Error("Expected type 'none'")
	}
}

func TestWRONGTYPE(t *testing.T) {
	store := NewStore()

	// Set a string
	store.SetString("mykey", "value")

	// Try to use as list
	_, err := store.GetOrCreateList("mykey")
	if err == nil {
		t.Error("Expected WRONGTYPE error")
	}

	if err.Error() != "WRONGTYPE" {
		t.Errorf("Expected 'WRONGTYPE', got %s", err.Error())
	}
}

func TestWRONGTYPE_ListAsString(t *testing.T) {
	store := NewStore()

	// Create a list
	list, _ := store.GetOrCreateList("mykey")
	list.LPush("item")

	// Try to use as string
	_, _, err := store.GetString("mykey")
	if err == nil {
		t.Error("Expected WRONGTYPE error")
	}

	if err.Error() != "WRONGTYPE" {
		t.Errorf("Expected 'WRONGTYPE', got %s", err.Error())
	}
}

func TestLPushMultipleValues(t *testing.T) {
	store := NewStore()
	list, _ := store.GetOrCreateList("mylist")

	// Push multiple values at once
	length := list.LPush("a", "b", "c")

	if length != 3 {
		t.Errorf("Expected length 3, got %d", length)
	}

	// Check order: c, b, a (reverse of input)
	expected := []string{"c", "b", "a"}
	node := list.Head
	for i, exp := range expected {
		if node == nil {
			t.Fatalf("List too short at index %d", i)
		}
		if node.Value != exp {
			t.Errorf("At index %d: expected %s, got %s", i, exp, node.Value)
		}
		node = node.Next
	}
}

func TestRPushMultipleValues(t *testing.T) {
	store := NewStore()
	list, _ := store.GetOrCreateList("mylist")

	// Push multiple values at once
	length := list.RPush("a", "b", "c")

	if length != 3 {
		t.Errorf("Expected length 3, got %d", length)
	}

	// Check order: a, b, c (same as input)
	expected := []string{"a", "b", "c"}
	node := list.Head
	for i, exp := range expected {
		if node == nil {
			t.Fatalf("List too short at index %d", i)
		}
		if node.Value != exp {
			t.Errorf("At index %d: expected %s, got %s", i, exp, node.Value)
		}
		node = node.Next
	}
}

// =============================================================================
// List Command Handler Tests
// =============================================================================

func TestHandleCommand_LPush(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"LPUSH", "mylist", "value1"})

	if result != ":1\r\n" {
		t.Errorf("expected ':1\\r\\n', got %q", result)
	}

	// Push another
	result = handleCommand([]string{"LPUSH", "mylist", "value2"})
	if result != ":2\r\n" {
		t.Errorf("expected ':2\\r\\n', got %q", result)
	}
}

func TestHandleCommand_LPushMultiple(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"LPUSH", "mylist", "a", "b", "c"})

	if result != ":3\r\n" {
		t.Errorf("expected ':3\\r\\n', got %q", result)
	}
}

func TestHandleCommand_LPushMissingArgs(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"LPUSH", "mylist"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_LPushWrongType(t *testing.T) {
	globalStore = NewStore()
	globalStore.SetString("strkey", "value")

	result := handleCommand([]string{"LPUSH", "strkey", "item"})

	if !strings.HasPrefix(result, "-WRONGTYPE") {
		t.Errorf("expected WRONGTYPE error, got %q", result)
	}
}

func TestHandleCommand_RPush(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"RPUSH", "mylist", "value1"})

	if result != ":1\r\n" {
		t.Errorf("expected ':1\\r\\n', got %q", result)
	}

	result = handleCommand([]string{"RPUSH", "mylist", "value2"})
	if result != ":2\r\n" {
		t.Errorf("expected ':2\\r\\n', got %q", result)
	}
}

func TestHandleCommand_RPushMissingArgs(t *testing.T) {
	globalStore = NewStore()

	result := handleCommand([]string{"RPUSH", "mylist"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_Type(t *testing.T) {
	globalStore = NewStore()

	// Non-existent
	result := handleCommand([]string{"TYPE", "nokey"})
	if result != "+none\r\n" {
		t.Errorf("expected '+none\\r\\n', got %q", result)
	}

	// String
	globalStore.SetString("strkey", "value")
	result = handleCommand([]string{"TYPE", "strkey"})
	if result != "+string\r\n" {
		t.Errorf("expected '+string\\r\\n', got %q", result)
	}

	// List
	handleCommand([]string{"LPUSH", "listkey", "item"})
	result = handleCommand([]string{"TYPE", "listkey"})
	if result != "+list\r\n" {
		t.Errorf("expected '+list\\r\\n', got %q", result)
	}
}

func TestHandleCommand_TypeMissingArgs(t *testing.T) {
	result := handleCommand([]string{"TYPE"})

	if !strings.HasPrefix(result, "-ERR") {
		t.Errorf("expected error response, got %q", result)
	}
}

func TestHandleCommand_GetOnList(t *testing.T) {
	globalStore = NewStore()
	handleCommand([]string{"LPUSH", "listkey", "item"})

	result := handleCommand([]string{"GET", "listkey"})

	if !strings.HasPrefix(result, "-WRONGTYPE") {
		t.Errorf("expected WRONGTYPE error, got %q", result)
	}
}

// =============================================================================
// List Doubly-Linked Structure Tests
// =============================================================================

func TestList_DoublyLinked(t *testing.T) {
	list := NewList()
	list.RPush("a", "b", "c")

	// Forward traversal
	forward := []string{}
	for node := list.Head; node != nil; node = node.Next {
		forward = append(forward, node.Value)
	}

	if len(forward) != 3 || forward[0] != "a" || forward[1] != "b" || forward[2] != "c" {
		t.Errorf("Forward traversal failed: %v", forward)
	}

	// Backward traversal
	backward := []string{}
	for node := list.Tail; node != nil; node = node.Prev {
		backward = append(backward, node.Value)
	}

	if len(backward) != 3 || backward[0] != "c" || backward[1] != "b" || backward[2] != "a" {
		t.Errorf("Backward traversal failed: %v", backward)
	}
}

func TestList_HeadTailConsistency(t *testing.T) {
	list := NewList()

	// Single element
	list.LPush("only")
	if list.Head != list.Tail {
		t.Error("Head and Tail should be same for single element")
	}

	// Add more
	list.LPush("first")
	if list.Head == list.Tail {
		t.Error("Head and Tail should differ for multiple elements")
	}

	if list.Head.Value != "first" {
		t.Errorf("Expected head 'first', got %s", list.Head.Value)
	}
	if list.Tail.Value != "only" {
		t.Errorf("Expected tail 'only', got %s", list.Tail.Value)
	}
}
