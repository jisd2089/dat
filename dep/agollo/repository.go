package agollo

import (
	"strconv"
	"github.com/coocood/freecache"
)

const (
	empty = ""

	//50m
	apolloConfigCacheSize = 50 * 1024 * 1024

	//1 minute
	configCacheExpireTime = 120
)

type Repository struct {
	currentConnApolloConfig *ApolloConnConfig
	apolloConfigCache       *freecache.Cache
}

var (
	currentConnApolloConfig = &ApolloConnConfig{}

	//config from apollo
	apolloConfigCache = freecache.NewCache(apolloConfigCacheSize)
)

func (r *Repository) updateApolloConfig(notifyChan chan *ChangeEvent, apolloConfig *ApolloConfig) {
	if apolloConfig == nil {
		logger.Error("apolloConfig is null,can't update!")
		return
	}
	//get change list
	changeList := r.updateApolloConfigCache(apolloConfig.Configurations, configCacheExpireTime)

	if len(changeList) > 0 {
		//create config change event base on change list
		event := createConfigChangeEvent(changeList, apolloConfig.NamespaceName)

		//push change event to channel
		pushChangeEvent(notifyChan, event)
	}

	//update apollo connection config

	r.currentConnApolloConfig.Lock()
	defer r.currentConnApolloConfig.Unlock()
	r.currentConnApolloConfig = &apolloConfig.ApolloConnConfig
}

func (r *Repository) updateApolloConfigCache(configurations map[string]string, expireTime int) map[string]*ConfigChange {
	if (configurations == nil || len(configurations) == 0) && apolloConfigCache.EntryCount() == 0 {
		return nil
	}

	//get old keys
	mp := map[string]bool{}
	it := r.apolloConfigCache.NewIterator()
	for en := it.Next(); en != nil; en = it.Next() {
		mp[string(en.Key)] = true
	}

	changes := make(map[string]*ConfigChange)

	if configurations != nil {
		// update new
		// keys
		for key, value := range configurations {
			//key state insert or update
			//insert
			if !mp[key] {
				changes[key] = createAddConfigChange(value)
			} else {
				//update
				oldValue, _ := r.apolloConfigCache.Get([]byte(key))
				if string(oldValue) != value {
					changes[key] = createModifyConfigChange(string(oldValue), value)
				}
			}

			r.apolloConfigCache.Set([]byte(key), []byte(value), expireTime)
			delete(mp, string(key))
		}
	}

	// remove del keys
	for key := range mp {
		//get old value and del
		oldValue, _ := r.apolloConfigCache.Get([]byte(key))
		changes[key] = createDeletedConfigChange(string(oldValue))

		r.apolloConfigCache.Del([]byte(key))
	}

	return changes
}

//base on changeList create Change event
func createConfigChangeEvent(changes map[string]*ConfigChange, nameSpace string) *ChangeEvent {
	return &ChangeEvent{
		Namespace: nameSpace,
		Changes:   changes,
	}
}

func (r *Repository) touchApolloConfigCache() error {
	r.updateApolloConfigCacheTime(configCacheExpireTime)
	return nil
}

func (r *Repository) updateApolloConfigCacheTime(expireTime int) {
	it := r.apolloConfigCache.NewIterator()
	for i := int64(0); i < r.apolloConfigCache.EntryCount(); i++ {
		entry := it.Next()
		if entry == nil {
			break
		}
		r.apolloConfigCache.Set([]byte(entry.Key), []byte(entry.Value), expireTime)
	}
}

func GetApolloConfigCache() *freecache.Cache {
	return apolloConfigCache
}

func GetCurrentApolloConfig() *ApolloConnConfig {
	currentConnApolloConfig.RLock()
	defer currentConnApolloConfig.RUnlock()
	return currentConnApolloConfig
}

func getConfigValue(key string) interface{} {
	value, err := apolloConfigCache.Get([]byte(key))
	if err != nil {
		logger.Error("get config value fail!err:", err)
		return empty
	}

	return string(value)
}

func getValue(key string) string {
	value := getConfigValue(key)
	if value == nil {
		return empty
	}

	return value.(string)
}

func GetStringValue(key string, defaultValue string) string {
	value := getValue(key)
	if value == empty {
		return defaultValue
	}

	return value
}

func GetIntValue(key string, defaultValue int) int {
	value := getValue(key)

	i, err := strconv.Atoi(value)
	if err != nil {
		logger.Debug("convert to int fail!error:", err)
		return defaultValue
	}

	return i
}

func GetFloatValue(key string, defaultValue float64) float64 {
	value := getValue(key)

	i, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Debug("convert to float fail!error:", err)
		return defaultValue
	}

	return i
}

func GetBoolValue(key string, defaultValue bool) bool {
	value := getValue(key)

	b, err := strconv.ParseBool(value)
	if err != nil {
		logger.Debug("convert to bool fail!error:", err)
		return defaultValue
	}

	return b
}
