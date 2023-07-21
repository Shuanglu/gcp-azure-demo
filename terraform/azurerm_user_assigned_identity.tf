resource "azurerm_user_assigned_identity" "gcp_demo" {
  location            = azurerm_resource_group.gcp_demo.location
  name                = "gcp-demo"
  resource_group_name = azurerm_resource_group.gcp_demo.name
  depends_on = [
    azurerm_resource_group.gcp_demo
  ]
}

output "gcp_demo_rm_user_assigned_identity" {
  value = azurerm_user_assigned_identity.gcp_demo.client_id
}
