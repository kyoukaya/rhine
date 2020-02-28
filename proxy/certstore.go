package proxy

import (
	"crypto/tls"
	"sync"

	"github.com/golang/groupcache/lru"
	"github.com/kyoukaya/rhine/log"
)

type certStore struct {
	mutex sync.Mutex
	lru   *lru.Cache
	log.Logger
}

// cmd/example/main.go runs weighs at around 35MB with a cache of 64 certs filled,
// which is pretty reasonable at full and is highly unlikely that a user would
// encounter so many hostnames while playing Arknights.
const certCacheSize = 64

func newCertStore(logger log.Logger) *certStore {
	store := &certStore{
		lru:    lru.New(certCacheSize),
		Logger: logger,
	}
	store.lru.OnEvicted = func(key lru.Key, value interface{}) {
		logger.Verbosef("certstore: %s evicted", key.(string))
	}
	return store
}

func (store *certStore) Fetch(hostname string, gen func() (*tls.Certificate, error)) (*tls.Certificate, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	cert := store.getCert(hostname)
	// On cache hit
	if cert != nil {
		return cert, nil
	}
	// On cache miss, cert == nil
	cert, err := gen()
	if err != nil {
		store.Warnf("Cached missed on %s and failed to gen cert with %s", hostname, err)
		return cert, err
	}
	store.lru.Add(hostname, cert)
	return cert, nil
}

func (store *certStore) getCert(hostname string) *tls.Certificate {
	cert, ok := store.lru.Get(hostname)
	if !ok {
		return nil
	}
	return cert.(*tls.Certificate)
}
