package services

type serverServices struct {
	extHandle func(u *Unit, ext ...string) error
	*services
	*handleHook
	*serverHook
}

func newServerServices(v func(u *Unit, ext ...string) error) *serverServices {
	return &serverServices{
		handleHook: newHandleHook(),
		services:   newService(),
		serverHook: new(serverHook),
		extHandle:  v,
	}
}

//Register 注册服务
func (s *serverServices) Register(name string, h interface{}, ext ...string) {
	groups, err := reflectHandle(name, h)
	if err != nil {
		panic(err)
	}
	if err := s.addGroup(groups, ext...); err != nil {
		panic(err)
	}
}
func (s *serverServices) handleExt(g *Unit, ext ...string) error {
	if s.extHandle == nil {
		return nil
	}
	return s.extHandle(g, ext...)
}

//addGroup 添加服务注册
func (s *serverServices) addGroup(g *UnitGroup, ext ...string) error {
	for _, u := range g.Services {

		//添加预处理函数
		if err := s.handleHook.AddHandling(u.Service, u.GetHandlings()...); err != nil {
			return err
		}

		//执行处理扩展函数
		if err := s.handleExt(u, ext...); err != nil {
			return err
		}

		//添加服务
		if err := s.services.AddHanler(u.Service, u.Handle); err != nil {
			return err

		}

		//添加后处理函数
		if err := s.handleHook.AddHandled(u.Service, u.GetHandleds()...); err != nil {
			return err
		}

		//添加降级函数
		if err := s.services.AddFallback(u.Service, u.Fallback); err != nil {
			return err

		}
	}

	//添加关闭函数
	if err := s.handleHook.AddClosingHanle(g.Closing); err != nil {
		return err
	}

	return nil
}