package actionplugin

import (
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	dto "../../../dto"
	"../../../settings"
)

type HandlerFunc func() (*dto.WechatReplyMsg, error)
type HandlerFactoryFunc func(*dto.WXBizMsg) (HandlerFunc, error)

type ActionPlugin struct {
	fileName       string
	handlerFactory HandlerFactoryFunc
	rwLock         *sync.RWMutex
}

func New(fileName string) *ActionPlugin {
	rwLock := new(sync.RWMutex)
	ap := ActionPlugin{
		rwLock:   rwLock,
		fileName: fileName,
	}
	ap.Load()
	return &ap
}

func (ap *ActionPlugin) Load() {
	ap.rwLock.Lock()
	defer ap.rwLock.Unlock()
	pluginPath := filepath.Join(settings.BaseDir, "plugin", fmt.Sprintf("%s.so", ap.fileName))
	p, err := plugin.Open(pluginPath)
	if err != nil {
		settings.GetLogger(nil).Println(err.Error())
		ap.handlerFactory = nil
		return
	}
	handler, err := p.Lookup("GetHandler")
	if err != nil {
		settings.GetLogger(nil).Println(err.Error())
		ap.handlerFactory = nil
		return
	}
	fn, ok := handler.(HandlerFactoryFunc)
	if !ok {
		settings.GetLogger(nil).Println("some err with plugin entry")
		ap.handlerFactory = nil
		return
	}
	ap.handlerFactory = fn

}

func (ap *ActionPlugin) GetHandler(msg *dto.WXBizMsg) (HandlerFunc, error) {
	ap.rwLock.RLock()
	defer ap.rwLock.RUnlock()
	if ap.handlerFactory == nil {
		return nil, fmt.Errorf("fail to load plugin:%s", ap.fileName)
	}
	fn, err := ap.handlerFactory(msg)
	return fn, err
}
