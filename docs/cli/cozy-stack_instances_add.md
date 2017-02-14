## cozy-stack instances add

Manage instances of a stack

### Synopsis



cozy-stack instances add allows to create an instance on the cozy for a
given domain.


```
cozy-stack instances add [domain]
```

### Options

```
      --apps stringSlice   Apps to be preinstalled
      --dev                To create a development instance
      --email string       The email of the owner
      --locale string      Locale of the new cozy instance (default "en")
      --tz string          The timezone for the user
```

### Options inherited from parent commands

```
      --admin-host string      administration server host (default "localhost")
      --admin-port int         administration server port (default 6060)
      --assets string          path to the directory with the assets (use the packed assets by default)
  -c, --config string          configuration file (default "$HOME/.cozy.yaml")
      --couchdb-url string     CouchDB URL (default "http://localhost:5984/")
      --fs-url string          filesystem url (default "file://localhost//storage")
      --host string            server host (default "localhost")
      --log-level string       define the log level (default "info")
      --mail-disable-tls       disable smtp over tls
      --mail-host string       mail smtp host (default "localhost")
      --mail-password string   mail smtp password
      --mail-port int          mail smtp port (default 465)
      --mail-username string   mail smtp username
  -p, --port int               server port (default 8080)
      --subdomains string      how to structure the subdomains for apps (can be nested or flat) (default "nested")
```

### SEE ALSO
* [cozy-stack instances](cozy-stack_instances.md)	 - Manage instances of a stack

###### Auto generated by spf13/cobra on 10-Feb-2017