resource "azuread_application" "gcp_demo" {
  display_name = "gcp-demo"
}



output "gcp_demo_ad_application" {
  value = azuread_application.gcp_demo.application_id
}
