up:
	@docker-compose up -d

down:
	@docker-compose down

redis:
	@docker exec -it redis redis-cli

stack:
	@docker-compose -f docker-compose-redis-stack.yaml up -d
