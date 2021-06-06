build:
	docker-compose -f docker-compose.yaml build

up:
	docker-compose -f docker-compose.yaml up

tests:
	docker-compose -f docker-compose-test.yaml build && docker-compose -f docker-compose-test.yaml up