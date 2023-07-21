resource "azuread_application_federated_identity_credential" "gcp_demo" {
  application_object_id = azuread_application.gcp_demo.object_id
  display_name          = "gcp-demo"
  audiences             = ["api://AzureADTokenExchange"]
  issuer                = var.gcp_issuer
  subject               = var.gcp_sa_id
  depends_on = [
    azuread_application.gcp_demo
  ]
}
