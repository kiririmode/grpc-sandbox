package common

import (
	"github.com/pkg/errors"
)

// Resource は初期化・終了処理を持つオブジェクトを抽象化する Interface
type Resource interface {
	// リソース名を返却する ("Configuration" 等)
	Name() string
	// 初期化処理を実行する
	Initialize() error
	// 終了処理を実行する
	Finalize() error
}

// ResourceManager はリソースの初期化・終了処理を管理するマネージャ
type ResourceManager struct {
	resources []Resource
}

// NewResourceManager は、管理対象として rs を含む新しい ResourceManager を返却する
func NewResourceManager(rs []Resource) *ResourceManager {
	return &ResourceManager{
		resources: rs,
	}
}

// AddResource はマネージャが管理する Resouce として r を追加する
func (m *ResourceManager) AddResource(r Resource) *ResourceManager {
	m.resources = append(m.resources, r)
	return m
}

// Initialize は管理しているリソースを追加された順に初期化する。
// エラーが起こった場合は、その時点で処理を打ち切る。
func (m *ResourceManager) Initialize() error {
	for _, r := range m.resources {
		err := r.Initialize()
		if err != nil {
			return errors.Wrapf(err, "failed to initialize %s", r.Name())
		}
	}
	return nil
}

// Finalize は管理しているリソースに対し、追加順と逆順に終了処理を行う。
// エラーが起こった場合も、一通りの終了処理を行う
func (m *ResourceManager) Finalize() []error {
	errArray := make([]error, 0)

	for i := len(m.resources) - 1; i >= 0; i-- {
		err := m.resources[i].Finalize()
		if err != nil {
			errArray = append(errArray, err)
		}
	}
	return errArray
}
