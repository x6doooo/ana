define build_env
	eGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/ana_$(1) ./main.go
	tar zcvf ./ana_$(1).tar.gz ./bin/ana_$(1) ./conf/conf.$(1).toml ./start_$(1).sh
endef

all: clean prod

clean: clean_test

clean_test:
	rm -rf ./ana_test.tar.gz

clean_prod:
	rm -rf ./ana_prod.tar.gz	

test:
	$(call build_env,test)

prod:
	$(call build_env,prod)
