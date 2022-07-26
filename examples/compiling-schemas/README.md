# Terraform JSON Params and Connections

[Details here](https://www.terraform.io/docs/configuration/syntax-json.html#variable-blocks)


Compile `variables.tf.json`:

```shell
xo provisioner compile terraform --schema=./variables.schema.json --output=./variables.tf.json
```

Run terraform:

```shell
terraform init
terraform apply -auto-approve
cat RESULTS.md
```
