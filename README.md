# Hopfencloud

Just another cloud project.

## Installation

Compilation and installation are managed via Makefile.

If you plan on using mariadb or postgres as database backend instead of
sqlite (which you really should!), modify `hopfencloud.service` and uncomment
the corresponding section before executing `sudo make install`.

```bash
# Compile the binary to bin/
make build

# Install the service
sudo make install

# Clean local binaries
make clean
```

After installation, copy `/etc/hopfencloud/example.config.toml` to
`/etc/hopfencloud/config.toml` and edit the file to match your needs.

After `/etc/hopfencloud/config.toml` is in place, you can start the service:

```bash
systemctl start hopfencloud.service
```

## Uninstallation

Uninstallation is also managed via Makefile.

There are 2 ways of uninstallation:
- `sudo make uninstall`:
Removes the binary and service, leaves all data in place
- `sudo make purge`:
Removes all files created while installation and during operation

If you are using mysql or postgres as database backend, you have to drop the
database to remove all data. You may also want to remove the database user.

## Developer guidelines

### Style guidelines - go

#### Imports

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
