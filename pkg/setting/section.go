package setting

import "time"

type ServerSettings struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var sections = make(map[string]interface{})

// ReadSection 读取单个作用域
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	if _, ok := sections[k]; ok {
		sections[k] = v
	}
	return nil
}

// ReloadAllSection 加载所有作用域
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
