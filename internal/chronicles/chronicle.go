package chronicles

import "time"

type Chronicle interface {
	RegisterRelease(id string, t time.Time, kind, name, revision, namespace string)
	UnregisterRelease(id string)
}
