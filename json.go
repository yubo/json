// not thread safe,
package json

import (
	"strings"

	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

func Decode(b []byte) (*Data, error) {
	data := map[string]interface{}{}
	if err := yaml.Unmarshal(b, &data); err != nil {
		return nil, err
	}

	return &Data{data: data}, nil
}

func Encode(obj interface{}) (*Data, error) {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return Decode(b)
}

type Data struct {
	data map[string]interface{}
	path []string
}

func (p *Data) Set(path string, v interface{}) error {
	if path == "" {
		b, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(b, &p.data)
	}

	ps := strings.Split(path, ".")
	src := map[string]interface{}{ps[len(ps)-1]: v}

	for i := len(ps) - 2; i >= 0; i-- {
		src = map[string]interface{}{ps[i]: src}
	}

	p.data = mergeValues(p.data, src)

	return nil
}

// unmarshal json path
func (p *Data) Read(path string, into interface{}) error {
	if into == nil {
		return nil
	}

	if v := p.GetRaw(path); v != nil {
		data, err := yaml.Marshal(v)
		//klog.V(5).InfoS("marshal", "v", v, "data", string(data), "err", err)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, into)
		if err != nil {
			klog.V(5).InfoS("unmarshal", "data", string(data), "err", err)
			if klog.V(5).Enabled() {
				panic(err)
			}
			return err
		}
	}

	if v, ok := into.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	if klog.V(10).Enabled() {
		b, _ := yaml.Marshal(into)
		klog.Infof("Read \n[%s]\n%s", path, string(b))
	}
	return nil
}

func (p *Data) GetData(path string) *Data {
	if p == nil {
		return nil
	}

	data, _ := p.GetRaw(path).(map[string]interface{})

	// noneed deepCopy
	out := new(Data)

	out.path = append(clonePath(p.path), parsePath(path)...)
	out.data = data

	return out
}

func (p *Data) GetRaw(path string) interface{} {
	if path == "" {
		return Values(p.data)
	}

	v, err := Values(p.data).PathValue(path)
	if err != nil {
		klog.V(5).InfoS("get pathValue err, ignored", "path", path, "v", v, "err", err)
		return nil
	}
	return v
}

func (p *Data) GetString(path string) string {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return ""
	}

	return ToString(v)
}

func (p *Data) GetBool(path string) (bool, error) {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return false, err
	}

	return ToBool(v), nil
}

func (p *Data) GetBoolDef(path string, def bool) bool {
	v, err := p.GetBool(path)
	if err != nil {
		return def
	}
	return v
}

func (p *Data) GetFloat64(path string) (float64, error) {
	v, err := Values(p.data).PathValue(path)
	if err != nil {
		return 0, err
	}

	return ToFloat64(v), nil
}

func (p *Data) GetFloat64Def(path string, def float64) float64 {
	v, err := p.GetFloat64(path)
	if err != nil {
		return def
	}

	return v
}

func (p *Data) GetInt64(path string) (int64, error) {
	v, err := p.GetFloat64(path)
	if err != nil {
		return 0, err
	}

	return ToInt64(v), nil
}

func (p *Data) GetInt64Def(path string, def int64) int64 {
	v, err := p.GetInt64(path)
	if err != nil {
		return def
	}
	return v
}

func (p *Data) GetInt(path string) (int, error) {
	v, err := p.GetFloat64(path)
	if err != nil {
		return 0, err
	}

	return ToInt(v), nil
}

func (p *Data) GetIntDef(path string, def int) int {
	v, err := p.GetInt(path)
	if err != nil {
		return def
	}
	return v
}

type validator interface {
	Validate() error
}

func (p *Data) IsSet(path string) bool {
	_, err := Values(p.data).PathValue(path)
	return err == nil
}

func (p *Data) String() string {
	buf, err := yaml.Marshal(p.data)
	if err != nil {
		return err.Error()
	}
	return string(buf)
}

func (p *Data) Json() ([]byte, error) {
	b, err := yaml.Marshal(p.data)
	if err != nil {
		return nil, err
	}

	return yaml.YAMLToJSON(b)
}

func (p *Data) Yaml() ([]byte, error) {
	return yaml.Marshal(p.data)
}

func (p *Data) GetDefault(path string) (interface{}, bool) {
	return "", false
}

func (p *Data) GetDescription(path string) (string, bool) {
	return "", false
}
