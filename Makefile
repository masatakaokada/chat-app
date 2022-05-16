heroku-go-deploy:
	docker build -t registry.heroku.com/chat-go-server/web -f docker/go/Dockerfile.production .
	docker push registry.heroku.com/chat-go-server/web
	git remote set-url heroku https://git.heroku.com/chat-go-server.git
	heroku container:release web

heroku-vue-deploy:
	docker build -t registry.heroku.com/chat-vue-client/web -f docker/vue/Dockerfile.production .
	docker push registry.heroku.com/chat-vue-client/web
	git remote set-url heroku https://git.heroku.com/chat-vue-client.git
	heroku container:release web
