

make build 
mkdir -p ../terraform-importer/.terraform/providers/registry.terraform.io/mrparkers/keycloak/99.0.0/darwin_amd64/
copy ./providers/registry.terraform.io/mrparkers/keycloak/99.0.0/darwin_amd64/terraform-provider-keycloak_v99.0.0  ../terraform-importer/.terraform/providers/registry.terraform.io/mrparkers/keycloak/99.0.0/darwin_amd64/terraform-provider-keycloak_v99.0.0

# make build-win
# mkdir -p ../terraform-importer/.terraform/providers/registry.terraform.io/mrparkers/keycloak/99.0.0/windows_amd64/
# copy ./providers/registry.terraform.io/mrparkers/keycloak/99.0.0/windows_amd64/terraform-provider-keycloak_v99.0.0.exe  ../terraform-importer/.terraform/providers/registry.terraform.io/mrparkers/keycloak/99.0.0/windows_amd64/terraform-provider-keycloak_v99.0.0.exe