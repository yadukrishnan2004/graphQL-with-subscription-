package graph

import (
	"gql/graph/model"
	"sync"

	"gorm.io/gorm"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct{
	DB *gorm.DB
	Subscriber map[string]chan*model.Todo
	Mu sync.Mutex
}


