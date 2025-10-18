package storage

import ( 
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Environment variable names used to bootstrap the MongoDB client.
const (
	mongoURIEnv        = "EXTENSION_OBSERVER_MONGODB_URI"
	mongoDatabaseEnv   = "EXTENSION_OBSERVER_MONGODB_DATABASE"
	mongoCollectionEnv = "EXTENSION_OBSERVER_MONGODB_COLLECTION"
)

// ExtensionEventDocument represents a stored observer event entry.
type ExtensionEventDocument struct {
	Event     string    `bson:"event"` 
	ClientID  string    `bson:"clientId"` 
	Time      string    `bson:"time"`
	CreatedAt time.Time `bson:"createdAt"`
}

var (
	syncOnce        sync.Once 
	mongoInitErr     error 
	mongoCollection  *mongo.Collection // 仅作一次成功初始化 单例
	defaultTimeout   = 5 * time.Second 
	errMissingConfig = errors.New("extension observer MongoDB configuration is incomplete")
)

// SaveExtensionEvent persists an extension event document using a lazily 
// initialised MongoDB collection.
// 把一个 ExtensionEventDocument 文档保存到 MongoDB 里。
func SaveExtensionEvent(ctx context.Context, doc ExtensionEventDocument) error { 
	// 先确保目标集合存在（ensureCollection），
	if err := ensureCollection(ctx); err != nil {
		return err
	}
	// 然后给文档加上当前 UTC 时间戳
	doc.CreatedAt = time.Now().UTC()
	// 最后用 InsertOne 把它插入到 MongoDB
	_, err := mongoCollection.InsertOne(ctx, doc)
	return err
}

// 初始化 MongoDB 连接，并拿到目标集合对象（mongoCollection）。
func ensureCollection(ctx context.Context) error {
	// fast-path：已就绪直接返回
	if mongoCollection != nil {
		return nil
	}

	uri := os.Getenv(mongoURIEnv) // 从环境变量读取 MongoDB 连接 URI
	database := os.Getenv(mongoDatabaseEnv) // 读取数据库名
	collection := os.Getenv(mongoCollectionEnv) // 读取集合名

	// 检查必要配置是否缺失（保持原有错误语义）
	if uri == "" || database == "" || collection == "" {
		mongoInitErr = fmt.Errorf("%w: expected %s, %s and %s", errMissingConfig, mongoURIEnv, mongoDatabaseEnv, mongoCollectionEnv) //🟢1（保留）
		return mongoInitErr
	}

	clientOpts := options.Client().ApplyURI(uri) // 根据 URI 创建客户端配置 

	timeoutCtx, cancel := context.WithTimeout(ctx, defaultTimeout) // 设置连接超时 
	defer cancel()

	client, err := mongo.Connect(timeoutCtx, clientOpts) // 连接 MongoDB
	if err != nil {
		mongoInitErr = fmt.Errorf("extension observer: connect mongo: %w", err)
		return mongoInitErr
	}

	// 尝试 ping 一下，确认连接可用
	if err := client.Ping(timeoutCtx, nil); err != nil {
		mongoInitErr = fmt.Errorf("extension observer: ping mongo: %w", err)
		return mongoInitErr
	}

	// 仅将“成功后的赋值”放到 once 里，失败不受 once 影响→可重试
	syncOnce.Do(func() {
		// 双检：避免并发下重复赋值
		if mongoCollection == nil {
			mongoCollection = client.Database(database).Collection(collection)
		}
	})

	// 若其他并发已成功设置，也算成功
	if mongoCollection != nil {
		mongoInitErr = nil
		return nil
	}

	// 理论上不会到这里；兜底保证可感知错误
	mongoInitErr = fmt.Errorf("extension observer: failed to initialize mongo collection")
	return mongoInitErr
}

