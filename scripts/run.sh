#!/bin/bash

set -eo pipefail

if [ $# -lt 1 ];
    then echo "illegal number of parameters"
    exit 1
fi

action=$1
shift

case $action in
    provision )
        command=apply
        ;;
    decommission )
        command=destroy
        ;;
    *)
        echo "Unsupported action: $action"
        exit 1
        ;;
esac

steps=[]
steps_from_file=$(yq '.steps' massdriver.yaml)
if [ "$steps_from_file" == null ] ; then
    echo "using default steps"
    steps=( '{"path": "src", "provisioner": "terraform"}' )
else
    echo "using steps from file"
      # Read in the steps (reverse them if this is destroy)
    if [ $command = "apply" ] ; then
        readarray steps < <(yq e -o=j -I=0 '.steps[]' massdriver.yaml)
    else
        readarray steps < <(yq e -o=j -I=0 '.steps | reverse | .[]' massdriver.yaml)
    fi
fi

# Iterate through each step
for step in "${steps[@]}"
do
    echo "on step $step"
    path=$(echo $step | yq '.path')
    echo "using path: $path"
    provisioner=$(echo $step | yq '.provisioner')
    echo "using provisioner: $provisioner"
    skip_on_delete=$(echo $step | yq '.skip_on_delete //false')
    echo "using skip_on_delete: $skip_on_delete"

    if [ $action = "decommission" ] && [ $skip_on_delete = "true" ]; then
        echo "Skipping step: 'skip_on_delete' is true"
        continue
    fi

    pushd $path
    echo "executing $action on $path"

    case $provisioner in
        terraform )
            tf_flags=""
            if [ $command = "destroy" ]; then
                tf_flags="-destroy"
            fi

            echo "executing terraform init"
            terraform init -no-color -input=false
            echo "executing terraform plan"
            terraform plan $tf_flags -out tf.plan -json | xo provisioner terraform report-progress -f -

            # check for invalid deletions if this isn't a destroy
            if [ $command = "apply" ]; then
                terraform show -json tf.plan > tfplan.json
                if [ -f validations.json ]; then
                    echo "evaluating OPA rules"
                    # opa eval --fail-defined --data /opa/terraform.rego --input tfplan.json --data validations.json "data.terraform.deletion_violations[x]" | xo provisioner opa report -f -
                fi
            fi

            echo "executing terraform apply"
            terraform apply $tf_flags -json tf.plan | xo provisioner terraform report-progress -f -
            ;;
        bicep )
            region_query=$(echo $step | yq '.region_query')
            service_principal_query=$(echo $step | yq '.service_principal_query // ""')
            delete_resource_group=$(echo $step | yq '.delete_resource_group // false')
            bicep_provision.sh $action $path $region_query $delete_resource_group $service_principal_query
            ;;
        *)
            echo "Unsupported provisioner: $provisioner"
            exit 1
            ;;
    esac

    popd
done

# Notify backend we completed
# xo deployment $action complete
