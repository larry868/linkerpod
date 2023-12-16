package yamlpod

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

// YamlStruct overall structure of linkerpod yaml configuration file
type YamlStruct struct {
	MiniPods    map[string]YamlMiniPod `yaml:"minipods,omitempty"` // minipods containing links
	SingleLinks map[string]YamlLink    `yaml:"links,omitempty"`    // standalone links
}

type YamlMiniPod struct {
	Separator string              `yaml:"separator,omitempty"`
	Name      string              `yaml:"name,omitempty"`
	Icon      string              `yaml:"icon,omitempty"`
	Links     map[string]YamlLink `yaml:"links,omitempty"`
	IsOpen    bool                `yaml:"isopen,omitempty"`
}

type YamlLink struct {
	Name string `yaml:"name,omitempty"`
	Link string `yaml:"link,omitempty"`
	Icon string `yaml:"icon,omitempty"`
}

var ErrGetYamlFile = errors.New("unable to get yaml setup file")

// DownloadYaml Get the yaml file at url and unmarshal its content
func DownloadYaml(url string) (*YamlStruct, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("DownloadYaml: %w, %w", ErrGetYamlFile, err)
	}
	req.Header.Set("Permissions-Policy", "interest-cohort=()")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DownloadYaml: %w, %w", ErrGetYamlFile, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("DownloadYaml: %w, %w", ErrGetYamlFile, errors.New(resp.Status))
	}

	return Unmarshal(url, resp.Body)
}

func Unmarshal(url string, r io.Reader) (*YamlStruct, error) {
	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, r)
	if err == nil {
		if n == 0 {
			err = errors.New("empty yaml file")
		} else if !strings.Contains(buf.String(), "links:") {
			err = fmt.Errorf("%q file is not a linkerpod yaml file", url)
		}
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
