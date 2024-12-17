package jenkins

import "encoding/xml"

type JenkinsJob struct {
	Name      string     `json:"name"`    // Job 名称
	Actions   []Action   `json:"actions"` // Job 的操作
	LastBuild Build      `json:"lastBuild"`
	Builds    []Build    `json:"builds"`
	Property  []Property `json:"property"`
}

// JenkinsView represents a Jenkins view
type JenkinsView struct {
	Name string `json:"name"` // 视图名称
	//Jobs []struct {
	//	Name string `json:"name"`
	//} `json:"jobs"`
	Jobs []JenkinsJob `json:"jobs"` // 视图中的 Job 列表
}
type Property struct {
	Class                string                `json:"_class,omitempty"`
	ParameterDefinitions []ParameterDefinition `json:"parameterDefinitions,omitempty"`
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
	Class                string                 ``
	ParameterDefinitions []ParameterDefinitions `json:"parameterDefinitions"`
}

type ParameterDefinitions struct {
	Name                  string `json:"name"`
	DefaultParameterValue struct {
		Value interface{} `json:"value"`
	} `json:"defaultParameterValue"`
}
type ParameterDefinition struct {
	Class                 string         `json:"_class"`
	DefaultParameterValue ParameterValue `json:"defaultParameterValue"`
	Description           string         `json:"description"`
	Name                  string         `json:"name"`
	Type                  string         `json:"type"`
}
type ParameterValue struct {
	Class string `json:"_class"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Build struct {
	Number    int       `json:"number"`    // 构建编号
	Result    string    `json:"result"`    // 构建结果
	Timestamp int64     `json:"timestamp"` // 构建时间戳
	URL       string    `json:"url"`       // 构建 URL
	ChangeSet ChangeSet `json:"changeSet"` // 变更集的信息
}
type ChangeSet struct {
	Kind  string          `json:"kind"`  // 变更集类型，如 "git" 或 "svn"
	Items []ChangeSetItem `json:"items"` // 变更项
}

type ChangeSetItem struct {
	CommitID      string   `json:"commitId"`      // 提交 ID
	Timestamp     int64    `json:"timestamp"`     // 提交时间
	Author        Author   `json:"author"`        // 提交作者
	Msg           string   `json:"msg"`           // 提交信息
	AffectedPaths []string `json:"affectedPaths"` // 受影响的文件路径
}
type Author struct {
	FullName string `json:"fullName"` // 作者全名
	Email    string `json:"email"`    // 作者邮箱
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
