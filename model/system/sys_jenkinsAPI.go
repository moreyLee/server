package system

import "encoding/xml"

type JenkinsJob struct {
	Actions   []Action `json:"actions"`
	LastBuild Build    `json:"lastBuild"`
	Name      string   `json:"name"`
	Builds    []Build  `json:"builds"`
}

// SCM 结构体用于存储 SCM 配置信息
type SCM struct {
	ConfigVersion     string `xml:"configVersion"`
	UserRemoteConfigs struct {
		URLs []string `xml:"hudson.plugins.git.UserRemoteConfig>url"`
	} `xml:"userRemoteConfigs"`
	Branches []string `xml:"branches>hudson.plugins.git.BranchSpec>name"`
}
type Action struct {
	ParameterDefinitions []ParameterDefinitions `json:"parameterDefinitions"`
}
type ParameterDefinitions struct {
	Name                  string `json:"name"`
	DefaultParameterValue struct {
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
}

// JenkinsView represents a Jenkins view
type JenkinsView struct {
	Name string `json:"name"`
	Jobs []struct {
		Name string `json:"name"`
	} `json:"jobs"`
}

// JobConfig 结构体用于存储 Jenkins Job 的配置信息
type JobConfig struct {
	XMLName xml.Name `xml:"project"`
	//SCM     SCM      `xml:"scm"`
	SCM struct {
		ConfigVersion     string `xml:"configVersion"`
		UserRemoteConfigs struct {
			URLs []string `xml:"hudson.plugins.git.UserRemoteConfig>url"`
		} `xml:"userRemoteConfigs"`
		Branches []string `xml:"branches>hudson.plugins.git.BranchSpec>name"`
	} `xml:"scm"`
}
type Build struct {
	Number    int    `json:"number"`
	Result    string `json:"result"`
	Timestamp int64  `json:"timestamp"`
	URL       string `json:"url"`
}
