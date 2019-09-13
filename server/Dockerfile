FROM nginx:1.17.3
COPY nginx.conf /etc/ngingx/ngingx.conf
COPY index.html /usr/share/nginx/html
RUN mkdir -p /usr/share/nginx/html/media
RUN mkdir -p /usr/share/nginx/html/logs
RUN chown -R $USER:$USER /usr/share/nginx/html
RUN chmod -R 755 /usr/share/nginx
RUN touch /usr/share/nginx/html/media/index.html
VOLUME /usr/share/nginx/html/logs
