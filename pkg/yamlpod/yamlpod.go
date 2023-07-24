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
	MiniPods map[string]YamlMiniPod `yaml:"minipods"`
}

type YamlLink struct {
	Name     string              `yaml:"name"`
	Link     string              `yaml:"link"`
	Minipods []YamlMinipodInLink `yaml:"minipods"`
}

type YamlMinipodInLink struct {
	MinipodKey string `yaml:"minipod"`
	ABC        string `yaml:"abc"`
}

type YamlMiniPod struct {
	Name     string              `yaml:"name"`
	Icon     string              `yaml:"icon"`
	ABC      string              `yaml:"abc"`
	Minipods []YamlMinipodInLink `yaml:"minipods"`
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

	// Write the body to the writer
	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, resp.Body)
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
