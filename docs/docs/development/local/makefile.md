## Makefile abstraction

To get the development process smoother, we can rely on Makefile abstract for some often-used commands:

---

## Example of most used commands

### For serving mkdocs
```shell
make serve-docs
```

### For Testing
```shell
make test
```

### For Deploying with docker-compose
```shell
make rebuild
```

---

## All Makefile's targets

Refer to the current Makefile to validate all options:
```Makefile
{% include "snippets/Makefile" %}
```

{% include "discussions.md" %}
