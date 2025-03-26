package generate

// Model 模型定义，其中fields为对应数据库表的全量字段，业务场景上使用的模型应该为该模型的子模型，即字段数量小于等于权限自动
type Model struct {
	ID          string  `json:"id" yaml:"id" validate:"required" description:"模型ID"`
	Code        string  `gorm:"type:varchar(255)" json:"code" yaml:"code" validate:"alpha" description:"模型编码: 用于生成代码的model和数据库表"`
	Name        string  `gorm:"type:varchar(255)" json:"name" yaml:"name" validate:"required" description:"模型名称: 用于生产页面提示信息，未来支持AI翻译成国际化"`
	Description string  `gorm:"type:longtext" json:"description" yaml:"description" description:"模型说明"`
	Fields      []Field `json:"fields" yaml:"fields" description:"模型字段: 数据库的列名"`
}

// ChildModel 子模型，从Model提取字段形成一个新的模型，
type ChildModel struct {
	ModelRef    string            `gorm:"type:varchar(255)" json:"modelFef" yaml:"modelRef" description:"关联模型"`
	Code        string            `gorm:"type:varchar(255)" json:"code" yaml:"code" validate:"alpha" description:"模型编码: 用于生成代码的model和数据库表"`
	Name        string            `gorm:"type:varchar(255)" json:"name" yaml:"name" validate:"required" description:"模型名称: 用于生产页面提示信息，未来支持AI翻译成国际化"`
	Description string            `gorm:"type:longtext" json:"description" yaml:"description" description:"模型说明"`
	Fields      map[string]Extend `json:"fields" yaml:"fields" description:"字段: key为字段名称，extend为额外说明"`
}

// Extend 子模型中字段额外说明
type Extend struct {
}

// Field 字段说明
type Field struct {
	Code         string            `gorm:"type:varchar(255)" json:"code" yaml:"code" validate:"alpha" description:"字段编码: 用于生成代码的model字段和数据库表字段"`
	Name         string            `gorm:"type:varchar(255)" json:"name" yaml:"name" validate:"required" description:"字段名称: 用于生产页面提示信息，未来支持AI翻译成国际化"`
	Description  string            `gorm:"type:longtext" json:"description" yaml:"description" description:"字段说明"`
	DataType     string            `gorm:"type:varchar(255)" json:"dataType" yaml:"dataType" validate:"oneof=string bool int uint int64 float64 time object" description:"数据类型"`
	DefaultValue string            `json:"defaultValue" yaml:"defaultValue" description:"默认值"`
	Enums        map[string]string `json:"enums" yaml:"enums" description:"可以选值: key为数据库值，value为值说明"`
	MaxLength    uint              `json:"maxLength" yaml:"maxLength" description:"长度: 字符串类型数据时长度，0或者负数表示不限制，前端采用多行文本"`

	Tips  string      `json:"tips" yaml:"tips" description:"提示信息"`
	Rules []FieldRule `json:"rules" yaml:"rules" description:"校验规则"`
}

// FieldRule 字段校验规则
type FieldRule struct {
	Required bool   `json:"required" yaml:"required" description:"是否必须"`
	Type     string `json:"type" yaml:"type" description:"输入信息类型"`
	Pattern  string `json:"pattern" yaml:"pattern" description:"正则"`
	Message  string `json:"message" yaml:"message" description:"提示信息"`
}
