# Godular

## Description

A modular dependency injection framework for Go applications.

## Features

- **Modular Architecture**: Organize your application into reusable modules
- **Dependency Injection**: Automatic dependency resolution between modules
- **Provider System**: Register and retrieve providers across modules
- **Event System**: Event-driven communication between modules
- **Lifecycle Hooks**: Initialize providers with `OnInit` hooks
- **Graph Visualization**: Generate dependency graphs
- **Generic Model Layer**: GORM-based generic repository pattern
- **Mocking Support**: Built-in mocking utilities for testing

## Installation

```bash
go get github.com/LordPax/godular
```

## Quick Start

### 1. Create a Module

```go
package user

import "github.com/LordPax/godular/core"

var module *UserModule

type UserModule struct {
    *core.Module
}

func NewUserModule() *UserModule {
    module := &UserModule{
        Module: core.NewModule("UserModule"),
    }

    module.AddModule(database.Module())
    module.AddProvider(NewUserModel(module))
    module.AddProvider(NewUserService(module))

    return module
}

func Module() *UserModule {
    if module == nil {
        module = NewUserModule()
    }
    return module
}
```

### 2. Create a Model

```go
package user

import "github.com/LordPax/godular/core"

type IUserModel interface {
    core.IModel[*User]
    QueryFindAll(q query.QueryFilter) ([]*User, error)
    CountAll() (int64, error)
}

type UserModel struct {
    *core.Model[*User]
    databaseService database.IDatabaseService
}

func NewUserModel(module core.IModule) *UserModel {
    service := &UserModel{
        Model:           core.NewModel[*User]("UserModel"),
        databaseService: module.Get("DatabaseService").(database.IDatabaseService),
    }

    // Register event handler
    module.On("db:migrate", service.Migrate)

    return service
}

func (um *UserModel) OnInit() error {
    um.SetDB(um.databaseService.GetDB())
    return nil
}

func (um *UserModel) CountAll() (int64, error) {
    var count int64
    err := um.databaseService.GetDB().Model(&User{}).
        Where("deleted_at IS NULL").
        Count(&count).Error
    return count, err
}
```

### 3. Create a Service

```go
package user

import "github.com/LordPax/godular/core"

type IUserService interface {
    core.IProvider
    FindByID(id string) (*User, error)
    Create(user *User) error
    IsUserExists(email, username string) bool
}

type UserService struct {
    *core.Provider
    userModel IUserModel
}

func NewUserService(module core.IModule) *UserService {
    return &UserService{
        Provider:  core.NewProvider("UserService"),
        userModel: module.Get("UserModel").(IUserModel),
    }
}

func (us *UserService) FindByID(id string) (*User, error) {
    return us.userModel.FindByID(id)
}

func (us *UserService) Create(user *User) error {
    return us.userModel.Create(user)
}

func (us *UserService) IsUserExists(email, username string) bool {
    emailExists, _ := us.userModel.CountBy("email", email)
    usernameExists, _ := us.userModel.CountBy("username", username)
    return emailExists > 0 || usernameExists > 0
}
```
