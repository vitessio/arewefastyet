[program:benchmark]
directory=/home/akilan/Documents/Vitess/arewefastyet
command=/home/akilan/Documents/Vitess/arewefastyet/benchmark/bin/gunicorn server:app -b localhost:8000
autostart=true
autorestart=true