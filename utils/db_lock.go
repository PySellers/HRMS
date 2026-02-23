package utils

import "sync"

// Global mutex to protect db.json
var DBMutex sync.Mutex