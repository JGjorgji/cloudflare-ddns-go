# Small utility to do ddns updates for Cloudflare

Provide the following variables to run it:  

- API_TOKEN - needs to have permissions to change the record
- ZONE_ID - your cloudflare zone id
- RECORD_NAME - which record you want to update, eg. home.example.com

Example systemd service below, it's recommended to use [systemd-creds](https://systemd.io/CREDENTIALS/) as in the example. You can also swap it out to `EnvironmentFile=` if you want to store them in a plain file.

```ini
[Unit]
Description=Update cloudflare dns record
After=network.target

[Service]
Type=oneshot
ProtectHome=read-only
ProtectSystem=strict
PrivateTmp=yes
RemoveIPC=yes
LoadCredentialEncrypted=vars:/path/to/my/creds/vars
ExecStart=/usr/bin/bash -c 'set -a && source "$CREDENTIALS_DIRECTORY/vars" && /usr/local/bin/cloudflare-ddns-go'

[Install]
WantedBy=multi-user.target
```

Timer on a 10 minute interval:  

```ini
[Unit]
Description=Periodically update cloudflare dns records if external IP changed

[Timer]
OnCalendar=*:0/10

[Install]
WantedBy=timers.target
```
