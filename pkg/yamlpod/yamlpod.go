package yamlpod

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type YamlStruct struct {
	Links    map[string]YamlLink    `yaml:"links"`
	MiniPods map[string]YamlMiniPod `yaml:"minipods,omitempty"`
}

type YamlLink struct {
	Name     string              `yaml:"name,omitempty"`
	Link     string              `yaml:"link,omitempty"`
	Icon     string              `yaml:"icon,omitempty"`
	Minipods []YamlMinipodInLink `yaml:"minipods,omitempty"`
}

type YamlMinipodInLink struct {
	MinipodKey string `yaml:"minipod"`
	ABC        string `yaml:"abc,omitempty"`
}

type YamlMiniPod struct {
	Name     string              `yaml:"name,omitempty"`
	Icon     string              `yaml:"icon,omitempty"`
	ABC      string              `yaml:"abc,omitempty"`
	Minipods []YamlMinipodInLink `yaml:"minipods,omitempty"`
}

func DownloadYaml(url string) (*YamlStruct, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Permissions-Policy", "interest-cohort=()")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return Unmarshal(resp.Body)
}

func Unmarshal(r io.Reader) (*YamlStruct, error) {
	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, r)
	if err == nil && n == 0 {
		err = errors.New("empty yaml file")
	}
	if err != nil {
		return nil, err
	}

	ys := &YamlStruct{}
	err = yaml.Unmarshal(buf.Bytes(), ys)
	if err != nil {
		return nil, err
	}

	return ys, nil
}

func Marshal(w io.Writer, ys *YamlStruct) error {
	out, err := yaml.Marshal(ys)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(out)
	n, err := io.Copy(w, buf)
	if err == nil && n == 0 {
		err = errors.New("empty yaml file")
	}
	return err
}
