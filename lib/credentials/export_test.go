package credentials

import "sync"

// ResetEnvWarningForTest lets a test observe the warning that is printed
// once per process.
func ResetEnvWarningForTest() { warnAPIKeyWins = sync.Once{} }
