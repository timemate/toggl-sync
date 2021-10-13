package sync

type SyncApi interface {
	Sync(tasks []Task) (err error)
}
