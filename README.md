# Hopfencloud

Just another cloud project.

## Style guidelines

### Imports

Imports should be split into three parts:
- stdlib
- project
- external

```go
import (
    // stdlib
    "fmt"
    "html/template"

    // project
    "github.com/myOmikron/hopfencloud/models/conf"
    "github.com/myOmikron/hopfencloud/models/db"

    // external
    "gorm.io/gorm"
)
```
