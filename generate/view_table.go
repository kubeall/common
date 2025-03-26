package generate

type ViewTable struct {
	ModelRef string   `gorm:"type:varchar(255)" json:"modelRef" yaml:"modelRef" description:"关联模型名"`
	Model    Model    `json:"model" yaml:"model" description:""`
	Searches []Search `json:"searches" yaml:"searches" description:"支持的搜索项"`
}
type Search struct {
	Type  string `json:"searchType" yaml:"searchType" enum:"equal|like" validate:"required" description:"搜索类型: 在支持搜索的情况下"`
	Field string `json:"field" yaml:"field" validate:"required" description:"搜索自动"`
}
