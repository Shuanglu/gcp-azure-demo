resource "azuread_service_principal" "gcp_demo" {
  application_id = azuread_application.gcp_demo.application_id
}
