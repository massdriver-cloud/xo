#!/bin/bash
action=$1
step=$2
region_query=$3
delete_resource_group=$4
service_principal_query=$5

params=$(cat parameters.json)
error_event=$(echo '{"@level": "error", "@message": "error"}' | jq -rc .)
success_event=$(echo '{"@level": "info", "@message": "Apply complete!"}' | jq -rc .)
service_principal_default_query='.connections.value.azure_service_principal'
service_principal_fallback_query="${service_principal_query:-${service_principal_default_query-default}}"
region=$(echo $params | jq -r $region_query)

echo "Authorizing to azure using service principal"
service_principal=$(echo $params | jq $service_principal_fallback_query)
AZURE_CLIENT_ID=$(echo $service_principal | jq -r .data.client_id)
AZURE_CLIENT_SECRET=$(echo $service_principal | jq -r .data.client_secret)
AZURE_TENANT_ID=$(echo $service_principal | jq -r .data.tenant_id)

az login --service-principal -u $AZURE_CLIENT_ID -p $AZURE_CLIENT_SECRET -t $AZURE_TENANT_ID > /dev/null

if [ $? -eq 0 ] 
  then 
    echo "Authorized"
  else
    echo $error_event
    exit 1 
fi

pkg_name=$(echo $params | jq -r .md_metadata.value.name_prefix)

if [ $action = "decommission" ]; then
  echo "Decommissioning resources"

  echo '' > empty.bicep

  az deployment group create \
    --name "$pkg_name-$step" \
    --resource-group $pkg_name \
    --template-file empty.bicep > /dev/null
  
  rm empty.bicep

  if [ $? -eq 0 ] 
    then 
      echo "Resources decommissioned"
    else
      echo $error_event
      exit 1 
  fi

  echo "Deleting deployment group $pkg_name-$step"
  az deployment group delete \
    --name "$pkg_name-$step" \
    --resource-group $pkg_name

  if [ $? -eq 0 ] 
    then 
      echo "Deployment group deleted"
    else
      echo $error_event
      exit 1 
  fi

  if [$delete_resource_group = "true"]; then
    echo "Deleting resource group $pkg_name in region $region"
    az group delete 
      --name $pkg_name
      --region $region

    if [ $? -eq 0 ] 
      then 
        echo "Resource group deleted"
      else
        echo $error_event
        exit 1 
    fi
  fi

  for artifact_file in artifact_*.jq; do
    field=$(echo $artifact_file | sed 's/^artifact_\(.*\).jq$/\1/')
    echo "Deleting artifact for field $field"
  done

  echo $success_event
  exit 0
fi


echo "Creating resource group $pkg_name in region $region"

az group create --name $pkg_name --location $region > /dev/null

if [ $? -eq 0 ] 
  then 
    echo "Resource group created"
  else 
    echo $error_event
    exit 1 
fi

echo "Deploying bundle"
OUTPUTS=$(az deployment group create \
  --name "$pkg_name-$step" \
  --resource-group $pkg_name \
  --template-file template.bicep \
  --parameters "$params")

if [ $? -eq 0 ] 
  then 
    echo "Deployment completed"
  else 
    echo $error_event
    exit 1 
fi

OUTPUTS=$(echo $OUTPUTS | jq -r '{outputs: .properties.outputs}')

echo "Creating artifacts"

inputs_and_outputs=$(echo $params $OUTPUTS | jq -s add)

for artifact_file in artifact_*.jq; do
  field=$(echo $artifact_file | sed 's/^artifact_\(.*\).jq$/\1/')
  echo "Creating artifact for field $field"
  echo $inputs_and_outputs | jq -rc -f $artifact_file
  echo "Artifact $field created"
done

echo $success_event 
exit 0
