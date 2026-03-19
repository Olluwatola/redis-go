
Golang Project 
Redis Server Challenge: Complete Stage-by-Stage Timeline
Overview
This document provides a comprehensive breakdown of every stage in building a Redis server from scratch. Each stage builds upon previous knowledge, gradually increasing in complexity from basic TCP connectivity to advanced features like geospatial commands and authentication.
Foundation Stages
Stage 1: Bind to a Port
What You'll Build:
Your first task is to create a TCP server that successfully binds to port 6379 (Redis's default port) and listens for incoming connections. This is the fundamental building block of any network server.
Technical Concepts:
Understanding TCP/IP networking fundamentals
Socket programming basics: creating sockets, binding to addresses, and listening
The concept of a server accepting connections on a specific port
How operating systems manage network ports and connections
What Success Looks Like:
Your server starts without errors
It successfully binds to port 6379
It enters a listening state, ready to accept connections
You can verify the port is open using tools like telnet or netcat
The server doesn't crash when clients attempt to connect
Basic logging shows the server is running and waiting for connections
Key Challenges:
Handling port conflicts if another process is using 6379
Understanding address binding (localhost vs all interfaces)
Proper error handling for socket operations
Ensuring clean shutdown and resource cleanup
Stage 2: Respond to PING
What You'll Build:
Implement the ability to accept a client connection, receive a PING command, and respond with PONG. This introduces you to the Redis Protocol (RESP).
Technical Concepts:
The Redis Serialization Protocol (RESP) format
Reading data from TCP sockets
Parsing simple RESP arrays and bulk strings
Writing formatted responses back to clients
The request-response cycle
What Success Looks Like:
A client can connect to your server
When the client sends a PING command (formatted as a RESP array), your server correctly parses it
Your server responds with "+PONG\r\n" (a RESP simple string)
The connection remains open after the exchange
You can use redis-cli or a telnet session to manually test this
Proper handling of connection lifecycle (accept, read, write, potentially close)
Key Challenges:
Understanding RESP protocol formatting with \r\n delimiters
Correctly parsing the incoming byte stream
Formatting the response according to RESP specifications
Handling potential network errors during read/write operations
Stage 3: Respond to Multiple PINGs
What You'll Build:
Enhance your server to handle multiple sequential PING commands from the same client without closing the connection.
Technical Concepts:
Persistent TCP connections
Request loops: continuously reading from a socket
Buffer management for incoming data
Knowing when a complete command has been received
What Success Looks Like:
A single client connects and sends multiple PING commands in succession
Your server responds with PONG to each PING
The connection stays alive between commands
No memory leaks from accumulating buffers
The server can gracefully detect when a client disconnects
Clean handling of connection closure
Key Challenges:
Implementing a loop that reads commands until the client disconnects
Handling partial reads (when a command arrives in multiple TCP packets)
Detecting end-of-command boundaries in the RESP protocol
Managing buffers efficiently without memory leaks
Recognizing client disconnection vs waiting for more data
Stage 4: Handle Concurrent Clients
What You'll Build:
Modify your server to handle multiple clients simultaneously. This is where your server becomes truly useful in real-world scenarios.
Technical Concepts:
Concurrency models (threading, async/await, event loops)
Thread-per-connection vs event-driven architectures
Shared state and race conditions
Resource management across multiple connections
What Success Looks Like:
Multiple redis-cli instances can connect simultaneously
Each client can send PING commands independently
Responses go to the correct client
The server doesn't block while handling one client
All clients receive prompt responses
Clean handling when one client disconnects doesn't affect others
No race conditions or data corruption
Key Challenges:
Choosing an appropriate concurrency model for your language
Ensuring thread safety or proper async handling
Managing per-client state separately
Avoiding blocking operations that prevent other clients from being served
Proper cleanup of resources when clients disconnect
Testing with multiple simultaneous connections
Stage 5: Implement the ECHO Command
What You'll Build:
Add support for the ECHO command, which takes a string argument and returns it back to the client.
Technical Concepts:
Parsing commands with arguments
Extracting data from RESP bulk strings
Routing different command types to appropriate handlers
Building a command dispatch system
What Success Looks Like:
Clients can send ECHO followed by any string
The server responds with that exact string in RESP format
PING still works correctly (command routing is in place)
Proper handling of edge cases: empty strings, special characters, large strings
Error messages for malformed ECHO commands
The response format matches Redis's RESP encoding for bulk strings
Key Challenges:
Parsing RESP arrays to extract command name and arguments
Distinguishing between different command types
Handling variable numbers of arguments
Encoding the response as a RESP bulk string with correct length prefix
Building a maintainable command dispatcher pattern
Validating argument counts
Stage 6: Implement the SET & GET Commands
What You'll Build:
Create an in-memory key-value store with SET (to store data) and GET (to retrieve data) commands. This is the core functionality of Redis as a cache.
Technical Concepts:
Hash tables or dictionary data structures
Key-value storage and retrieval
Data persistence in memory (not disk, yet)
String encoding and storage
What Success Looks Like:
SET key value stores data in memory
GET key retrieves the previously stored value
GET on a non-existent key returns null
Values can be overwritten with subsequent SET commands
Multiple key-value pairs can coexist
Concurrent clients can SET and GET without conflicts
RESP formatting is correct for both commands
Key Challenges:
Choosing an efficient data structure for the key-value store
Thread-safe access to the shared data store (if using threading)
Properly encoding null responses for GET on missing keys
Memory management as data accumulates
Handling special characters in keys and values
Case sensitivity considerations
Stage 7: Expiry
What You'll Build:
Add time-to-live (TTL) functionality where keys automatically expire and are removed after a specified duration.
Technical Concepts:
Time-based operations and scheduling
The SET command with PX (milliseconds) or EX (seconds) flags
Background cleanup of expired keys
Timestamp comparisons and time arithmetic
What Success Looks Like:
SET key value PX 1000 sets a key that expires after 1000 milliseconds
GET on an expired key returns null (as if the key never existed)
Expired keys are eventually cleaned up from memory
Non-expired keys remain accessible
Multiple keys can have different expiration times
Edge cases handled: very short TTLs, very long TTLs, negative values
Key Challenges:
Storing expiration metadata alongside values
Checking expiration on every GET operation
Implementing background cleanup (lazy deletion vs active deletion)
Handling time zone and clock issues
Parsing PX and EX flags from SET command
Memory cleanup of expired data
Race conditions between expiration checks
Lists Section
Stage 8: Create a List
What You'll Build:
Implement the LPUSH command to create a new Redis list data structure and add initial elements to it.
Technical Concepts:
List data structures (doubly-linked lists or array-based)
Type differentiation in your key-value store
List initialization
Left-side insertion semantics
What Success Looks Like:
LPUSH mylist value creates a new list with one element
Lists are a distinct type from strings
TYPE command can distinguish lists from strings
Setting a string key and list key with same name handles type conflicts appropriately
The list maintains insertion order
Key Challenges:
Extending your storage system to support multiple data types
Type checking when commands operate on existing keys
Implementing efficient list data structures
Memory allocation for variable-length lists
Returning appropriate error messages for type mismatches
Stage 9: Append an Element
What You'll Build:
Extend LPUSH to add elements to an existing list, inserting at the head (left side).
Technical Concepts:
List prepending operations
Maintaining list order
Return values indicating list length after operation
What Success Looks Like:
LPUSH on an existing list adds the new element to the front
The operation returns the new list length as an integer
Multiple LPUSH operations result in reverse insertion order
Large lists handle additional elements efficiently
The original list elements remain intact
Key Challenges:
Efficient insertion at the head of your list structure
Updating list metadata (size/length)
Maintaining pointer integrity in linked lists
Handling very large lists without performance degradation
Stage 10: Append Multiple Elements
What You'll Build:
Allow LPUSH to accept multiple values in a single command, inserting them atomically.
Technical Concepts:
Variadic command arguments
Atomic multi-element operations
Argument parsing for variable-length commands
What Success Looks Like:
LPUSH mylist a b c inserts all three elements
Elements are inserted in the order specified in the arguments
The command is atomic (all or nothing)
Returns the final list length
Works correctly with any number of elements
Key Challenges:
Parsing variable numbers of arguments from RESP
Maintaining atomicity (all insertions succeed together)
Preserving the correct insertion order
Efficient batch insertion
Memory allocation for multiple elements
Stage 11: List Elements (Positive Indexes)
What You'll Build:
Implement LRANGE to retrieve a range of elements from a list using positive (zero-based) indexes.
Technical Concepts:
Array indexing conventions
Range queries on ordered data
Inclusive range semantics (both start and end included)
RESP array encoding for multiple values
What Success Looks Like:
LRANGE mylist 0 2 returns the first three elements
LRANGE mylist 1 1 returns just the second element
Out-of-range indexes are handled gracefully
Empty ranges return empty arrays
The response is a RESP array of bulk strings
Key Challenges:
Efficient range extraction from your list structure
Boundary checking for index validation
Handling edge cases: empty lists, single-element ranges, full list ranges
Properly encoding multi-element responses in RESP
Iterator or index-based access depending on your list implementation
Stage 12: List Elements (Negative Indexes)
What You'll Build:
Extend LRANGE to support negative indexes, which count from the end of the list.
Technical Concepts:
Negative indexing semantics (-1 is last element, -2 is second-to-last)
Index normalization
Mixed positive and negative indexes in ranges
What Success Looks Like:
LRANGE mylist -3 -1 returns the last three elements
LRANGE mylist 0 -1 returns the entire list
LRANGE mylist -2 -2 returns the second-to-last element
Negative indexes correctly convert to positive equivalents
All combinations of positive and negative indexes work correctly
Key Challenges:
Converting negative indexes to positive positions
Handling lists where negative index magnitude exceeds list length
Edge cases with very short lists
Ensuring reversed ranges return empty arrays
Maintaining consistency with Redis's exact behavior
Stage 13: Prepend Elements
What You'll Build:
Implement RPUSH to add elements to the right side (tail) of a list.
Technical Concepts:
Tail insertion operations
Dual-ended list capabilities
Symmetric operations to LPUSH
What Success Looks Like:
RPUSH mylist value appends to the end of the list
Supports single and multiple values
Returns the new list length
Works on new and existing lists
LPUSH and RPUSH can be used on the same list
Key Challenges:
Efficient tail insertion depending on your list structure
Maintaining pointers to both head and tail
Ensuring consistent behavior with LPUSH
Performance considerations for append operations
Stage 14: Query List Length
What You'll Build:
Implement LLEN to return the number of elements in a list.
Technical Concepts:
Metadata storage
Constant-time length queries
Handling non-existent keys
What Success Looks Like:
LLEN mylist returns the correct count
Returns 0 for non-existent keys
Returns error for non-list types
Updates correctly as list is modified
Very fast operation regardless of list size
Key Challenges:
Maintaining an accurate length counter
Updating length on every modification
Type checking before returning length
Integer encoding in RESP responses
Stage 15: Remove an Element
What You'll Build:
Implement LPOP to remove and return the first element from the left (head) of a list.
Technical Concepts:
Destructive read operations
List modification with return values
Empty list handling
What Success Looks Like:
LPOP mylist removes and returns the head element
List length decreases by one
Repeated LPOP operations retrieve elements in FIFO order
LPOP on empty list returns null
List is removed when last element is popped
Type errors on non-list keys
Key Challenges:
Atomic remove-and-return operation
Proper memory cleanup of removed elements
Updating list metadata
Removing the list key when empty
Thread-safe modification
Stage 16: Remove Multiple Elements
What You'll Build:
Extend LPOP to accept a count argument for removing multiple elements at once.
Technical Concepts:
Batch removal operations
Returning multiple values
Partial success handling
What Success Looks Like:
LPOP mylist 3 removes and returns three elements
Returns available elements if fewer than count exist
Returns empty array if list is empty
Returns RESP array with removed elements
List is cleaned up if fully drained
Key Challenges:
Handling count larger than list size
Efficient multi-element removal
Array response formatting
Maintaining atomicity
Edge case: count of zero or negative
Stage 17: Blocking Retrieval
What You'll Build:
Implement BLPOP for blocking left-side pop operations that wait for elements if the list is empty.
Technical Concepts:
Blocking operations and timeouts
Client waiting/suspension
Event notification when data becomes available
Multi-client coordination
What Success Looks Like:
BLPOP mylist 0 blocks until an element is available (0 means wait forever)
When another client executes LPUSH, the blocked client immediately receives the element
Returns the key name and value
Can block on multiple lists, returning from first available
Proper client unblocking when data arrives
Key Challenges:
Implementing blocking waits without busy-waiting
Queuing blocked clients
Waking up correct clients when data arrives
Handling disconnection of blocked clients
Fair ordering of multiple blocked clients
Testing synchronization between clients
Stage 18: Blocking Retrieval with Timeout
What You'll Build:
Extend BLPOP to support timeout values, returning null if timeout expires before data arrives.
Technical Concepts:
Timeout-based blocking
Timer management
Graceful timeout handling
What Success Looks Like:
BLPOP mylist 5 waits up to 5 seconds
Returns null if timeout expires
Returns data immediately if available
Accurate timeout duration (not significantly early or late)
Client remains responsive after timeout
Key Challenges:
Implementing accurate timers
Cleaning up timed-out waiters
Race conditions between timeout and data arrival
Precision requirements for timeout values
Resource cleanup after timeout
Streams Section
Stage 19: The TYPE Command
What You'll Build:
Implement the TYPE command that returns the data type of a key's value.
Technical Concepts:
Runtime type introspection
Type system in your Redis implementation
String representations of types
What Success Looks Like:
TYPE mykey returns "string", "list", "stream", etc.
TYPE nonexistent returns "none"
Correctly identifies all supported types
Response is a simple string in RESP format
Key Challenges:
Maintaining type information with stored values
Supporting new types as you add features
Consistent type naming
Efficient type lookup
Stage 20: Create a Stream
What You'll Build:
Implement XADD to create a new stream and add the first entry with an auto-generated or specified ID.
Technical Concepts:
Stream data structures (append-only logs)
Entry IDs with millisecond-timestamp and sequence number format
Field-value pairs in stream entries
Stream initialization
What Success Looks like:
XADD mystream * field1 value1 creates a stream and returns generated ID
Entry IDs follow format: timestamp-sequence (e.g., 1609459200000-0)
Multiple field-value pairs can be stored in one entry
Stream maintains insertion order
TYPE mystream returns "stream"
Key Challenges:
Implementing the ID generation algorithm
Storing variable field-value pairs per entry
Efficient append operations
Timestamp precision and clock handling
Data structure selection for ordered entries
Stage 21: Validating Entry IDs
What You'll Build:
Add validation for explicitly specified entry IDs to ensure they're monotonically increasing.
Technical Concepts:
ID ordering and validation
Error handling for invalid IDs
ID comparison logic
What Success Looks Like:
XADD mystream 1234567890000-0 field value accepts valid IDs
Rejects IDs that aren't greater than the last ID
Proper error messages for invalid formats
IDs must be chronologically increasing
Sequence numbers must increase when timestamps are equal
Key Challenges:
Parsing ID format correctly
Implementing ID comparison (lexicographic with numeric components)
Clear error messages
Edge cases: first entry, equal timestamps
Validating both timestamp and sequence components
Stage 22: Partially Auto-generated IDs
What You'll Build:
Support IDs where timestamp is specified but sequence is auto-generated (format: timestamp-*).
Technical Concepts:
Hybrid manual/automatic ID generation
Sequence number management within a timestamp
ID parsing with wildcard
What Success Looks Like:
XADD mystream 1234567890000-* field value generates sequence automatically
Sequence increments correctly for same timestamp
Resets to 0 for new timestamps
Validation still ensures monotonic increase
Handles edge case where timestamp equals last but sequence needs incrementing
Key Challenges:
Tracking last ID to generate next sequence
Handling transition between timestamps
Sequence overflow considerations
Parsing asterisk in ID position
Ensuring generated IDs maintain ordering
Stage 23: Fully Auto-generated IDs
What You'll Build:
Complete auto-generation using * for entire ID, using current system time.
Technical Concepts:
System time retrieval with millisecond precision
Clock monotonicity handling
Automatic timestamp and sequence generation
What Success Looks Like:
XADD mystream * field value uses current time
Handles multiple entries in same millisecond with incrementing sequences
Clock going backwards handled gracefully
Generated IDs always increase
Works reliably under high insertion rates
Key Challenges:
Getting precise timestamps
Handling rapid insertions (same timestamp)
Clock synchronization issues
Sequence overflow within a timestamp
Testing time-based logic
Stage 24: Query Entries from Stream
What You'll Build:
Implement XRANGE to retrieve entries within an ID range.
Technical Concepts:
Range queries on ordered data
Inclusive range boundaries
Stream traversal
Multi-entry response formatting
What Success Looks Like:
XRANGE mystream start end returns matching entries
Each entry includes ID and field-value pairs
Results maintain insertion order
Handles empty ranges
Efficient retrieval without scanning entire stream
Key Challenges:
Efficient range search in your stream structure
Complex RESP encoding for nested arrays (entries with field-value pairs)
Boundary inclusion/exclusion logic
Large range performance
Memory efficiency for result sets
Stage 25: Query with -
What You'll Build:
Support the special - symbol meaning "start from the beginning" in XRANGE.
Technical Concepts:
Special range boundary markers
Stream minimum ID semantics
What Success Looks Like:
XRANGE mystream - + returns all entries
XRANGE mystream - someID returns entries from start to someID
Handles empty streams
Correctly interprets - as smallest possible ID
Key Challenges:
Parsing special symbols in range arguments
Representing minimum ID in your system
Consistent handling across different queries
Stage 26: Query with +
What You'll Build:
Support the special + symbol meaning "until the end" in XRANGE.
Technical Concepts:
Maximum ID semantics
Open-ended ranges
What Success Looks Like:
XRANGE mystream someID + returns entries from someID to end
XRANGE mystream - + returns entire stream
Works with dynamic streams (new entries don't affect ongoing query)
Key Challenges:
Representing maximum ID
Snapshot consistency for queries
Efficient "to end" traversal
Stage 27: Query Single Stream Using XREAD
What You'll Build:
Implement XREAD to read new entries from a stream starting after a specified ID.
Technical Concepts:
Tail-reading patterns
Last-seen ID tracking
Different semantics from XRANGE
What Success Looks Like:
XREAD STREAMS mystream lastID returns entries after lastID
Returns only newer entries, not including the ID itself
Returns empty if no new entries
Multiple entries returned in order
Key Challenges:
Exclusive start boundary (after, not including)
Different response format than XRANGE
Multiple streams support preparation
Handling non-existent IDs
Stage 28: Query Multiple Streams Using XREAD
What You'll Build:
Extend XREAD to query multiple streams simultaneously, each with its own starting ID.
Technical Concepts:
Multi-stream operations
Parallel range queries
Grouped results
What Success Looks Like:
XREAD STREAMS stream1 stream2 id1 id2 returns results from both
Each stream section shows new entries
Order maintained within each stream
Empty streams handled gracefully
Response groups results by stream
Key Challenges:
Parsing multiple streams and IDs
Complex nested response structure
Efficient multi-stream query
Maintaining per-stream context
Response formatting complexity
Stage 29: Blocking Reads
What You'll Build:
Add BLOCK option to XREAD to wait for new entries if none are currently available.
Technical Concepts:
Blocking stream reads
Event-driven notification
Timeout handling for streams
Multi-client stream waiting
What Success Looks Like:
XREAD BLOCK 1000 STREAMS mystream lastID waits up to 1000ms
Returns immediately if entries exist
Waits and returns when new entry added via XADD
Timeout returns null if no data arrives
Multiple blocked readers all receive new entries
Key Challenges:
Notifying waiting clients when XADD occurs
Managing timeouts accurately
Queue of blocked clients per stream
Race conditions between check and block
Cleanup on client disconnect
Stage 30: Blocking Reads Without Timeout
What You'll Build:
Support BLOCK 0 for indefinite waiting until data arrives.
Technical Concepts:
Infinite blocking
Resource management for long-lived waits
Graceful shutdown with blocked clients
What Success Looks Like:
XREAD BLOCK 0 STREAMS mystream id waits forever
Returns as soon as data arrives
Survives server load and delays
Properly cleaned up on client disconnect
No resource leaks from long waits
Key Challenges:
Preventing resource exhaustion
Tracking indefinitely blocked clients
Testing long-duration blocks
Cleanup mechanisms
Stage 31: Blocking Reads Using $
What You'll Build:
Support the special $ ID meaning "wait for entries newer than what currently exists".
Technical Concepts:
Dynamic ID resolution
"New data only" semantics
Snapshot at blocking time
What Success Looks Like:
XREAD BLOCK 1000 STREAMS mystream $ waits only for brand new entries
Ignores all existing entries
$ is resolved to stream's current maximum ID at block time
Commonly used for consumer patterns
Key Challenges:
Resolving $ to current max ID at block invocation
Differentiating from regular ID blocking
Edge case: empty stream with $
Timing of ID resolution vs actual blocking
Transactions Section
Stage 32: The INCR Command (1/3)
What You'll Build:
Implement basic INCR command that increments an integer value stored at a key.
Technical Concepts:
Atomic read-modify-write operations
Type coercion and validation
Integer parsing and arithmetic
What Success Looks Like:
INCR counter increases value by 1
Creates key with value 1 if it doesn't exist
Returns the new value
Works with any valid integer string
Thread-safe increment operation
Key Challenges:
Parsing string values as integers
Handling non-integer values (error)
Atomic increment (no race conditions)
Integer overflow considerations
Type checking
Stage 33: The INCR Command (2/3)
What You'll Build:
Handle edge cases and errors for INCR.
Technical Concepts:
Error handling and reporting
Type errors vs value errors
Error message formatting in RESP
What Success Looks Like:
INCR on non-integer string returns error
INCR on list/stream returns type error
Proper RESP error format
Clear error messages
No data corruption on errors
Key Challenges:
Distinguishing error types
Preserving data on failed operations
RESP error encoding
Comprehensive error coverage
Stage 34: The INCR Command (3/3)
What You'll Build:
Optimize INCR and ensure full compatibility.
Technical Concepts:
Performance optimization
Edge case handling
Numerical limits
What Success Looks Like:
Very fast increment operations
Handles maximum integer values
Consistent behavior with Redis
Works under concurrent load
Proper overflow handling
Key Challenges:
Performance tuning
Integer boundary testing
Stress testing concurrent increments
Exact Redis compatibility
Stage 35: The MULTI Command
What You'll Build:
Implement MULTI to begin a transaction block.
Technical Concepts:
Transaction state management
Client session state
Command queuing initiation
What Success Looks Like:
MULTI puts client into transaction mode
Returns "OK"
Subsequent commands are queued, not executed
Client state tracked per connection
Can only be in one transaction at a time
Key Challenges:
Per-client state tracking
State transitions and validation
Handling MULTI within MULTI
Nested transaction prevention
Stage 36: The EXEC Command
What You'll Build:
Implement EXEC to execute queued transaction commands atomically.
Technical Concepts:
Atomic batch execution
Command queue processing
Transaction commit semantics
What Success Looks Like:
EXEC runs all queued commands in order
Returns array of results
All commands execute atomically
State returns to normal after EXEC
Cannot EXEC without prior MULTI
Key Challenges:
Ensuring true atomicity
Collecting results from multiple commands
Error handling during execution
Resetting transaction state
Isolation from other clients
Stage 37: Empty Transaction
What You'll Build:
Handle transactions with no queued commands.
Technical Concepts:
Degenerate case handling
Empty queue processing
What Success Looks Like:
MULTI followed immediately by EXEC returns empty array
No errors occur
State cleaned up properly
Works correctly every time
Key Challenges:
Empty result array formatting
State cleanup with no commands
Not treating empty as error
Stage 38: Queueing Commands
What You'll Build:
Properly queue commands between MULTI and EXEC without executing them.
Technical Concepts:
Command buffering
Deferred execution
QUEUED responses
What Success Looks Like:
Commands after MULTI return "QUEUED"
Commands are not actually executed yet
Command queue maintained in order
All command types can be queued
Queue persists until EXEC or DISCARD
Key Challenges:
Storing command arguments for later execution
Returning QUEUED instead of normal responses
Memory management for queued commands
Maintaining queue order
Supporting all command types in queue
Stage 39: Executing a Transaction
What You'll Build:
Execute all queued commands atomically with proper result collection.
Technical Concepts:
Atomic execution guarantees
Multi-command isolation
Result aggregation
What Success Looks Like:
All queued commands execute in order
No interleaving with other clients' commands
Each command result captured correctly
Array response with all results
True atomicity under concurrent load
Key Challenges:
Preventing other client commands during execution
Locking or isolation mechanisms
Handling mix of successful and error results
Performance with long transaction queues
Testing atomicity guarantees
Stage 40: The DISCARD Command
What You'll Build:
Implement DISCARD to abort a transaction and clear the queue.
Technical Concepts:
Transaction rollback
State cleanup
Resource freeing
What Success Looks Like:
DISCARD exits transaction mode
Queued commands are cleared
Returns "OK"
No commands are executed
Can start new transaction after DISCARD
Cannot DISCARD outside transaction
Key Challenges:
Complete queue cleanup
Memory deallocation
State reset
Error handling for DISCARD outside MULTI
Ensuring no partial execution
Stage 41: Failures Within Transactions
What You'll Build:
Handle command errors during transaction execution.
Technical Concepts:
Partial failure handling
Error propagation in transactions
Continue-on-error semantics
What Success Looks Like:
If a command fails during EXEC, that command's result is an error
Other commands in transaction still execute
Error doesn't abort entire transaction
Results array includes errors in correct positions
Consistent with Redis behavior
Key Challenges:
Capturing errors without stopping execution
Mixed success/error result arrays
Maintaining execution order with errors
Testing various error scenarios
Proper error formatting in result array
Stage 42: Multiple Transactions
What You'll Build:
Support multiple clients each running their own transactions concurrently.
Technical Concepts:
Per-client transaction isolation
Concurrent transaction management
State independence
What Success Looks Like:
Each client can be in independent transaction state
One client's MULTI doesn't affect another
Transactions execute without interference
Commands queued per client
EXEC executes only that client's queue
Key Challenges:
Isolating transaction state per client
Testing concurrent transactions
Ensuring atomicity across concurrent transactions
Resource management per client
No cross-transaction contamination
Replication Section
Stage 43: Configure Listening Port
What You'll Build:
Add command-line option to configure the port your Redis server listens on.
Technical Concepts:
Command-line argument parsing
Configuration management
Port specification
What Success Looks Like:
Server starts with --port 6380 (or similar flag)
Binds to specified port instead of default 6379
Works with any valid port number
Error handling for invalid ports
Documentation of configuration option
Key Challenges:
Argument parsing
Port validation
Default vs configured behavior
Port conflict handling
Stage 44: The INFO Command
What You'll Build:
Implement INFO command to return server information and statistics.
Technical Concepts:
Server introspection
Formatted text responses
Metadata collection
What Success Looks Like:
INFO returns server details
Includes role (master by default)
Shows various statistics
RESP bulk string response
Formatted as key:value lines
Key Challenges:
Gathering server statistics
Formatting multi-line response
Keeping stats updated
Determining what information to include
RESP encoding for long strings
Stage 45: The INFO Command on a Replica
What You'll Build:
Make INFO report different information when server is configured as a replica.
Technical Concepts:
Replica vs master role differentiation
Configuration flags
Runtime role tracking
What Success Looks Like:
Server started with --replicaof host port flag
INFO shows role:slave (or replica)
Shows master host and port
Different output than master INFO
Role correctly identified
Key Challenges:
Parsing replicaof configuration
Storing master connection info
Role-specific INFO formatting
Configuration validation
Stage 46: Initial Replication ID and Offset
What You'll Build:
Generate and track a replication ID and offset for the master server.
Technical Concepts:
Replication identifiers
Offset tracking
Random ID generation
What Success Looks Like:
Master generates a unique replication ID on startup
Offset starts at 0
INFO shows master_replid and master_repl_offset
ID persists during server lifetime
Offset updates as commands are processed
Key Challenges:
Generating unique IDs (typically 40-char hex)
Offset incrementation logic
Including in INFO output
Understanding replication metadata
Stage 47: Send Handshake (1/3)
What You'll Build:
Replica sends PING to master as first step of replication handshake.
Technical Concepts:
TCP client connection (replica to master)
Replication handshake protocol
Master-replica communication
What Success Looks Like:
Replica connects to configured master
Sends PING command
Receives PONG response
Connection stays open
Handshake initiation successful
Key Challenges:
Initiating outbound connection from replica
Handling connection failures
Protocol matching for handshake
Error handling and retries
Stage 48: Send Handshake (2/3)
What You'll Build:
Replica sends REPLCONF commands to exchange configuration with master.
Technical Concepts:
Configuration exchange
Capability negotiation
Multi-step handshake
What Success Looks Like:
REPLCONF listening-port sent after PING
REPLCONF capa psync2 sent to declare capabilities
Master responds with OK to each
Handshake progresses correctly
Configuration stored by master
Key Challenges:
Proper command sequencing
Parsing master responses
Handling unsupported configurations
Error recovery in handshake
Stage 49: Send Handshake (3/3)
What You'll Build:
Replica sends PSYNC command to initiate synchronization.
Technical Concepts:
Partial resynchronization protocol
Full synchronization vs partial
Replication stream initiation
What Success Looks Like:
PSYNC replicationID offset sent
For initial sync, sends PSYNC ? -1
Master responds with FULLRESYNC
Replica ready to receive data
Handshake complete
Key Challenges:
PSYNC parameter formatting
Parsing FULLRESYNC response
Extracting new replication ID and offset
State transition to syncing
Stage 50: Receive Handshake (1/2)
What You'll Build:
Master accepts and responds to replica's PING and REPLCONF commands.
Technical Concepts:
Handshake acceptance
Replica registration
Connection state tracking
What Success Looks Like:
Master responds to replica PING with PONG
Accepts REPLCONF commands with OK
Tracks replica connection state
Stores replica configuration
Maintains replica list
Key Challenges:
Identifying replica connections
Storing replica metadata
Multiple simultaneous replicas
Validating REPLCONF parameters
Stage 51: Receive Handshake (2/2)
What You'll Build:
Master responds to PSYNC with FULLRESYNC and prepares data transfer.
Technical Concepts:
Synchronization negotiation
RDB transfer preparation
Replication stream setup
What Success Looks Like:
Master sends FULLRESYNC replID offset
Prepares to send database snapshot
Replica added to active replication list
Ready to transmit data
Offset tracking begins
Key Challenges:
Generating FULLRESYNC response
Snapshot preparation
Managing replication state per replica
Thread safety with multiple replicas
Stage 52: Empty RDB Transfer
What You'll Build:
Master sends an empty RDB file to replica to complete initial synchronization.
Technical Concepts:
RDB file format basics
Binary data transfer over TCP
RESP encoding for binary data
What Success Looks Like:
Master generates empty RDB file (minimal valid RDB)
Sends RDB as RESP bulk string
Replica receives and parses RDB
Synchronization completes
Replica becomes live
Key Challenges:
Creating valid (though empty) RDB format
Binary data transmission
RDB format specification compliance
RESP bulk string with binary data
Replica RDB parsing
Stage 53: Single-Replica Propagation
What You'll Build:
Master forwards write commands to connected replica in real-time.
Technical Concepts:
Command propagation
Write command identification
Real-time replication
What Success Looks Like:
When master receives SET, it executes locally and sends to replica
Replica receives and applies command
Data stays synchronized
Replica eventually consistent with master
Offset incremented on both sides
Key Challenges:
Identifying which commands to propagate (writes only)
Asynchronous sending to replica
Handling propagation failures
Offset tracking
Replica command application
Stage 54: Multi-Replica Propagation
What You'll Build:
Master propagates commands to multiple replicas simultaneously.
Technical Concepts:
Fan-out message distribution
Multi-replica coordination
Scalable propagation
What Success Looks Like:
All connected replicas receive write commands
Each replica applies commands independently
Master tracks offset per replica
Failure of one replica doesn't affect others
Efficient broadcast mechanism
Key Challenges:
Iterating over replica connections
Parallel/asynchronous sending
Per-replica offset tracking
Handling slow/failed replicas
Resource management with many replicas
Stage 55: Command Processing
What You'll Build:
Replicas correctly process propagated commands and update local state.
Technical Concepts:
Command execution on replica
Read-only vs writable replica modes
State consistency
What Success Looks Like:
Replica applies received commands to its data store
Data matches master (eventually)
Replica updates its offset
Commands not re-propagated (no loops)
Replica serves read queries with current data
Key Challenges:
Distinguishing propagated commands from client commands
Preventing write commands from clients to replicas
Offset calculation and updates
Ensuring command application order
Testing data consistency
Stage 56: ACKs with No Commands
What You'll Build:
Implement REPLCONF GETACK command where master requests replica to acknowledge its offset.
Technical Concepts:
Acknowledgment protocol
Offset reporting
Health checking
What Success Looks Like:
Master sends REPLCONF GETACK *
Replica responds with REPLCONF ACK offset
Master receives current replica offset
Works even with no commands propagated
Enables offset tracking
Key Challenges:
Implementing REPLCONF GETACK on master
Implementing REPLCONF ACK on replica
Accurate offset reporting
Asynchronous request/response handling
Stage 57: ACKs with Commands
What You'll Build:
Track replica offset as commands are propagated and acknowledged.
Technical Concepts:
Offset calculation
Byte counting
Lag detection
What Success Looks Like:
Replica offset increases as commands received
REPLCONF ACK reports correct offset
Master can calculate replica lag
Offset matches bytes received
Accurate tracking under load
Key Challenges:
Calculating byte length of propagated commands
Offset arithmetic
Handling partial receives
Testing offset accuracy
Command size calculation
Stage 58: WAIT with No Replicas
What You'll Build:
Implement WAIT command that waits for replicas to acknowledge data, handling case of no replicas.
Technical Concepts:
Synchronous replication waiting
Replica acknowledgment aggregation
Immediate return for no replicas
What Success Looks Like:
WAIT numreplicas timeout with no replicas returns 0 immediately
Doesn't block unnecessarily
Returns replica count that acknowledged
Fast path for zero replicas
Key Challenges:
Detecting no replicas condition
Immediate return logic
Proper return value
Stage 59: WAIT with No Commands
What You'll Build:
Handle WAIT when no commands have been issued (all replicas are up to date).
Technical Concepts:
Current offset comparison
Immediate acknowledgment
What Success Looks Like:
WAIT after no writes returns immediately
Returns count of all connected replicas
No unnecessary blocking
Replicas already synchronized
Key Challenges:
Checking if replicas are caught up
Determining "no new data" condition
Instant return path
Stage 60: WAIT with Multiple Commands
What You'll Build:
Implement full WAIT functionality that blocks until sufficient replicas acknowledge.
Technical Concepts:
Blocking until acknowledgment threshold met
Timeout handling
Aggregating ACKs from multiple replicas
What Success Looks Like:
WAIT 2 1000 blocks until 2 replicas ACK or 1000ms timeout
Sends GETACK to all replicas
Collects ACK responses
Returns count of replicas that acknowledged
Times out correctly if insufficient ACKs
Key Challenges:
Broadcasting GETACK to all replicas
Collecting and counting ACK responses
Timeout implementation
Threshold checking
Partial acknowledgment handling
Testing synchronization guarantees
RDB Persistence Section
Stage 61: RDB File Config
What You'll Build:
Add command-line options for specifying RDB file location.
Technical Concepts:
Configuration parameters
File path handling
Database persistence setup
What Success Looks Like:
Server accepts --dir and --dbfilename flags
Stores RDB file at configured location
Default values if not specified
Path validation and error handling
Key Challenges:
Parsing file path arguments
Directory existence validation
Permission checking
Default configuration
Stage 62: Read a Key
What You'll Build:
Parse an RDB file and load a single key-value pair into memory.
Technical Concepts:
RDB binary format parsing
File I/O operations
Deserialization
What Success Looks Like:
Server reads RDB file on startup
Parses one string key-value pair
Loads it into memory
Key is accessible via GET
Handles basic RDB format
Key Challenges:
Understanding RDB binary format
Byte-level parsing
Length encoding in RDB
Type indicators
Checksum validation (optional initially)
Stage 63: Read a String Value
What You'll Build:
Correctly parse and load string values from RDB files.
Technical Concepts:
String encoding in RDB
Variable-length encoding
Integer-encoded strings
What Success Looks Like:
Loads plain string values
Handles different string encodings
Correctly reconstructs original values
Supports various string lengths
Key Challenges:
Multiple string encoding formats in RDB
Length prefixes
Integer-encoded strings
LZF compressed strings (advanced)
Stage 64: Read Multiple Keys
What You'll Build:
Extend RDB parser to load all key-value pairs from file.
Technical Concepts:
Iterative parsing
Database section parsing
EOF detection
What Success Looks Like:
All keys from RDB file loaded
Multiple databases supported (DB 0, DB 1, etc.)
Complete database restoration
No keys missed
Key Challenges:
Iterating through RDB entries
Database selector parsing
End-of-file detection
Handling database numbers
Stage 65: Read Multiple String Values
What You'll Build:
Handle RDB files with many string keys and values.
Technical Concepts:
Batch loading
Memory management for large datasets
Performance optimization
What Success Looks Like:
Large RDB files load successfully
All string pairs accessible
Reasonable load time
Memory efficient loading
Key Challenges:
Handling large files
Memory efficiency
Load performance
Testing with realistic datasets
Stage 66: Read Value with Expiry
What You'll Build:
Parse expiry metadata from RDB and restore TTLs for keys.
Technical Concepts:
Expiry encoding in RDB
Timestamp parsing (seconds or milliseconds)
Expiry restoration
What Success Looks Like:
Keys with expiry load correctly
TTLs preserved from RDB
Expired keys not loaded (or immediately expire)
Expiry timers work after loading
Key Challenges:
Parsing expiry timestamps
Converting RDB timestamps to TTL
Handling already-expired keys
Both second and millisecond precision
Time zone handling
Pub/Sub Section
Stage 67: Subscribe to a Channel
What You'll Build:
Implement SUBSCRIBE command to subscribe to a single pub/sub channel.
Technical Concepts:
Publish/subscribe pattern
Channel subscription tracking
Client mode switching
What Success Looks Like:
SUBSCRIBE channelname adds subscription
Returns subscription confirmation
Client enters pub/sub mode
Only pub/sub commands work now
Key Challenges:
Tracking subscriptions per client
Client state management (pub/sub mode)
Subscription confirmation format
Channel name storage
Stage 68: Subscribe to Multiple Channels
What You'll Build:
Allow subscribing to multiple channels in one command and across commands.
Technical Concepts:
Multi-channel subscription
Subscription aggregation
Set-based tracking
What Success Looks Like:
SUBSCRIBE chan1 chan2 chan3 works
Multiple SUBSCRIBE commands accumulate
Each channel confirmed separately
Client subscribed to all channels
Key Challenges:
Storing multiple subscriptions per client
Avoiding duplicate subscriptions
Efficient channel lookup
Subscription count tracking
Stage 69: Enter Subscribed Mode
What You'll Build:
Implement pub/sub mode restrictions where only pub/sub commands are allowed.
Technical Concepts:
Modal client behavior
Command filtering
State enforcement
What Success Looks Like:
After SUBSCRIBE, only pub/sub commands accepted
Regular commands return error
PING, SUBSCRIBE, UNSUBSCRIBE, QUIT allowed
Other commands rejected with appropriate error
Key Challenges:
Command allowlist in pub/sub mode
Mode tracking per client
Clear error messages
Mode exit conditions
Stage 70: PING in Subscribed Mode
What You'll Build:
Allow PING command to work in pub/sub mode with special response format.
Technical Concepts:
Command allowlisting
Modified response formats
Mode-specific behavior
What Success Looks Like:
PING works in pub/sub mode
Returns PONG in pub/sub format (array)
Keeps client in pub/sub mode
Useful for connection keepalive
Key Challenges:
Different PING response format for pub/sub
Maintaining consistency
Testing mode-specific command behavior
Stage 71: Publish a Message
What You'll Build:
Implement PUBLISH command to send messages to a channel.
Technical Concepts:
Message broadcasting
Channel-based routing
Subscriber notification
What Success Looks Like:
PUBLISH channel message sends to all subscribers
Returns count of subscribers that received message
Non-subscribed clients can publish
Works with zero subscribers
Key Challenges:
Finding all subscribers for a channel
Message delivery to multiple clients
Counting recipients
Handling publish to non-existent channel
Stage 72: Deliver Messages
What You'll Build:
Deliver published messages to all subscribed clients in correct format.
Technical Concepts:
Message fan-out
Async message delivery
Pub/sub message format
What Success Looks Like:
Subscribers receive messages on their channels
Message format: array with "message", channel, content
All subscribers receive message
Messages delivered promptly
No cross-channel contamination
Key Challenges:
Maintaining channel→subscriber mapping
Efficient broadcast to many clients
Message format encoding
Testing message delivery
Handling slow or disconnected subscribers
Stage 73: Unsubscribe
What You'll Build:
Implement UNSUBSCRIBE to remove channel subscriptions.
Technical Concepts:
Subscription removal
Resource cleanup
Mode exit conditions
What Success Looks Like:
UNSUBSCRIBE channel removes that subscription
UNSUBSCRIBE with no args removes all
Returns unsubscription confirmation
Exit pub/sub mode when no subscriptions remain
Can re-subscribe after unsubscribing
Key Challenges:
Removing subscriptions from tracking
Handling unsubscribe from non-subscribed channel
Detecting zero-subscription state
Mode transition back to normal
Confirmation message format
Sorted Sets Section
Stage 74: Create a Sorted Set
What You'll Build:
Implement ZADD to create a sorted set with scored members.
Technical Concepts:
Sorted set data structure (often skip list + hash table)
Score-member pairs
Automatic sorting by score
What Success Looks Like:
ZADD myzset 1.0 member1 creates sorted set
TYPE returns "zset"
Member associated with score
Single member added successfully
Key Challenges:
Choosing efficient data structure (dual structure recommended)
Score storage and association
Initial insertion
Type system extension
Stage 75: Add Members
What You'll Build:
Add multiple members with scores to sorted set, handling updates and duplicates.
Technical Concepts:
Score updates for existing members
Multiple insertions
Maintaining sort order
What Success Looks Like:
ZADD myzset 2.0 member2 3.0 member3 adds multiple
Returns count of newly added members
Updating existing member's score works
Sorted order maintained automatically
Key Challenges:
Efficient insertion maintaining sort order
Detecting new vs updated members
Batch insertion performance
Score comparison (floating point)
Duplicate detection
Stage 76: Retrieve Member Rank
What You'll Build:
Implement ZRANK to get the position (rank) of a member in sorted order.
Technical Concepts:
Rank calculation
Zero-based indexing
Efficient rank lookup
What Success Looks Like:
ZRANK myzset member returns its rank (0 for lowest score)
Returns null for non-existent member
Fast operation even for large sets
Rank updates when scores change
Key Challenges:
Efficient rank calculation
Maintaining rank information
Handling ties (same score)
Performance optimization
Stage 77: List Sorted Set Members
What You'll Build:
Implement ZRANGE to retrieve members in score order by rank range.
Technical Concepts:
Range queries on sorted data
Optional score inclusion
Rank-based slicing
What Success Looks Like:
ZRANGE myzset 0 -1 returns all members in order
ZRANGE myzset 0 2 returns first three
WITHSCORES option includes scores
Results properly ordered
Key Challenges:
Efficient range extraction
Converting ranks to members
Score inclusion formatting
Large range performance
Stage 78: ZRANGE with Negative Indexes
What You'll Build:
Support negative indexes in ZRANGE counting from the end.
Technical Concepts:
Negative index conversion
Reverse counting
Index normalization
What Success Looks Like:
ZRANGE myzset -3 -1 returns last three members
ZRANGE myzset 0 -1 returns entire set
Negative indexes work with WITHSCORES
Consistent with list negative indexing
Key Challenges:
Index conversion logic
Set size tracking
Edge cases with small sets
Combining positive and negative indexes
Stage 79: Count Sorted Set Members
What You'll Build:
Implement ZCARD to return the number of members in a sorted set.
Technical Concepts:
Cardinality tracking
Fast count retrieval
What Success Looks Like:
ZCARD myzset returns member count
Constant time operation
Returns 0 for non-existent key
Updates as members added/removed
Key Challenges:
Maintaining accurate count
Updating on all modifications
Type checking
Stage 80: Retrieve Member Score
What You'll Build:
Implement ZSCORE to get the score of a specific member.
Technical Concepts:
Hash-based member lookup
Score retrieval
Null handling
What Success Looks Like:
ZSCORE myzset member returns its score
Returns null for non-existent member
Fast lookup regardless of set size
Returns score as string
Key Challenges:
Efficient member→score mapping
Floating point formatting
Null response formatting
Hash table implementation
Stage 81: Remove a Member
What You'll Build:
Implement ZREM to remove members from a sorted set.
Technical Concepts:
Member deletion
Dual structure maintenance
Count updates
What Success Looks Like:
ZREM myzset member removes it
Returns count of removed members
Supports removing multiple members
Non-existent members ignored
Set deleted when empty
Key Challenges:
Removing from both structures (skip list + hash)
Maintaining sort order
Count tracking
Memory cleanup
Batch removal
Geospatial Commands Section
Stage 82: Respond to GEOADD
What You'll Build:
Implement basic GEOADD command structure to accept longitude, latitude, and member name.
Technical Concepts:
Geospatial data storage (typically backed by sorted set)
Coordinate parsing
Geohash encoding (preparation)
What Success Looks Like:
GEOADD locations 13.361389 38.115556 "Palermo" accepted
Returns count of added locations
Basic structure in place
Command parsing works
Key Challenges:
Parsing longitude, latitude, name tuples
Coordinate value extraction
Multiple location additions
Underlying storage design
Stage 83: Validate Coordinates
What You'll Build:
Add validation for longitude and latitude ranges.
Technical Concepts:
Geographic coordinate constraints
Input validation
Error reporting
What Success Looks Like:
Longitude must be -180 to 180
Latitude must be -85.05112878 to 85.05112878
Invalid coordinates return error
Clear error messages
No data stored for invalid input
Key Challenges:
Precise range validation
Floating point comparison
Appropriate error messages
Batch validation (one invalid fails all)
Stage 84: Store a Location
What You'll Build:
Actually store location data in the underlying sorted set structure.
Technical Concepts:
Geohash calculation
Sorted set as storage backend
Encoding coordinates as score
What Success Looks Like:
Locations stored and retrievable
Underlying sorted set created
TYPE shows "zset" (geospatial uses sorted set)
Locations persist correctly
Key Challenges:
Understanding geo-backed sorted sets
Location to member mapping
Ensuring ZSET commands work on geo data
Storage format
Stage 85: Calculate Location Score
What You'll Build:
Implement geohash algorithm to convert coordinates to sorted set scores.
Technical Concepts:
Geohash encoding algorithm
Interleaving latitude and longitude bits
52-bit geohash for Redis
Score as sortable representation
What Success Looks Like:
Coordinates converted to geohash scores
Scores enable spatial proximity queries
Precise geohash calculation
Reversible encoding
Key Challenges:
Implementing geohash algorithm correctly
Bit manipulation
Precision handling
Testing geohash accuracy
Understanding interleaving
Stage 86: Respond to GEOPOS
What You'll Build:
Implement GEOPOS to retrieve stored coordinates for location names.
Technical Concepts:
Reverse lookup (member to coordinates)
Geohash decoding
Coordinate reconstruction
What Success Looks Like:
GEOPOS locations Palermo returns its coordinates
Multiple locations supported
Returns null for non-existent members
Coordinates match stored values (within precision)
Key Challenges:
Retrieving geohash score from sorted set
Preparing for decoding (next stage)
Null handling
Response formatting
Stage 87: Decode Coordinates
What You'll Build:
Implement geohash decoding to convert scores back to longitude/latitude.
Technical Concepts:
Reverse geohash algorithm
Bit de-interleaving
Precision loss handling
What Success Looks Like:
Geohash scores decoded to coordinates
Longitude and latitude extracted
Reasonable precision (within meters)
Matches original input closely
Key Challenges:
Implementing decode algorithm
Bit manipulation in reverse
Precision considerations
Floating point representation
Testing accuracy
Stage 88: Calculate Distance
What You'll Build:
Implement GEODIST to calculate distance between two stored locations.
Technical Concepts:
Haversine formula
Great circle distance
Unit conversion (m, km, mi, ft)
What Success Looks Like:
GEODIST locations Palermo Catania returns distance
Default unit is meters
Supports m, km, mi, ft units
Accurate distance calculation
Returns null if member doesn't exist
Key Challenges:
Implementing Haversine formula
Trigonometric calculations
Unit conversions
Earth radius constant
Precision and accuracy
Edge cases (antipodal points)
Stage 89: Search Within Radius
What You'll Build:
Implement GEORADIUS to find all members within a specified radius of coordinates.
Technical Concepts:
Spatial range query
Radius search
Bounding box optimization
Distance filtering
What Success Looks Like:
GEORADIUS locations 15 37 100 km returns nearby locations
All members within radius included
Distance calculations correct
Optional distance and coordinate inclusion
Sorted by distance option
Key Challenges:
Efficient spatial search (using geohash ranges)
Filtering by actual distance
Optional result enrichment
Sorting options
Performance with many locations
Unit handling in radius
Authentication Section
Stage 90: Respond to ACL WHOAMI
What You'll Build:
Implement ACL WHOAMI to return the current authenticated username.
Technical Concepts:
User context tracking
Authentication state
Default user concept
What Success Looks Like:
ACL WHOAMI returns "default" initially
Returns authenticated username after AUTH
Per-client tracking
Simple string response
Key Challenges:
Per-client user tracking
Default user handling
State management
Stage 91: Respond to ACL GETUSER
What You'll Build:
Implement ACL GETUSER to retrieve user account information.
Technical Concepts:
User account storage
ACL metadata
Structured responses
What Success Looks Like:
ACL GETUSER default returns user info
Shows flags, passwords, commands
Array response format
Returns null for non-existent users
Key Challenges:
User data structure
Response formatting
Null handling
Default user initialization
Stage 92: The nopass Flag
What You'll Build:
Support the nopass flag indicating a user requires no password.
Technical Concepts:
Password requirement flags
Authentication exemptions
Default user configuration
What Success Looks Like:
Default user has nopass flag initially
ACL GETUSER shows flags including "on" and "nopass"
Allows access without authentication
Flag array formatting
Key Challenges:
Flag representation
Multiple flag handling
Default configuration
Response formatting
Stage 93: The passwords Property
What You'll Build:
Display password hashes in ACL GETUSER response.
Technical Concepts:
Password hashing (SHA256)
Hash representation
Secure storage
What Success Looks Like:
ACL GETUSER shows passwords array
Empty array for nopass users
Password hashes shown for users with passwords
Hash format with prefix
Key Challenges:
Password storage
Hash computation
Array formatting
Security considerations
Stage 94: Setting Default User Password
What You'll Build:
Implement CONFIG SET to set a password for the default user via requirepass.
Technical Concepts:
Runtime configuration
Password setting
ACL modification
What Success Looks Like:
CONFIG SET requirepass mypassword sets default user password
Removes nopass flag
Adds password hash to default user
ACL GETUSER reflects changes
Returns OK
Key Challenges:
CONFIG command implementation
requirepass parameter handling
ACL updates
Password hashing
Removing nopass flag
Authentication requirement activation
Stage 95: The AUTH Command
What You'll Build:
Implement AUTH command for client authentication.
Technical Concepts:
Password-based authentication
Client authentication state
Password verification
What Success Looks Like:
AUTH password authenticates client
AUTH username password for non-default users
Returns OK on success
Returns error on failure
Client state updated on success
Key Challenges:
Password hash comparison
Two-argument vs one-argument forms
Client state updates
Error messages
Secure comparison
Stage 96: Enforce Authentication
What You'll Build:
Require authentication before allowing commands when authentication is configured.
Technical Concepts:
Authentication enforcement
Command blocking
Allowlist of pre-auth commands
What Success Looks Like:
Commands rejected with NOAUTH error before AUTH
AUTH and HELLO commands work without authentication
All commands work after successful AUTH
Per-client authentication tracking
Key Challenges:
Pre-authentication command allowlist
Error message consistency
State checking on every command
HELLO command support
Testing authentication flow
Stage 97: Authenticate Using AUTH
What You'll Build:
Complete end-to-end authentication flow with proper state management.
Technical Concepts:
Complete authentication cycle
Session management
Security best practices
What Success Looks Like:
Full workflow: CONFIG SET requirepass → commands rejected → AUTH → commands work
Multiple clients authenticate independently
Re-authentication works
Failed auth doesn't affect state
Secure and reliable
Key Challenges:
Integration of all auth components
Testing complete flows
Multiple client scenarios
Security testing
Error recovery
State consistency
Conclusion
This challenge takes you from basic networking concepts to advanced distributed systems features. Each stage builds essential skills:
Foundation: TCP servers, protocols, concurrency
Data structures: Strings, lists, streams, sorted sets
Advanced features: Transactions, replication, pub/sub
Persistence: RDB file format
Specialized: Geospatial indexing, authentication
Success at each stage means working code that passes tests, but also understanding the underlying concepts. The real learning comes from debugging issues, optimizing performance, and seeing how Redis's design choices solve real problems.
Take your time with each stage, understand why things work the way they do, and don't hesitate to research Redis's actual implementation for insights. By the end, you'll have built a substantial piece of infrastructure software and gained deep knowledge of systems programming.
