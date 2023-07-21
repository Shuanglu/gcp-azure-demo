resource "azurerm_role_assignment" "gcp_demo_ad_application" {
  scope                = "/subscriptions/${var.azure_sub_id}"
  role_definition_name = "Reader"
  principal_id         = azuread_service_principal.gcp_demo.object_id
  depends_on = [
    azuread_application.gcp_demo,
    azuread_service_principal.gcp_demo
  ]
}

resource "azurerm_role_assignment" "gcp_demo_rm_user_assigned_identity" {
  scope                = "/subscriptions/${var.azure_sub_id}"
  role_definition_name = "Reader"
  principal_id         = azurerm_user_assigned_identity.gcp_demo.principal_id
  depends_on = [
    azurerm_user_assigned_identity.gcp_demo
  ]
}
