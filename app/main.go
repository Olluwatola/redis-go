package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)


// Store provides thread-safe in-memory key-value storage
type Store struct {
    mu   sync.RWMutex
    data map[string]string
    expires map[string]time.Time
}

// NewStore creates a new Store
func NewStore() *Store {
    return &Store{
        data: make(map[string]string),
        expires: make(map[string]time.Time),
    }
}

// Set stores a key-value pair (thread-safe)
func (s *Store) Set(key, value string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.data[key] = value
    delete(s.expires, key)
}

// SetWithExpire sets a key with expiration time
func (s *Store) SetWithExpire(key, value string, expireAt time.Time) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.data[key] = value
    s.expires[key] = expireAt
}

// Get retrieves a value by key (thread-safe)
// Get retrieves a value, checking expiration
func (s *Store) Get(key string) (string, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Lazy expiration check
    if s.isExpiredLocked(key) {
        delete(s.data, key)
        delete(s.expires, key)
        return "", false
    }

    value, exists := s.data[key]
    return value, exists
}

// isExpiredLocked checks if key is expired (must hold lock)
func (s *Store) isExpiredLocked(key string) bool {
    expireAt, hasExpiry := s.expires[key]
    if !hasExpiry {
        return false
    }

    return time.Now().After(expireAt)
}

// StartActiveExpiration begins background cleanup
func (s *Store) StartActiveExpiration(interval time.Duration) {
    ticker := time.NewTicker(interval)

    go func() {
        defer ticker.Stop()

        for range ticker.C {
            s.deleteExpiredKeys()
        }
    }()
}

// deleteExpiredKeys performs active expiration
func (s *Store) deleteExpiredKeys() {
    const (
        SAMPLE_SIZE   = 20
        EXPIRED_RATIO = 0.25
    )

    s.mu.Lock()
    defer s.mu.Unlock()

    now := time.Now()
    sampled := 0
    deleted := 0

    // Sample random keys
    for key, expireAt := range s.expires {
        if sampled >= SAMPLE_SIZE {
            break
        }
        sampled++

        if now.After(expireAt) {
            delete(s.data, key)
            delete(s.expires, key)
            deleted++
        }
    }

    // Could implement iterative cleanup here if ratio > 25%
}

// Keys returns all keys in the store (thread-safe)
func (s *Store) Keys() []string {
    s.mu.RLock()
    defer s.mu.RUnlock()

    keys := make([]string, 0, len(s.data))
    for k := range s.data {
        keys = append(keys, k)
    }
    return keys
}

// Size returns number of keys (thread-safe)
func (s *Store) Size() int {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return len(s.data)
}

// Global store shared by all connections
var globalStore = NewStore()

func main() {

    // Start active expiration every 100ms
    globalStore.StartActiveExpiration(100 * time.Millisecond)

	fmt.Println("Redis server starting...")


	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	//Accept connections in a loop
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		//handle each connection in a goroutine
		fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer fmt.Printf("Client %s disconnected\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	
	for {
		line, err := parseRESP(reader)
		if err != nil {
			if(err == io.EOF){
				fmt.Println("Client closed the connection.")
				break
			}
			fmt.Println("Error reading line: ", err.Error())
			break
		}
		
		fmt.Printf("Received command: %v\n", line)
		response := handleCommand(line)
		fmt.Printf("Sending response: %q\n", response)
		
		// send response back to client
		_,err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing response: ", err.Error())
			break
		}
	}
}

func parseRESP(reader *bufio.Reader) ([]string, error) {
    // Read first byte to determine type
    typeByte, err := reader.ReadByte()
    if err != nil {
        return nil, err
    }
    
    switch typeByte {
    case '*': // Array
        return parseArray(reader)
    default:
        return nil, fmt.Errorf("unknown RESP type: %c", typeByte)
    }
}

func parseArray(reader *bufio.Reader) ([]string, error) {
    // Read array length
    line, err := reader.ReadString('\n')
    if err != nil {
        return nil, err
    }
    
    line = strings.TrimSpace(line)
    var count int
    _, err = fmt.Sscanf(line, "%d", &count)
    if err != nil {
        return nil, fmt.Errorf("invalid array count: %v", err)
    }
    
    // Read each element
    result := make([]string, count)
    for i := 0; i < count; i++ {
        element, err := parseBulkString(reader)
        if err != nil {
            return nil, err
        }
        result[i] = element
    }
    
    return result, nil
}

func parseBulkString(reader *bufio.Reader) (string, error) {
    // Read '$'
    typeByte, err := reader.ReadByte()
    if err != nil {
        return "", err
    }
    if typeByte != '$' {
        return "", fmt.Errorf("expected bulk string, got %c", typeByte)
    }
    
    // Read length
    line, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }
    
    line = strings.TrimSpace(line)
    var length int
    _, err = fmt.Sscanf(line, "%d", &length)
    if err != nil {
        return "", fmt.Errorf("invalid bulk string length: %v", err)
    }
    
    // Read exact number of bytes
    data := make([]byte, length)
    _, err = io.ReadFull(reader, data)
    if err != nil {
        return "", err
    }
    
    // Read trailing \r\n
    reader.ReadByte() // \r
    reader.ReadByte() // \n
    
    return string(data), nil
}

func handleCommand(command []string) string {
    if len(command) == 0 {
        return "-ERR empty command\r\n"
    }
    
    cmd := strings.ToUpper(command[0])
    
    switch cmd {
    case "PING":
        return "+PONG\r\n"
    case "ECHO":
        if len(command) < 2 {
            return "-ERR wrong number of arguments\r\n"
        }
        return fmt.Sprintf("$%d\r\n%s\r\n", len(command[1]), command[1])
    case "SET":
        return handleSet(command)
    case "GET":
        return handleGet(command)

    default:
        return fmt.Sprintf("-ERR unknown command '%s'\r\n", cmd)
    }

    
}


func handleSet(command []string) string {
    if len(command) < 3 {
        return "-ERR wrong number of arguments for 'set' command\r\n"
    }

    key := command[1]
    value := command[2]
    var expireAt time.Time

    // Parse optional EX/PX flags
    i := 3
    for i < len(command) {
        flag := strings.ToUpper(command[i])

        switch flag {
        case "EX":
            if i+1 >= len(command) {
                return "-ERR syntax error\r\n"
            }
            seconds, err := strconv.ParseInt(command[i+1], 10, 64)
            if err != nil || seconds <= 0 {
                return "-ERR value is not an integer or out of range\r\n"
            }
            expireAt = time.Now().Add(time.Duration(seconds) * time.Second)
            i += 2

        case "PX":
            if i+1 >= len(command) {
                return "-ERR syntax error\r\n"
            }
            milliseconds, err := strconv.ParseInt(command[i+1], 10, 64)
            if err != nil || milliseconds <= 0 {
                return "-ERR value is not an integer or out of range\r\n"
            }
            expireAt = time.Now().Add(time.Duration(milliseconds) * time.Millisecond)
            i += 2

        default:
            return fmt.Sprintf("-ERR unsupported option '%s'\r\n", flag)
        }
    }

    // Store with or without expiration
    if expireAt.IsZero() {
        globalStore.Set(key, value)
    } else {
        globalStore.SetWithExpire(key, value, expireAt)
    }

    return "+OK\r\n"
}

func handleGet(command []string) string {
    if len(command) < 2 {
        return "-ERR wrong number of arguments for 'get' command\r\n"
    }

    key := command[1]

    // Retrieve from global store (thread-safe)
    value, exists := globalStore.Get(key)

    if !exists {
        // Key doesn't exist - return null bulk string
        return "$-1\r\n"
    }


    // Key exists - return value as bulk string
    return encodeBulkString(value)
}

// Encode a string as RESP bulk string
func encodeBulkString(s string) string {
    return fmt.Sprintf("$%d\r\n%s\r\n", len(s), s)
}

// Encode a simple string
func encodeSimpleString(s string) string {
    return fmt.Sprintf("+%s\r\n", s)
}

// Encode an error
func encodeError(msg string) string {
    return fmt.Sprintf("-%s\r\n", msg)
}