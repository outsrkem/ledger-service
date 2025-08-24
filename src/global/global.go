// Defined in a global package (e.g. global.go)
package global

import "ledger/src/pkg/cache"

// GCache is the global instance ID cache
// (initialized when the project starts)
var GCache = cache.NewInstanceCache()
