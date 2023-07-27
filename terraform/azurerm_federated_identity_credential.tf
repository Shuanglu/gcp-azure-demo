resource "azurerm_federated_identity_credential" "gcp_demo" {
  name                = "gcp-demo"
  resource_group_name = azurerm_resource_group.gcp_demo.name
  audience            = [var.azure_exchange_audience]
  issuer              = var.gcp_issuer
  parent_id           = azurerm_user_assigned_identity.gcp_demo.id
  subject             = var.gcp_sa_id
  depends_on = [
    azuread_application.gcp_demo
  ]
}
