# Dobrify Bot

Notify about new prizes appearing in Dobry Cola contest

## Deploy

### Install `supervisor`

```bash
sudo apt update && sudo apt install supervisor
```

```conf
; /etc/supervisor/conf.d/dobrify.conf
[program:dobrify-bot]
directory=/var/www/app
command=/var/www/app/dobrify bot
autostart=true
autorestart=true
stderr_logfile=/var/log/dobrify-bot.err.log
stdout_logfile=/var/log/dobrify-bot.out.log

[program:dobrify-cron]
directory=/var/www/app
command=/var/www/app/dobrify cron
autostart=true
autorestart=true
stderr_logfile=/var/log/dobrify-cron.err.log
stdout_logfile=/var/log/dobrify-cron.out.log
```

```bash
sudo supervisorctl reread
sudo supervisorctl update
```

### Add after deploy script to restart services

```sh
# /var/www/app/after_deploy.sh
sudo supervisorctl stop dobrify-bot dobrify-cron
mv /var/www/app/dobrify-linux /var/www/app/dobrify
sudo supervisorctl start dobrify-bot dobrify-cron
```

```bash
chmod +x /var/www/app/after_deploy.sh
```

### Setup deploy environment variables

Copy `.env.deploy.example` to `.env.deploy` and fill in the values

### Do deploy

From local machine

```bash
make deploy
```