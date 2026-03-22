package script

import (
	"errors"

	"github.com/dop251/goja"
)

type Engine struct {
	runtime *goja.Runtime
}

func NewEngine() *Engine {
	return &Engine{
		runtime: goja.New(),
	}
}

func (e *Engine) Execute(script string, config map[string]interface{}) (map[string]interface{}, error) {
	e.runtime.Set("config", config)
	value, err := e.runtime.RunString(script)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, errors.New("script returned nil")
	}
	result, ok := value.Export().(map[string]interface{})
	if !ok {
		return nil, errors.New("script must return an object")
	}
	return result, nil
}

func (e *Engine) ExecuteWithCallback(script string, config map[string]interface{}, callback func(map[string]interface{}) error) error {
	result, err := e.Execute(script, config)
	if err != nil {
		return err
	}
	return callback(result)
}
