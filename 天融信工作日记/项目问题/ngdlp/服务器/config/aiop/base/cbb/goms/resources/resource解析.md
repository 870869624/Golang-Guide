# 1. 配置文件结构

## 终端连接日志消息处理器资源配置

```yaml
processor_resources:
    - label: "processor_resource_sensitive_category_level"
      bloblang: |
       meta catalog   = "sensitive_category_level"
       root           = this
       root.message_type           = deleted()
```

# 2. 字段说明

- *label*: 处理器的唯一标识名称
  - 命名规范：*processor_resource_{资源类型}*
  - 示例：*processor_resource_sensitive_category_level*
- *bloblang*: 数据转换和路由脚本
  - 使用 *Bloblang* 语言定义如何处理和路由消息

# 3. Bloblang 脚本解析

```bloblang
# 设置元数据 - catalog 标识资源类型
meta catalog = "sensitive_category_level"

# 将输入数据赋值给 root（消息主体）
root = this

# 设置 message_type 字段，并删除原始字段
root.message_type = deleted()
```

# 4. 完整的敏感文件资源配置示例

根据你代码中的 4 种资源类型，应该有以下配置文件：

## (1) files 资源配置
```yaml
# processor_resource_sensitive_file.yaml
processor_resources:
    - label: "processor_resource_sensitive_file"
      bloblang: |
        meta catalog = "sensitive_file"
        root = this
        root.message_type = deleted()
```

## (2) categories 资源配置
```yaml
# processor_resource_sensitive_file_category.yaml
processor_resources:
    - label: "processor_resource_sensitive_file_category"
      bloblang: |
        meta catalog = "sensitive_file_category"
        root = this
        root.message_type = deleted()
```

## (3) levels 资源配置
```yaml
# processor_resource_sensitive_file_level.yaml
processor_resources:
    - label: "processor_resource_sensitive_file_level"
      bloblang: |
        meta catalog = "sensitive_file_level"
        root = this
        root.message_type= deleted()
```

## (4) category_levels 资源配置（你已提供）
```yaml
# processor_resource_sensitive_category_level.yaml
processor_resources:
    - label: "processor_resource_sensitive_category_level"
      bloblang: |
        meta catalog = "sensitive_category_level"
        root = this
        root.message_type= deleted()
```

# 5. 与发送代码的对应关系

你的 Go 代码中：

```go
data := map[string][]map[string]any{
    "files":           sensitiveFiles,           // → processor_resource_sensitive_file
    "categories":      sensitiveFileCategories,  // → processor_resource_sensitive_file_category
    "levels":          sensitiveFileLevels,      // → processor_resource_sensitive_file_level
    "category_levels": sensitiveFileCategoryLevels, // → processor_resource_sensitive_category_level
}
```

每个 key 对应一个 goms processor resource 配置，通过 meta catalog 来路由到不同的处理器。

# 6. 高级用法（可选）

如果需要更复杂的路由或数据转换：
```yaml
processor_resources:
    - label: "processor_resource_sensitive_file"
      bloblang: |
        # 设置资源类型
        meta catalog = "sensitive_file"
        
        # 数据转换
        root.id = this.id
        root.filename = this.filename
        root.md5 = this.md5
        root.sha1 = this.sha1
        root.category = this.category
        root.level = this.level
        root.client_id = this.client_id
        root.tenant_id = this.tenant_id
        root.created_at = this.created_at
        root.updated_at = this.updated_at
        
        # 添加处理标记
        root.processed_by = "goms"
        root.process_timestamp = now()
        
        # 删除不需要的字段
        root.message_type= deleted()
        root.internal_field = deleted()
```

# 7. 文件位置规范

根据你提供的路径：
```bash
F:\gopath\src\aiop\deploy\config\aiop\base\cbb\goms\resources\
├── processor_resource_sensitive_file.yaml
├── processor_resource_sensitive_file_category.yaml
├── processor_resource_sensitive_file_level.yaml
└── processor_resource_sensitive_category_level.yaml
```

总结：goms 的 resource 配置就是使用 YAML 格式定义处理器，通过 Bloblang 脚本来指定数据的路由（meta catalog）和转换规则（root）。每个资源类型对应一个独立的配置文件。