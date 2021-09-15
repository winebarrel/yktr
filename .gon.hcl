source    = ["./dist/yktr-macos_darwin_amd64/yktr"]
bundle_id = "com.github.winebarrel.yktr"

apple_id {
  username = "sugawara@winebarrel.jp"
  password = "@keychain:altool"
}

sign {
  application_identity = "00C6D6A53ED057485749B871ADADA1A6E475D024"
  entitlements_file    = "entitlements.plist"
}

zip {
  output_path = "./dist/yktr_0.9.4_darwin_amd64.zip"
}
