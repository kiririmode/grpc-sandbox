package common

import (
	"errors"
	"testing"
)

// 正常に初期化・終了処理を行う Resouce
type successResouce struct{}

func (r *successResouce) Initialize() error {
	return nil
}
func (r *successResouce) Finalize() error {
	return nil
}
func (r *successResouce) Name() string {
	return "success"
}

// 初期化処理、終了処理に失敗する Resource
type failureResource struct{}

func (r *failureResource) Initialize() error {
	return errors.New("failed")
}
func (r *failureResource) Finalize() error {
	return errors.New("failed")
}
func (r *failureResource) Name() string {
	return "failure"
}

func TestResourceManager_Initialize(t *testing.T) {
	t.Run("すべての初期化処理が正常に完了すれば、正常終了する", func(t *testing.T) {
		sut := &ResourceManager{}
		sut.AddResource(&successResouce{}).AddResource(&successResouce{})

		err := sut.Initialize()
		if err != nil {
			t.Errorf("error should be nil, but got %s", err)
		}
	})
	t.Run("初期化処理が異常終了したものが１つでもあれば異常終了する", func(t *testing.T) {
		sut := &ResourceManager{}
		sut.AddResource(&failureResource{}).
			AddResource(&successResouce{})

		err := sut.Initialize()
		if err == nil {
			t.Error("error should be occured, but got success")
		}
	})
}

func TestResourceManager_Finalize(t *testing.T) {
	t.Run("すべての終了処理が正常に完了すれば、正常終了する", func(t *testing.T) {
		sut := &ResourceManager{}
		sut.AddResource(&successResouce{}).AddResource(&successResouce{})

		errArray := sut.Finalize()
		if len(errArray) > 0 {
			for _, err := range errArray {
				t.Errorf("error should be nil, but got %s", err)
			}
		}
	})
	t.Run("終了処理が異常終了したものが１つでもあれば異常終了する", func(t *testing.T) {
		sut := &ResourceManager{}
		sut.AddResource(&failureResource{}).
			AddResource(&successResouce{})

		errs := sut.Finalize()
		if len(errs) == 0 {
			t.Error("error should be occured, but got success")
		}
	})
	t.Run("終了処理が異常終了するリソースがあっても続行する", func(t *testing.T) {
		sut := &ResourceManager{}
		sut.AddResource(&failureResource{}).
			AddResource(&failureResource{})

		errs := sut.Finalize()
		if len(errs) != 2 {
			t.Errorf("error count must be 2, but got %d", len(errs))
		}
	})
}
