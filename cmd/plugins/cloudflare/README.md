# Cloudflare Plugin

Config as follows:

```yaml
provider: cloudflare
username: # your cloudflare email
password: # your cloudflare global api key
```

**OR**

```yaml
provider: cloudflare
username: 
password: # your cloudflare API Token
```

If you use the API Token instead of the global API key you must give the token permission to `List Zones` & `List, Update, Create DNS Records`