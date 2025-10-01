package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/db"
)

var (
	_db *db.Client

	// TODO: 定义 Model 变量
	// Subscription *db.Model[ModelSubscription]
	// Node *db.Model[ModelNode]
	// Route *db.Model[ModelRoute]
)

// ConnectDatabase 连接数据库并初始化表
func ConnectDatabase() (err error) {
	log.Info("try connect database")

	// 创建数据库客户端（暂时不传入 Model，等定义后再添加）
	_db, err = db.New(State.Config.Db)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// TODO: 初始化 Model
	// Subscription = db.NewModel[ModelSubscription](Db()).
	// 	SetNotFound(xerror.NewError(3000)).
	// 	SetDuplicatedKeyError(xerror.NewError(3001))

	log.Info("database connected successfully")

	return nil
}

// Db 返回数据库客户端
func Db() *db.Client {
	return _db
}

// NewScoop 创建新的数据库操作作用域
func NewScoop() *db.Scoop {
	return _db.NewScoop()
}

// Begin 开始事务
func Begin() *db.Scoop {
	return NewScoop().Begin()
}

// CommitOrRollback 执行事务逻辑，自动提交或回滚
func CommitOrRollback(logic func(tx *db.Scoop) error) error {
	return NewScoop().CommitOrRollback(NewScoop().Begin(), logic)
}
