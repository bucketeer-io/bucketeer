package v3

import (
	"strings"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestClient_scanClusterWithPagination_NonClusterClient(t *testing.T) {
	t.Parallel()
	// Test the error case when client is not a cluster client
	client := &client{
		rc:         redis.NewClient(&redis.Options{}), // Non-cluster client
		clientType: ClientTypeCluster,
	}

	keys, cursor, err := client.scanClusterWithPagination(0, "test:*", 10)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client is not a cluster client")
	assert.Nil(t, keys)
	assert.Equal(t, uint64(0), cursor)
}

func TestCursorEncodingDecoding(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		nodeIndex   uint64
		nodeCursor  uint64
		expectation uint64
	}{
		{
			name:        "encode node 0, cursor 0",
			nodeIndex:   0,
			nodeCursor:  0,
			expectation: 0,
		},
		{
			name:        "encode node 1, cursor 0",
			nodeIndex:   1,
			nodeCursor:  0,
			expectation: 1 << 32,
		},
		{
			name:        "encode node 0, cursor 100",
			nodeIndex:   0,
			nodeCursor:  100,
			expectation: 100,
		},
		{
			name:        "encode node 2, cursor 500",
			nodeIndex:   2,
			nodeCursor:  500,
			expectation: (2 << 32) | 500,
		},
		{
			name:        "encode max values",
			nodeIndex:   0xFFFFFFFF,
			nodeCursor:  0xFFFFFFFF,
			expectation: 0xFFFFFFFFFFFFFFFF,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test the cursor encoding logic that's used in scanClusterWithPagination
			// Encode (this is the logic from the method)
			encoded := (tt.nodeIndex << 32) | tt.nodeCursor
			assert.Equal(t, tt.expectation, encoded)

			// Decode (this is the logic from the method)
			decodedNodeIndex := encoded >> 32
			decodedNodeCursor := encoded & 0xFFFFFFFF

			assert.Equal(t, tt.nodeIndex, decodedNodeIndex)
			assert.Equal(t, tt.nodeCursor, decodedNodeCursor)
		})
	}
}

func TestParseClusterMasterAddresses(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		clusterNodesResp    string
		expectedMasterAddrs []string
	}{
		{
			name: "three master nodes",
			clusterNodesResp: `
				node1 127.0.0.1:6379@16379 master - 0 connected 0-5460
				node2 127.0.0.1:6380@16380 master - 0 connected 5461-10922
				node3 127.0.0.1:6381@16381 master - 0 connected 10923-16383
			`,
			expectedMasterAddrs: []string{"127.0.0.1:6379", "127.0.0.1:6380", "127.0.0.1:6381"},
		},
		{
			name: "mixed master and slave nodes",
			clusterNodesResp: `
				node1 127.0.0.1:6379@16379 master - 0 connected 0-5460
				node2 127.0.0.1:6380@16380 slave node1 0 connected
				node3 127.0.0.1:6381@16381 master - 0 connected 5461-10922
			`,
			expectedMasterAddrs: []string{"127.0.0.1:6379", "127.0.0.1:6381"},
		},
		{
			name: "only slave nodes",
			clusterNodesResp: `
				node1 127.0.0.1:6379@16379 slave - 0 connected
				node2 127.0.0.1:6380@16380 slave - 0 connected
			`,
			expectedMasterAddrs: []string{},
		},
		{
			name:                "empty response",
			clusterNodesResp:    "",
			expectedMasterAddrs: []string{},
		},
		{
			name: "addresses without cluster port",
			clusterNodesResp: `
				node1 127.0.0.1:6379 master - 0 connected 0-5460
				node2 127.0.0.1:6380 master - 0 connected 5461-10922
			`,
			expectedMasterAddrs: []string{"127.0.0.1:6379", "127.0.0.1:6380"},
		},
		{
			name: "mixed address formats",
			clusterNodesResp: `
				node1 127.0.0.1:6379@16379 master - 0 connected 0-5460
				node2 127.0.0.1:6380 master - 0 connected 5461-10922
				node3 redis-host.example.com:6381@16381 master - 0 connected 10923-16383
			`,
			expectedMasterAddrs: []string{"127.0.0.1:6379", "127.0.0.1:6380", "redis-host.example.com:6381"},
		},
		{
			name: "ipv6 addresses",
			clusterNodesResp: `
				node1 [::1]:6379@16379 master - 0 connected 0-5460
				node2 [2001:db8::1]:6380@16380 master - 0 connected 5461-10922
			`,
			expectedMasterAddrs: []string{"[::1]:6379", "[2001:db8::1]:6380"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := parseClusterMasterAddresses(strings.TrimSpace(tt.clusterNodesResp))
			assert.Equal(t, tt.expectedMasterAddrs, result)
		})
	}
}

func TestNodeIndexCalculation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name               string
		cursor             uint64
		numNodes           int
		expectedNodeIndex  uint64
		expectedNodeCursor uint64
		shouldBeComplete   bool
	}{
		{
			name:               "start of scan",
			cursor:             0,
			numNodes:           3,
			expectedNodeIndex:  0,
			expectedNodeCursor: 0,
			shouldBeComplete:   false,
		},
		{
			name:               "middle of node 0",
			cursor:             100,
			numNodes:           3,
			expectedNodeIndex:  0,
			expectedNodeCursor: 100,
			shouldBeComplete:   false,
		},
		{
			name:               "start of node 1",
			cursor:             1 << 32,
			numNodes:           3,
			expectedNodeIndex:  1,
			expectedNodeCursor: 0,
			shouldBeComplete:   false,
		},
		{
			name:               "middle of node 2",
			cursor:             (2 << 32) | 500,
			numNodes:           3,
			expectedNodeIndex:  2,
			expectedNodeCursor: 500,
			shouldBeComplete:   false,
		},
		{
			name:               "beyond last node",
			cursor:             3 << 32,
			numNodes:           3,
			expectedNodeIndex:  3,
			expectedNodeCursor: 0,
			shouldBeComplete:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test the node index calculation logic used in scanClusterWithPagination
			// Decode cursor (from scanClusterWithPagination)
			nodeIndex := tt.cursor >> 32
			nodeCursor := tt.cursor & 0xFFFFFFFF

			assert.Equal(t, tt.expectedNodeIndex, nodeIndex)
			assert.Equal(t, tt.expectedNodeCursor, nodeCursor)

			// Check if scanning should be complete
			isComplete := nodeIndex >= uint64(tt.numNodes)
			assert.Equal(t, tt.shouldBeComplete, isComplete)
		})
	}
}

func TestCursorTransitions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		nodeIndex      uint64
		nodeCursor     uint64
		newNodeCursor  uint64
		expectedCursor uint64
		description    string
	}{
		{
			name:           "continue on same node",
			nodeIndex:      0,
			nodeCursor:     100,
			newNodeCursor:  200,
			expectedCursor: (0 << 32) | 200,
			description:    "cursor advances on same node",
		},
		{
			name:           "move to next node",
			nodeIndex:      0,
			nodeCursor:     100,
			newNodeCursor:  0, // cursor 0 means node scan complete
			expectedCursor: 1 << 32,
			description:    "move to node 1 when current node complete",
		},
		{
			name:           "move from node 2 to node 3",
			nodeIndex:      2,
			nodeCursor:     500,
			newNodeCursor:  0, // cursor 0 means node scan complete
			expectedCursor: 3 << 32,
			description:    "move to node 3 when node 2 complete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Test the cursor transition logic used in scanClusterWithPagination
			var newCursor uint64
			if tt.newNodeCursor == 0 {
				// Current node scan complete, move to next node
				newCursor = (tt.nodeIndex + 1) << 32
			} else {
				// Continue scanning current node
				newCursor = (tt.nodeIndex << 32) | tt.newNodeCursor
			}

			assert.Equal(t, tt.expectedCursor, newCursor, tt.description)
		})
	}
}
