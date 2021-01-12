EXAMPLE_PATH=WORKS_ON_MY_MACHINE.md
.PHONY: all
all:
	terraform init
	terraform apply -auto-approve -var=path=${EXAMPLE_PATH}
	cat ${EXAMPLE_PATH}