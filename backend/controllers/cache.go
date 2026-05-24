package controllers

import (
	"backend/config"
)

// 清除所有商品列表缓存
func ClearAllProductListCache() {
	keys, _ := config.RDB.Keys(config.Ctx, "product:list:*").Result()
	if len(keys) > 0 {
		config.RDB.Del(config.Ctx, keys...)
	}
}

// 清除单个商品缓存
func ClearSingleProductCache(id string) {
	config.RDB.Del(config.Ctx, "product:item:"+id)
}

// 清除所有分类缓存
func ClearAllCategoryCache() {
	config.RDB.Del(config.Ctx, "category:parent")
	childKeys, _ := config.RDB.Keys(config.Ctx, "category:child:*").Result()
	if len(childKeys) > 0 {
		config.RDB.Del(config.Ctx, childKeys...)
	}
}
